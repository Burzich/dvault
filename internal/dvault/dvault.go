package dvault

import (
	"context"
	"log/slog"

	"github.com/Burzich/dvault/internal/dvault/kv/disc"
)

type DVault struct {
	EncryptionKey string
	Logger        *slog.Logger

	diskKV disc.KV
}

func NewDVault(logger *slog.Logger) *DVault {
	return &DVault{
		EncryptionKey: "",
		Logger:        logger,
	}
}

func (d *DVault) Unseal(ctx context.Context) error {
	return nil
}

func (d *DVault) Seal(ctx context.Context) error {
	return nil
}

func (d *DVault) SealStatus(ctx context.Context) error {
	return nil
}
