package audit

import (
	"fmt"
	"io"
)

type AuditJournal struct {
	io io.Writer
}

func (j AuditJournal) Write() error {
	fmt.Fprint(j.io, "hello")
	return nil
}
