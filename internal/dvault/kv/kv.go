package kv

import (
	"context"
	"errors"
	"time"
)

type Record struct {
	Data     map[string]interface{} `json:"data"`
	Metadata struct {
		CreatedTime    time.Time   `json:"created_time"`
		CustomMetadata interface{} `json:"custom_metadata"`
		DeletionTime   string      `json:"deletion_time"`
		Destroyed      bool        `json:"destroyed"`
		Version        int         `json:"version"`
	} `json:"metadata"`
}

type Config struct {
	CasRequired        bool   `json:"cas_required"`
	DeleteVersionAfter string `json:"delete_version_after"`
	MaxVersions        int    `json:"max_versions"`
}

type Meta struct {
	CasRequired        bool                   `json:"cas_required"`
	CreatedTime        time.Time              `json:"created_time"`
	CurrentVersion     int                    `json:"current_version"`
	DeleteVersionAfter string                 `json:"delete_version_after"`
	MaxVersions        int                    `json:"max_versions"`
	OldestVersion      int                    `json:"oldest_version"`
	UpdatedTime        time.Time              `json:"updated_time"`
	CustomMetadata     map[string]interface{} `json:"custom_metadata"`
	Versions           map[string]struct {
		CreatedTime  time.Time `json:"created_time"`
		DeletionTime string    `json:"deletion_time"`
		Destroyed    bool      `json:"destroyed"`
	} `json:"versions"`
}

type KV interface {
	Get(ctx context.Context, secretPath string) (Record, error)
	GetVersion(ctx context.Context, secretPath string, version int) (Record, error)
	Save(ctx context.Context, secretPath string, data map[string]interface{}, cas int) error
	Delete(ctx context.Context, secretPath string) error
	Undelete(ctx context.Context, secretPath string) error
	DeleteVersion(ctx context.Context, secretPath string, version []int) error
	UndeleteVersion(ctx context.Context, secretPath string, version int) error
	Destroy(ctx context.Context, secretPath string, version []int) error

	UpdateConfig(ctx context.Context, config Config) error
	GetConfig(ctx context.Context) (Config, error)
	GetMeta(ctx context.Context, secretPath string) (Meta, error)
	UpdateMeta(ctx context.Context, secretPath string, meta Meta) error
	DeleteMeta(ctx context.Context, secretPath string) error
}

func CreateConfigFromMap(m map[string]interface{}) (Config, error) {
	casRequired := false
	deleteVersionAfter := ""
	maxVersions := 0

	if _, ok := m["casRequired"]; ok {
		v, ok := m["casRequired"].(bool)
		if !ok {
			return Config{}, errors.New("casRequired must be a bool")
		}
		casRequired = v
	}

	if _, ok := m["deleteVersionAfter"]; ok {
		v, ok := m["deleteVersionAfter"].(string)
		if !ok {
			return Config{}, errors.New("deleteVersionAfter must be a string")
		}
		deleteVersionAfter = v
	}

	if _, ok := m["maxVersions"]; ok {
		v, ok := m["maxVersions"].(int)
		if !ok {
			return Config{}, errors.New("maxVersions must be an integer")
		}
		maxVersions = v
	}

	return Config{
		CasRequired:        casRequired,
		DeleteVersionAfter: deleteVersionAfter,
		MaxVersions:        maxVersions,
	}, nil
}
