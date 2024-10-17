package handler

type UnsealRequest struct {
	Key     string `json:"key"`
	Migrate bool   `json:"migrate"`
	Reset   bool   `json:"reset"`
}

type InitRequest struct {
	PgpKeys           []string `json:"pgp_keys"`
	RecoveryPgpKeys   []string `json:"recovery_pgp_keys"`
	RecoveryShares    int      `json:"recovery_shares"`
	RecoveryThreshold int      `json:"recovery_threshold"`
	RootTokenPgpKey   string   `json:"root_token_pgp_key"`
	SecretShares      int      `json:"secret_shares"`
	SecretThreshold   int      `json:"secret_threshold"`
	StoredShares      int      `json:"stored_shares"`
}

type UpdateKVConfigRequest struct {
	CasRequired        bool   `json:"cas_required"`
	DeleteVersionAfter string `json:"delete_version_after"`
	MaxVersions        int    `json:"max_versions"`
}

type CreateKVSecretRequest struct {
	Data    map[string]interface{} `json:"data"`
	Options struct {
		CAS int `json:"cas"`
	} `json:"options"`
}

type UpdateKVMetadata struct {
	CasRequired        bool                   `json:"cas_required"`
	CustomMetadata     map[string]interface{} `json:"custom_metadata"`
	DeleteVersionAfter string                 `json:"delete_version_after"`
	MaxVersions        int                    `json:"max_versions"`
}

type DeleteKVSecret struct {
	Versions []int `json:"versions"`
}

type DestroyKVSecret struct {
	Versions []int `json:"versions"`
}

type CreateMount struct {
	Config                map[string]interface{} `json:"config"`
	Description           string                 `json:"description"`
	ExternalEntropyAccess bool                   `json:"external_entropy_access"`
	Local                 bool                   `json:"local"`
	Options               map[string]interface{} `json:"options"`
	PluginName            string                 `json:"plugin_name"`
	PluginVersion         string                 `json:"plugin_version"`
	SealWrap              bool                   `json:"seal_wrap"`
	Type                  string                 `json:"type"`
}
