package pki

import (
	"github.com/baetyl/baetyl-cloud/models"
	"io"
	"os"
)

//go:generate mockgen -destination=../../../mock/plugin/default/pvc.go -package=plugin github.com/baetyl/baetyl-cloud/plugin/default/pki PVC

const (
	File       = "file"
	Database   = "database"
	Kubernetes = "kubernetes"
)

type PVC interface {
	CreateCert(cert models.Cert) error
	DeleteCert(certId string) error
	UpdateCert(cert models.Cert) error
	GetCert(certId string) (*models.Cert, error)
	CountCertByParentId(parentId string) (int, error)
	io.Closer
}

func NewPVC(cfg Persistent) (PVC, error) {
	switch cfg.Kind {
	case Database:
		return NewPVCDatabase(cfg)
	default:
		return nil, os.ErrInvalid
	}
}
