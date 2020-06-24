package pki

import (
	"encoding/base64"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-go/log"
	"time"
)

type Cert struct {
	CertId      string    `db:"cert_id"`
	ParentId    string    `db:"parent_id"`
	Type        string    `db:"type"`
	CommonName  string    `db:"common_name"`
	Csr         string    `db:"csr"`
	Content     string    `db:"content"`
	PrivateKey  string    `db:"private_key"`
	Description string    `db:"description"`
	NotBefore   time.Time `db:"not_before"`
	NotAfter    time.Time `db:"not_after"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func ToCertModel(c *Cert) *models.Cert {
	csr, err := base64.StdEncoding.DecodeString(c.Csr)
	if err != nil {
		log.L().Error("ToCertModel csr", log.Error(err))
		return nil
	}
	content, err := base64.StdEncoding.DecodeString(c.Content)
	if err != nil {
		log.L().Error("ToCertModel content", log.Error(err))
		return nil
	}
	priv, err := base64.StdEncoding.DecodeString(c.PrivateKey)
	if err != nil {
		log.L().Error("ToCertModel priv", log.Error(err))
		return nil
	}
	return &models.Cert{
		CertId:      c.CertId,
		ParentId:    c.ParentId,
		Type:        c.Type,
		CommonName:  c.CommonName,
		Csr:         csr,
		Content:     content,
		Priv:        priv,
		Description: c.Description,
		NotBefore:   c.NotBefore,
		NotAfter:    c.NotAfter,
	}
}

func FromCertModel(c *models.Cert) *Cert {
	return &Cert{
		CertId:      c.CertId,
		ParentId:    c.ParentId,
		Type:        c.Type,
		CommonName:  c.CommonName,
		Csr:         base64.StdEncoding.EncodeToString(c.Csr),
		Content:     base64.StdEncoding.EncodeToString(c.Content),
		PrivateKey:  base64.StdEncoding.EncodeToString(c.Priv),
		Description: c.Description,
		NotBefore:   c.NotBefore,
		NotAfter:    c.NotAfter,
	}
}
