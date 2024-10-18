package disc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/Burzich/dvault/internal/dvault/kv"
	"github.com/Burzich/dvault/internal/dvault/tools"
)

type KV struct {
	configPath string
	dataPath   string
	encryptor  tools.Encryptor
	m          sync.Mutex
}

func NewKV(configPath string, dataPath string, config kv.KVConfig, encryptor tools.Encryptor) (*KV, error) {
	k := KV{
		configPath: configPath,
		dataPath:   dataPath,
		encryptor:  encryptor,
	}

	err := os.MkdirAll(k.dataPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(k.configPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = k.writeConfig(config)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func RestoreKV(configPath string, dataPath string, encryptor tools.Encryptor) (*KV, error) {
	k := KV{
		configPath: configPath,
		dataPath:   dataPath,
		encryptor:  encryptor,
	}

	return &k, nil
}

func (k *KV) Save(_ context.Context, secretPath string, data map[string]interface{}, cas int) error {
	k.m.Lock()
	defer k.m.Unlock()

	oldData, err := k.readData(secretPath)
	if errors.Is(err, os.ErrNotExist) {
		data := Data{
			Records: []kv.KVRecord{
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
			Meta: kv.KVMeta{
				CasRequired:        false,
				CreatedTime:        time.Now(),
				CurrentVersion:     1,
				DeleteVersionAfter: "",
				MaxVersions:        0,
				OldestVersion:      1,
				UpdatedTime:        time.Now(),
			},
		}

		return k.writeData(secretPath, data)
	}

	if oldData.Meta.CurrentVersion != cas && oldData.Meta.CasRequired == true {
		return errors.New("cas version does not match")
	}

	oldData.Records = append(oldData.Records, kv.KVRecord{
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
	oldData.Meta.OldestVersion++

	return k.writeData(secretPath, oldData)
}

func (k *KV) Update(_ context.Context, secretPath string, data map[string]interface{}) error {
	k.m.Lock()
	defer k.m.Unlock()

	oldData, err := k.readData(secretPath)
	if errors.Is(err, os.ErrNotExist) {
		data := Data{
			Records: []kv.KVRecord{
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
			Meta: kv.KVMeta{
				CasRequired:        false,
				CreatedTime:        time.Now(),
				CurrentVersion:     1,
				DeleteVersionAfter: "",
				MaxVersions:        0,
				OldestVersion:      1,
				UpdatedTime:        time.Now(),
			},
		}

		return k.writeData(secretPath, data)
	}

	newRecord := kv.KVRecord{
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
	}

	oldData.Meta.UpdatedTime = time.Now()
	for i := len(oldData.Records); i != 0; i-- {
		record := oldData.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" {
			oldData.Records[i-1] = newRecord

			return k.writeData(secretPath, oldData)
		}
	}

	return k.writeData(secretPath, oldData)
}

func (k *KV) UpdateConfig(_ context.Context, config kv.KVConfig) error {
	k.m.Lock()
	defer k.m.Unlock()

	return k.writeConfig(config)
}

func (k *KV) Destroy(_ context.Context, secretPath string, versions []int) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && slices.Contains(versions, record.Metadata.Version) {
			record.Metadata.Destroyed = true
			record.Data = nil
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return errors.New("not found")
}

func (k *KV) GetConfig(_ context.Context) (kv.KVConfig, error) {
	k.m.Lock()
	defer k.m.Unlock()

	return k.readConfig()
}

func (k *KV) GetMeta(_ context.Context, secretPath string) (kv.KVMeta, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return kv.KVMeta{}, nil
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

func (k *KV) UpdateMeta(_ context.Context, secretPath string, meta kv.KVMeta) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	data.Meta.MaxVersions = meta.MaxVersions
	data.Meta.CasRequired = meta.CasRequired
	data.Meta.CustomMetadata = meta.CustomMetadata
	data.Meta.DeleteVersionAfter = meta.DeleteVersionAfter

	return k.writeData(secretPath, data)
}

func (k *KV) DeleteMeta(_ context.Context, secretPath string) error {
	k.m.Lock()
	defer k.m.Unlock()

	return k.deleteData(secretPath)
}

func (k *KV) UndeleteVersion(_ context.Context, secretPath string, version int) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime != "" && record.Metadata.Version == version {
			record.Metadata.DeletionTime = ""
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return nil
}

func (k *KV) DeleteVersion(ctx context.Context, secretPath string, versions []int) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" && slices.Contains(versions, record.Metadata.Version) {
			record.Metadata.DeletionTime = time.Now().String()
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return errors.New("not found")
}

func (k *KV) Undelete(_ context.Context, secretPath string) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime != "" {
			record.Metadata.DeletionTime = ""
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return nil
}

func (k *KV) Delete(_ context.Context, secretPath string) error {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" {
			record.Metadata.DeletionTime = time.Now().String()
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return errors.New("not found")
}

func (k *KV) Get(_ context.Context, secretPath string) (kv.KVRecord, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return kv.KVRecord{}, err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if record.Metadata.Destroyed == false && record.Metadata.DeletionTime == "" {
			record.Metadata.CustomMetadata = data.Meta.CustomMetadata
			return record, nil
		}
	}

	return kv.KVRecord{}, errors.New("not found")
}

func (k *KV) GetVersion(_ context.Context, secretPath string, version int) (kv.KVRecord, error) {
	k.m.Lock()
	defer k.m.Unlock()

	data, err := k.readData(secretPath)
	if err != nil {
		return kv.KVRecord{}, err
	}

	index := slices.IndexFunc(data.Records, func(record kv.KVRecord) bool {
		return record.Metadata.Version == version && record.Metadata.DeletionTime == "" && record.Metadata.Destroyed == false
	})

	if index == -1 {
		return kv.KVRecord{}, errors.New("version not found")
	}

	record := data.Records[index]
	record.Metadata.CustomMetadata = data.Meta.CustomMetadata
	return record, nil
}

func (k *KV) readConfig() (kv.KVConfig, error) {
	pathEncoded := base64.StdEncoding.EncodeToString([]byte("config"))
	p := filepath.Join(k.configPath, pathEncoded)

	b, err := os.ReadFile(p)
	if err != nil {
		return kv.KVConfig{}, err
	}

	var data kv.KVConfig
	err = json.Unmarshal(b, &data)
	if err != nil {
		return kv.KVConfig{}, err
	}

	return data, nil
}

func (k *KV) deleteConfig(secretPath string) error {
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(secretPath))
	p := filepath.Join(k.configPath, pathEncoded)

	return os.Remove(p)
}

func (k *KV) writeConfig(data kv.KVConfig) error {
	pathEncoded := base64.StdEncoding.EncodeToString([]byte("config"))
	p := filepath.Join(k.configPath, pathEncoded)

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(p, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (k *KV) readData(secretPath string) (Data, error) {
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(secretPath))
	p := filepath.Join(k.dataPath, pathEncoded)

	b, err := os.ReadFile(p)
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

func (k *KV) deleteData(secretPath string) error {
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(secretPath))
	p := filepath.Join(k.dataPath, pathEncoded)

	return os.Remove(p)
}

func (k *KV) writeData(secretPath string, data Data) error {
	pathEncoded := base64.StdEncoding.EncodeToString([]byte(secretPath))
	p := filepath.Join(k.dataPath, pathEncoded)

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(p, d, 0644)
	if err != nil {
		return err
	}

	return nil
}
