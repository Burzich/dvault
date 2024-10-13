package kv

import "time"

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
