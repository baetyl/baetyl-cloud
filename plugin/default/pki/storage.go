package pki

import (
	"io"
	"os"

	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../../../mock/plugin/default/pki/storage.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin/default/pki Storage

const (
	File       = "file"
	Database   = "database"
	Kubernetes = "kubernetes"
)

type Storage interface {
	CreateCert(cert plugin.Cert) error
	DeleteCert(certId string) error
	UpdateCert(cert plugin.Cert) error
	GetCert(certId string) (*plugin.Cert, error)
	CountCertByParentId(parentId string) (int, error)
	io.Closer
}

func NewStorage(cfg Persistent) (Storage, error) {
	switch cfg.Kind {
	case Database:
		return NewStorageDatabase(cfg)
	default:
		return nil, os.ErrInvalid
	}
}
