package disc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/Burzich/dvault/internal/dvault/kv"
)

type KV struct {
	mountPath string
	m         sync.Mutex
}

func NewKV(mountPath string) *KV {
	return &KV{
		mountPath: mountPath,
	}
}

func (k *KV) Save(_ context.Context, secretPath string, data map[string]interface{}, cas int) error {
	k.m.Lock()
	defer k.m.Unlock()

	oldData, err := k.read(secretPath)
	if errors.Is(err, os.ErrNotExist) {
		data := Data{
			Records: []kv.Record{
				{
					Data: data,
					Metadata: struct {
						CreatedTime    time.Time   `json:"created_time"`
						CustomMetadata interface{} `json:"custom_metadata"`
						DeletionTime   string      `json:"deletion_time"`
						Destroyed      bool        `json:"destroyed"`
						Version        int         `json:"version"`
					}{
						CreatedTime:    time.Now(),
						CustomMetadata: nil,
						DeletionTime:   "",
						Destroyed:      false,
						Version:        1,
					}},
			},
			Meta: kv.Meta{
				CasRequired:        false,
				CreatedTime:        time.Now(),
				CurrentVersion:     1,
				DeleteVersionAfter: "",
				MaxVersions:        0,
				OldestVersion:      1,
				UpdatedTime:        time.Now(),
			},
		}

		return k.write(secretPath, data)
	}

	if oldData.Meta.CurrentVersion != cas && oldData.Meta.CasRequired == true {
		return errors.New("cas version does not match")
	}

	oldData.Records = append(oldData.Records, kv.Record{
		Data: data,
		Metadata: struct {
			CreatedTime    time.Time   `json:"created_time"`
			CustomMetadata interface{} `json:"custom_metadata"`
			DeletionTime   string      `json:"deletion_time"`
			Destroyed      bool        `json:"destroyed"`
			Version        int         `json:"version"`
		}{
			CreatedTime:    time.Now(),
			CustomMetadata: nil,
			DeletionTime:   "",
			Destroyed:      false,
			Version:        0,
		},
	})
	oldData.Meta.CurrentVersion++
	oldData.Meta.UpdatedTime = time.Now()
	oldData.Meta.CurrentVersion = len(oldData.Records)

	return k.write(secretPath, oldData)
}

func (k *KV) UpdateConfig(_ context.Context, secretPath string, config kv.Config) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return nil
	}

	data.Meta.MaxVersions = config.MaxVersions
	data.Meta.DeleteVersionAfter = config.DeleteVersionAfter
	data.Meta.CasRequired = config.CasRequired

	return k.write(secretPath, data)
}

func (k *KV) ReadConfig(_ context.Context, secretPath string) (kv.Config, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return kv.Config{}, nil
	}

	return kv.Config{
		CasRequired:        data.Meta.CasRequired,
		DeleteVersionAfter: data.Meta.DeleteVersionAfter,
		MaxVersions:        data.Meta.MaxVersions,
	}, nil
}

func (k *KV) GetMeta(_ context.Context, secretPath string) (kv.Meta, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return kv.Meta{}, nil
	}

	for i, record := range data.Records {
		data.Meta.Versions[strconv.Itoa(i)] = struct {
			CreatedTime  time.Time `json:"created_time"`
			DeletionTime string    `json:"deletion_time"`
			Destroyed    bool      `json:"destroyed"`
		}{
			CreatedTime:  record.Metadata.CreatedTime,
			DeletionTime: record.Metadata.DeletionTime,
			Destroyed:    record.Metadata.Destroyed,
		}
	}

	return data.Meta, nil
}

func (k *KV) UpdateMeta(_ context.Context, secretPath string, meta kv.Meta) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return err
	}

	data.Meta.MaxVersions = meta.MaxVersions
	data.Meta.CasRequired = meta.CasRequired
	data.Meta.CustomMetadata = meta.CustomMetadata
	data.Meta.DeleteVersionAfter = meta.DeleteVersionAfter

	return k.write(secretPath, data)
}

func (k *KV) DeleteMeta(_ context.Context, secretPath string) error {
	k.m.Lock()
	defer k.m.Unlock()

	return k.delete(secretPath)
}

func (k *KV) UndeleteVersion(_ context.Context, secretPath string, version int) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime != "" && record.Metadata.Version == version {
			record.Metadata.DeletionTime = ""
			data.Records[i-1] = record

			return k.write(secretPath, data)
		}
	}

	return nil
}

func (k *KV) DeleteVersion(_ context.Context, secretPath string, version int) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" && record.Metadata.Version == version {
			record.Metadata.DeletionTime = time.Now().String()
			data.Records[i-1] = record

			return k.write(secretPath, data)
		}
	}

	return errors.New("not found")
}

func (k *KV) Undelete(_ context.Context, secretPath string) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime != "" {
			record.Metadata.DeletionTime = ""
			data.Records[i-1] = record

			return k.write(secretPath, data)
		}
	}

	return nil
}

func (k *KV) Delete(_ context.Context, secretPath string) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" {
			record.Metadata.DeletionTime = time.Now().String()
			data.Records[i-1] = record

			return k.write(secretPath, data)
		}
	}

	return errors.New("not found")
}

func (k *KV) Get(_ context.Context, secretPath string) (kv.Record, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return kv.Record{}, err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" {
			record.Metadata.CustomMetadata = data.Meta.CustomMetadata
			return record, nil
		}
	}

	return kv.Record{}, errors.New("not found")
}

func (k *KV) GetVersion(_ context.Context, secretPath string, version int) (kv.Record, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.read(secretPath)
	if err != nil {
		return kv.Record{}, err
	}

	index := slices.IndexFunc(data.Records, func(record kv.Record) bool {
		return record.Metadata.Version == version && record.Metadata.DeletionTime == "" && record.Metadata.Destroyed == false
	})

	if index == -1 {
		return kv.Record{}, errors.New("version not found")
	}

	record := data.Records[index]
	record.Metadata.CustomMetadata = data.Meta.CustomMetadata
	return record, nil
}

func (k *KV) read(secretPath string) (Data, error) {
	p := path.Join(k.mountPath, secretPath)
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(p))

	b, err := os.ReadFile(pathEncoded)
	if err != nil {
		return Data{}, err
	}

	var data Data
	err = json.Unmarshal(b, &data)
	if err != nil {
		return Data{}, err
	}

	return data, nil
}

func (k *KV) delete(secretPath string) error {
	p := path.Join(k.mountPath, secretPath)
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(p))

	return os.Remove(pathEncoded)
}

func (k *KV) write(secretPath string, data Data) error {
	p := path.Join(k.mountPath, secretPath)
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(p))

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(pathEncoded, d, 0644)
	if err != nil {
		return err
	}

	return nil
}
