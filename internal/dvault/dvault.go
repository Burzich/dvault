package dvault

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	kv2 "github.com/Burzich/dvault/internal/dvault/kv"
	"github.com/Burzich/dvault/internal/dvault/kv/disc"
	"github.com/Burzich/dvault/internal/dvault/tools"
	"github.com/cloudflare/circl/group"
	"github.com/cloudflare/circl/secretsharing"
)

type DVault struct {
	EncryptionKey string
	Logger        *slog.Logger
	mountPath     string

	buildDate     time.Time
	isSealed      atomic.Bool
	isInitialized atomic.Bool

	mu sync.RWMutex

	kv         map[string]kv2.KV
	shareKeys  []string
	N          int
	T          int
	commitment secretsharing.SecretCommitment
}

func NewDVault(logger *slog.Logger, mountPath string) *DVault {
	d := DVault{
		Logger:        logger,
		mountPath:     mountPath,
		buildDate:     time.Now(),
		isSealed:      atomic.Bool{},
		isInitialized: atomic.Bool{},
		mu:            sync.RWMutex{},
		kv:            make(map[string]kv2.KV),
		shareKeys:     nil,
		N:             0,
		T:             0,
		commitment:    nil,
	}

	d.isSealed.Store(true)

	return &d
}

func (d *DVault) Unseal(ctx context.Context, unseal Unseal) (UnsealResponse, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isInitialized.Load() {
		return UnsealResponse{}, errors.New("already unsealed")
	}

	if unseal.Reset {
		d.shareKeys = make([]string, 0)
	}

	d.shareKeys = append(d.shareKeys, unseal.Key)

	if len(d.shareKeys) == d.T {
		err := d.tryUnseal(d.shareKeys)
		if err != nil {
			return UnsealResponse{}, err
		}

		d.isSealed.Store(false)

		return UnsealResponse{
			BuildDate:         d.buildDate.String(),
			ClusterId:         "dvault",
			ClusterName:       "dvault",
			HcpLinkResourceID: "",
			HcpLinkStatus:     "",
			Initialized:       d.isInitialized.Load(),
			Migration:         false,
			N:                 d.N,
			T:                 d.T,
			Progress:          0,
			Nonce:             "",
			RecoverySeal:      false,
			Sealed:            d.isSealed.Load(),
			StorageType:       "file",
			Type:              "shamir",
			Version:           "1.0.0",
		}, nil
	}

	return UnsealResponse{
		BuildDate:         d.buildDate.String(),
		ClusterId:         "dvault",
		ClusterName:       "dvault",
		HcpLinkResourceID: "",
		HcpLinkStatus:     "",
		Initialized:       d.isInitialized.Load(),
		Migration:         false,
		N:                 d.N,
		T:                 d.T,
		Progress:          len(d.shareKeys),
		Nonce:             "",
		RecoverySeal:      false,
		Sealed:            d.isSealed.Load(),
		StorageType:       "file",
		Type:              "shamir",
		Version:           "1.0.0",
	}, nil
}

func (d *DVault) Seal(ctx context.Context) (Response, error) {
	var response Response
	response.RequestId = tools.GenerateXRequestID()

	d.isSealed.Store(true)

	return response, nil
}

func (d *DVault) Init(_ context.Context, init Init) (InitResponse, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	g := group.P256
	t := uint(2)
	n := uint(5)

	if init.SecretShares != 0 {
		n = uint(init.SecretShares)
	}
	if init.SecretThreshold != 0 {
		t = uint(init.SecretThreshold)
	}

	secret := g.RandomScalar(rand.Reader)
	ss := secretsharing.New(rand.Reader, t, secret)
	shares := ss.Share(n)

	com := ss.CommitSecret()

	var sharesValuesBase64 []string

	for _, share := range shares {
		shareValueBytes, err := share.Value.MarshalBinary()
		if err != nil {
			return InitResponse{}, err
		}
		shareIdBytes, err := share.ID.MarshalBinary()
		if err != nil {
			return InitResponse{}, err
		}

		shareValueBase64 := base64.StdEncoding.EncodeToString(shareValueBytes)
		shareIdBase64 := base64.StdEncoding.EncodeToString(shareIdBytes)
		sharesValuesBase64 = append(sharesValuesBase64, shareValueBase64+"#"+shareIdBase64)
	}

	secretBytes, err := secret.MarshalBinary()
	if err != nil {
		return InitResponse{}, err
	}

	err = d.generateAndSaveEncryptKey(secretBytes, n, t)
	if err != nil {
		return InitResponse{}, err
	}
	d.commitment = com

	return InitResponse{
		Keys:       sharesValuesBase64,
		KeysBase64: sharesValuesBase64,
		RootToken:  base64.StdEncoding.EncodeToString(secretBytes),
	}, nil
}

func (d *DVault) SealStatus(ctx context.Context) (SealStatus, error) {
	return SealStatus{
		Type:         "shamir",
		Initialized:  d.isInitialized.Load(),
		Sealed:       d.isSealed.Load(),
		T:            d.T,
		N:            d.N,
		Progress:     len(d.shareKeys),
		Nonce:        "",
		Version:      "1.0.0",
		BuildDate:    d.buildDate,
		Migration:    false,
		ClusterName:  "dvault",
		ClusterId:    "dvault",
		RecoverySeal: false,
		StorageType:  "file",
	}, nil
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

func (d *DVault) generateAndSaveEncryptKey(secret []byte, shares uint, threshold uint) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	encryptKey := make([]byte, 256)
	_, err := rand.Read(encryptKey)
	if err != nil {
		return err
	}

	encryptedEncryptedKey, err := tools.Encrypt(encryptKey, secret)
	if err != nil {
		return err
	}

	encryptedEncryptedKeyBase64 := base64.StdEncoding.EncodeToString(encryptedEncryptedKey)
	keyPath := filepath.Join(d.mountPath, "key")
	err = os.WriteFile(keyPath, []byte(encryptedEncryptedKeyBase64), 0600)

	if err != nil {
		return err
	}

	return nil
}

func (d *DVault) tryUnseal(keysBase64Encoded []string) error {
	valueKeys := make([][]byte, len(keysBase64Encoded))
	idKeys := make([][]byte, len(keysBase64Encoded))
	for i := range valueKeys {
		valueBase64, idBase64, ok := strings.Cut(keysBase64Encoded[i], "#")
		if !ok {
			return errors.New("invalid share")
		}

		{
			base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(valueBase64)))

			n, err := base64.StdEncoding.Decode(base64Text, []byte(valueBase64))
			if err != nil {
				return err
			}

			valueKeys[i] = base64Text[:n]
		}
		{
			base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(idBase64)))

			n, err := base64.StdEncoding.Decode(base64Text, []byte(idBase64))
			if err != nil {
				return err
			}

			idKeys[i] = base64Text[:n]
		}
	}

	var values []group.Scalar
	for i := range valueKeys {
		g := group.P256
		scalar := g.NewScalar()
		err := scalar.UnmarshalBinary(valueKeys[i])
		if err != nil {
			return err
		}
		values = append(values, scalar)
	}

	var ids []group.Scalar
	for i := range idKeys {
		g := group.P256
		scalar := g.NewScalar()
		err := scalar.UnmarshalBinary(idKeys[i])
		if err != nil {
			return err
		}
		ids = append(values, scalar)
	}

	var shares []secretsharing.Share
	for i := range valueKeys {
		shares = append(shares, secretsharing.Share{
			ID:    ids[i],
			Value: values[i],
		})
	}

	for i := range shares {
		ok := secretsharing.Verify(uint(d.T), shares[i], d.commitment)
		if !ok {
			return errors.New("invalid share")
		}
	}

	return nil
}
