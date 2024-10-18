package handler

type UnsealResponse struct {
	BuildDate         string `json:"build_date"`
	ClusterId         string `json:"cluster_id"`
	ClusterName       string `json:"cluster_name"`
	HcpLinkResourceID string `json:"hcp_link_resource_ID"`
	HcpLinkStatus     string `json:"hcp_link_status"`
	Initialized       bool   `json:"initialized"`
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

type SealStatusResponse struct {
	BuildDate         string `json:"build_date"`
	ClusterId         string `json:"cluster_id"`
	ClusterName       string `json:"cluster_name"`
	HcpLinkResourceID string `json:"hcp_link_resource_ID"`
	HcpLinkStatus     string `json:"hcp_link_status"`
	Initialized       bool   `json:"initialized"`
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

type Error struct {
	Errors []string `json:"errors"`
}
