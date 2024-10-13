package dvault

import "log/slog"

type DVault struct {
	EncryptionKey string
	Logger        *slog.Logger
}

func NewDVault(logger *slog.Logger) *DVault {
	return &DVault{
		EncryptionKey: "",
		Logger:        logger,
	}
}

func (d DVault) Unseal() {

}

func (d DVault) Seal() {

}
