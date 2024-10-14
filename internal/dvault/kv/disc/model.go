package disc

import (
	"github.com/Burzich/dvault/internal/dvault/kv"
)

type Data struct {
	Records []kv.KVRecord `json:"records"`
	Meta    kv.KVMeta     `json:"meta"`
}
