package standart

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/Burzich/dvault/internal/dvault/kv"
	"github.com/Burzich/dvault/internal/dvault/storage"
	"github.com/Burzich/dvault/internal/tools"
)

type KV struct {
	configPath string
	dataPath   string
	storage    storage.Storage
	encryptor  tools.Encryptor
}

func NewKV(configPath string, dataPath string, config kv.Config, storage storage.Storage, encryptor tools.Encryptor) (*KV, error) {
	k := KV{
		configPath: configPath,
		dataPath:   dataPath,
		storage:    storage,
		encryptor:  encryptor,
	}

	err := k.writeConfig(config)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func RestoreKV(configPath string, dataPath string, s storage.Storage, encryptor tools.Encryptor) (*KV, error) {
	k := KV{
		configPath: configPath,
		dataPath:   dataPath,
		storage:    s,
		encryptor:  encryptor,
	}

	return &k, nil
}

func (k *KV) Save(_ context.Context, secretPath string, data map[string]interface{}, cas int) error {
	oldData, err := k.readData(secretPath)
	if errors.Is(err, storage.ErrPathNotFound) {
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

		return k.writeData(secretPath, data)
	}

	if oldData.Meta.CurrentVersion != cas && oldData.Meta.CasRequired {
		return kv.ErrCas
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
	oldData.Meta.OldestVersion++

	return k.writeData(secretPath, oldData)
}

func (k *KV) UpdateConfig(_ context.Context, config kv.Config) error {
	return k.writeConfig(config)
}

func (k *KV) Destroy(_ context.Context, secretPath string, versions []int) error {
	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if !record.Metadata.Destroyed && slices.Contains(versions, record.Metadata.Version) {
			record.Metadata.Destroyed = true
			record.Data = nil
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return kv.ErrVersionNotFound
}

func (k *KV) GetConfig(_ context.Context) (kv.Config, error) {
	return k.readConfig()
}

func (k *KV) GetMeta(_ context.Context, secretPath string) (kv.Meta, error) {
	data, err := k.readData(secretPath)
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
	return k.deleteData(secretPath)
}

func (k *KV) UndeleteVersion(_ context.Context, secretPath string, version int) error {
	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if !record.Metadata.Destroyed && record.Metadata.DeletionTime != "" && record.Metadata.Version == version {
			record.Metadata.DeletionTime = ""
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return nil
}

func (k *KV) DeleteVersion(_ context.Context, secretPath string, versions []int) error {
	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if !record.Metadata.Destroyed && record.Metadata.DeletionTime == "" && slices.Contains(versions, record.Metadata.Version) {
			record.Metadata.DeletionTime = time.Now().String()
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return kv.ErrVersionNotFound
}

func (k *KV) Undelete(_ context.Context, secretPath string) error {
	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if !record.Metadata.Destroyed && record.Metadata.DeletionTime != "" {
			record.Metadata.DeletionTime = ""
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return nil
}

func (k *KV) Delete(_ context.Context, secretPath string) error {
	data, err := k.readData(secretPath)
	if err != nil {
		return err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if !record.Metadata.Destroyed && record.Metadata.DeletionTime == "" {
			record.Metadata.DeletionTime = time.Now().String()
			data.Records[i-1] = record

			return k.writeData(secretPath, data)
		}
	}

	return kv.ErrPathNotFound
}

func (k *KV) Get(_ context.Context, secretPath string) (kv.Record, error) {
	data, err := k.readData(secretPath)
	if err != nil {
		return kv.Record{}, err
	}

	for i := len(data.Records); i != 0; i-- {
		record := data.Records[i-1]
		if !record.Metadata.Destroyed && record.Metadata.DeletionTime == "" {
			record.Metadata.CustomMetadata = data.Meta.CustomMetadata
			return record, nil
		}
	}

	return kv.Record{}, kv.ErrPathNotFound
}

func (k *KV) GetVersion(_ context.Context, secretPath string, version int) (kv.Record, error) {
	data, err := k.readData(secretPath)
	if err != nil {
		return kv.Record{}, err
	}

	index := slices.IndexFunc(data.Records, func(record kv.Record) bool {
		return record.Metadata.Version == version && record.Metadata.DeletionTime == "" && !record.Metadata.Destroyed
	})

	if index == -1 {
		return kv.Record{}, kv.ErrVersionNotFound
	}

	record := data.Records[index]
	record.Metadata.CustomMetadata = data.Meta.CustomMetadata
	return record, nil
}

func (k *KV) DestroyKV(_ context.Context) error {
	err := k.deleteConfig()
	if err != nil {
		return err
	}

	err = k.deleteAllData()
	if err != nil {
		return err
	}

	return nil
}

func (k *KV) readConfig() (kv.Config, error) {
	p := filepath.Join(k.configPath, "config")

	b, err := k.storage.Get(context.Background(), p)
	if errors.Is(err, storage.ErrPathNotFound) {
		return kv.Config{}, kv.ErrPathNotFound
	}
	if err != nil {
		return kv.Config{}, err
	}

	decryptedData, err := k.encryptor.Decrypt(b)
	if err != nil {
		return kv.Config{}, err
	}

	var data kv.Config
	err = json.Unmarshal(decryptedData, &data)
	if err != nil {
		return kv.Config{}, err
	}

	return data, nil
}

func (k *KV) deleteConfig() error {
	p := filepath.Join(k.configPath, "config")

	return k.Delete(context.Background(), p)
}

func (k *KV) writeConfig(data kv.Config) error {
	p := filepath.Join(k.configPath, "config")

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	d, err = k.encryptor.Encrypt(d)
	if err != nil {
		return err
	}

	err = k.storage.Put(context.Background(), p, d)
	if err != nil {
		return err
	}

	return nil
}

func (k *KV) readData(secretPath string) (Data, error) {
	p := filepath.Join(k.dataPath, secretPath)

	b, err := k.storage.Get(context.Background(), p)
	if errors.Is(err, storage.ErrPathNotFound) {
		return Data{}, kv.ErrPathNotFound
	}
	if err != nil {
		return Data{}, err
	}

	decryptedData, err := k.encryptor.Decrypt(b)
	if err != nil {
		return Data{}, err
	}

	var data Data
	err = json.Unmarshal(decryptedData, &data)
	if err != nil {
		return Data{}, err
	}

	return data, nil
}

func (k *KV) deleteData(secretPath string) error {
	p := filepath.Join(k.dataPath, secretPath)

	return k.storage.Delete(context.Background(), p)
}

func (k *KV) deleteAllData() error {
	return k.storage.Delete(context.Background(), k.dataPath)
}

func (k *KV) writeData(secretPath string, data Data) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encryptedData, err := k.encryptor.Encrypt(d)
	if err != nil {
		return err
	}

	err = k.storage.Put(context.Background(), filepath.Join(k.dataPath, secretPath), encryptedData)
	if err != nil {
		return err
	}

	return nil
}
