package dvault

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	kv2 "github.com/Burzich/dvault/internal/dvault/kv"
	"github.com/Burzich/dvault/internal/dvault/kv/disc"
	"github.com/Burzich/dvault/internal/dvault/tools"
)

type DVault struct {
	EncryptionKey string
	Logger        *slog.Logger
	mountPath     string

	mu sync.RWMutex
	kv map[string]kv2.KV
}

func NewDVault(logger *slog.Logger, mountPath string) *DVault {
	return &DVault{
		EncryptionKey: "",
		mountPath:     mountPath,
		Logger:        logger,
		kv:            make(map[string]kv2.KV),
	}
}

func (d *DVault) Unseal(ctx context.Context) (Response, error) {
	return Response{}, nil
}

func (d *DVault) Seal(ctx context.Context) (Response, error) {
	return Response{}, nil
}

func (d *DVault) SealStatus(ctx context.Context) (Response, error) {
	return Response{}, nil
}

func (d *DVault) GetKVSecret(ctx context.Context, mount string, secretPath string) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	data, err := d.kv[mount].Get(ctx, secretPath)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.Data = data
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) GetKVSecretByVersion(ctx context.Context, mount string, secretPath string, version int) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	data, err := d.kv[mount].GetVersion(ctx, secretPath, version)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.Data = data
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) SaveKVSecret(ctx context.Context, mount string, secretPath string, data map[string]interface{}, cas int) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].Save(ctx, secretPath, data, cas)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.Data = data
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) UpdateKVSecret(ctx context.Context, mount string, secretPath string, data map[string]interface{}) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].Update(ctx, secretPath, data)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.Data = data
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) DeleteKVSecret(ctx context.Context, mount string, secretPath string) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].Delete(ctx, secretPath)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) UndeleteKVSecret(ctx context.Context, mount string, secretPath string) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].Undelete(ctx, secretPath)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) DeleteKVSecretByVersion(ctx context.Context, mount string, secretPath string, version []int) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].DeleteVersion(ctx, secretPath, version)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) UndeleteKVSecretByVersion(ctx context.Context, mount string, secretPath string, version int) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].UndeleteVersion(ctx, secretPath, version)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) DestroyKVSecret(ctx context.Context, mount string, secretPath string, version []int) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].Destroy(ctx, secretPath, version)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) UpdateKVConfig(ctx context.Context, mount string, config kv2.KVConfig) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].UpdateConfig(ctx, config)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) GetKVConfig(ctx context.Context, mount string) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	data, err := d.kv[mount].GetConfig(ctx)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.Data = data
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) GetKVMeta(ctx context.Context, mount string, secretPath string) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	data, err := d.kv[mount].GetMeta(ctx, secretPath)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.Data = data
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) UpdateKVMeta(ctx context.Context, mount string, secretPath string, meta kv2.KVMeta) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].UpdateMeta(ctx, secretPath, meta)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) DeleteKVMeta(ctx context.Context, mount string, secretPath string) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.kv[mount]; !ok {
		return Response{}, fmt.Errorf("kv %s does not exist", mount)
	}

	err := d.kv[mount].DeleteMeta(ctx, secretPath)
	if err != nil {
		return Response{}, err
	}

	var response Response
	response.MountType = "kv"
	response.RequestId = tools.GenerateXRequestID()

	return response, nil
}

func (d *DVault) CreateMount(_ context.Context, path string, mount CreateMount) (Response, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var response Response
	response.RequestId = tools.GenerateXRequestID()

	if strings.Contains(path, ".") {
		return response, errors.New("path can't contain '.'")
	}

	path = filepath.Clean(path)
	if _, ok := d.kv[path]; ok {
		return response, errors.New("mount already exist")
	}

	switch mount.Type {
	case "kv":
		cfg, err := kv2.CreateConfigFromMap(mount.Config)
		if err != nil {
			return response, err
		}

		kv, err := disc.NewKV(filepath.Join(d.mountPath, path), filepath.Join(d.mountPath, "data", path), cfg, d.EncryptionKey)
		if err != nil {
			return response, err
		}
		d.kv[path] = kv
	default:
		return response, errors.New("unknown mount type")
	}

	response.MountType = "kv"

	return response, nil
}
