package handler

type UnsealRequest struct {
	Key     string `json:"key"`
	Migrate bool   `json:"migrate"`
	Reset   bool   `json:"reset"`
}
