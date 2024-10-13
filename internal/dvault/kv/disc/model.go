package disc

import "github.com/Burzich/dvault/internal/dvault/kv"

type Data struct {
	Records []kv.Record `json:"records"`
	Meta    kv.Meta     `json:"meta"`
}
