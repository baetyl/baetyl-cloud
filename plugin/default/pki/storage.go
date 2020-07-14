package pki

import (
	"os"

	"github.com/baetyl/baetyl-go/v2/pki"
)

const (
	File       = "file"
	Database   = "database"
	Kubernetes = "kubernetes"
)

func NewStorage(cfg Persistent) (pki.Storage, error) {
	switch cfg.Kind {
	case Database:
		return NewStorageDatabase(cfg)
	default:
		return nil, os.ErrInvalid
	}
}
