package dvault

import "time"

type Response struct {
	RequestId     string      `json:"request_id"`
	LeaseId       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          interface{} `json:"data"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      interface{} `json:"warnings"`
	Auth          interface{} `json:"auth"`
	MountType     string      `json:"mount_type"`
}

type InitResponse struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base64"`
	RootToken  string   `json:"root_token"`
}

type Unseal struct {
	Key     string `json:"key"`
	Migrate bool   `json:"migrate"`
	Reset   bool   `json:"reset"`
}

type Init struct {
	PgpKeys           []string `json:"pgp_keys"`
	RecoveryPgpKeys   []string `json:"recovery_pgp_keys"`
	RecoveryShares    int      `json:"recovery_shares"`
	RecoveryThreshold int      `json:"recovery_threshold"`
	RootTokenPgpKey   string   `json:"root_token_pgp_key"`
	SecretShares      int      `json:"secret_shares"`
	SecretThreshold   int      `json:"secret_threshold"`
	StoredShares      int      `json:"stored_shares"`
}

type SealStatus struct {
	Type         string    `json:"type"`
	Initialized  bool      `json:"isInitialized"`
	Sealed       bool      `json:"sealed"`
	T            int       `json:"t"`
	N            int       `json:"n"`
	Progress     int       `json:"progress"`
	Nonce        string    `json:"nonce"`
	Version      string    `json:"version"`
	BuildDate    time.Time `json:"build_date"`
	Migration    bool      `json:"migration"`
	ClusterName  string    `json:"cluster_name"`
	ClusterId    string    `json:"cluster_id"`
	RecoverySeal bool      `json:"recovery_seal"`
	StorageType  string    `json:"storage_type"`
}

type UnsealResponse struct {
	BuildDate         string `json:"build_date"`
	ClusterId         string `json:"cluster_id"`
	ClusterName       string `json:"cluster_name"`
	HcpLinkResourceID string `json:"hcp_link_resource_ID"`
	HcpLinkStatus     string `json:"hcp_link_status"`
	Initialized       bool   `json:"isInitialized"`
	Migration         bool   `json:"migration"`
	N                 int    `json:"n"`
	Nonce             string `json:"nonce"`
	Progress          int    `json:"progress"`
	RecoverySeal      bool   `json:"recovery_seal"`
	Sealed            bool   `json:"sealed"`
	StorageType       string `json:"storage_type"`
	T                 int    `json:"t"`
	Type              string `json:"type"`
	Version           string `json:"version"`
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

type Mounts struct {
	RequestId     string               `json:"request_id"`
	LeaseId       string               `json:"lease_id"`
	Renewable     bool                 `json:"renewable"`
	LeaseDuration int                  `json:"lease_duration"`
	Data          map[string]MountData `json:"data"`
	WrapInfo      interface{}          `json:"wrap_info"`
	Warnings      interface{}          `json:"warnings"`
	Auth          interface{}          `json:"auth"`
	MountType     string               `json:"mount_type"`
}

type MountData struct {
	Accessor string `json:"accessor"`
	Config   struct {
		DefaultLeaseTtl   int    `json:"default_lease_ttl"`
		ForceNoCache      bool   `json:"force_no_cache"`
		ListingVisibility string `json:"listing_visibility"`
		MaxLeaseTtl       int    `json:"max_lease_ttl"`
	} `json:"config"`
	DeprecationStatus     string `json:"deprecation_status"`
	Description           string `json:"description"`
	ExternalEntropyAccess bool   `json:"external_entropy_access"`
	Local                 bool   `json:"local"`
	Options               struct {
		Version string `json:"version"`
	} `json:"options"`
	PluginVersion        string `json:"plugin_version"`
	RunningPluginVersion string `json:"running_plugin_version"`
	RunningSha256        string `json:"running_sha256"`
	SealWrap             bool   `json:"seal_wrap"`
	Type                 string `json:"type"`
	Uuid                 string `json:"uuid"`
}
