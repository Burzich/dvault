package dvault

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
