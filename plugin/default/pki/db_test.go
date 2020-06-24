package pki

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	csr       = "csr"
	content   = "content"
	priv      = "priv"
	timestamp = time.Unix(1000, 1000)

	certificateTables = []string{
		`
CREATE TABLE baetyl_certificate
(
    cert_id          varchar(128)  PRIMARY KEY,
    parent_id        varchar(128)  NOT NULL DEFAULT '',
    type             varchar(64)   NOT NULL DEFAULT '',
    common_name      varchar(128)  NOT NULL DEFAULT '',
    description      varchar(256)  NOT NULL DEFAULT '',
    csr              varchar(2048) DEFAULT '',
    content          varchar(2048) DEFAULT '',
    private_key      varchar(2048) DEFAULT '',
    not_before       timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    not_after        timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    create_time      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time      timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func mockNewDB(cfg Persistent) (*dbStorage, error) {
	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &dbStorage{db: db}, nil
}

func genCertificate() *models.Cert {
	return &models.Cert{
		CertId:      "123",
		ParentId:    "456",
		Type:        TypeIssuingCA,
		CommonName:  "cn",
		Csr:         []byte(csr),
		Content:     []byte(content),
		Priv:        []byte(priv),
		Description: "desc",
		NotBefore:   timestamp,
		NotAfter:    timestamp,
	}
}

func (d *dbStorage) MockCreateCertificateTable() {
	for _, sql := range certificateTables {
		_, err := d.db.Exec(sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestCertificate(t *testing.T) {
	var cfg Persistent
	cfg.Database.Type = "sqlite3"
	cfg.Database.URL = ":memory:"

	db, err := mockNewDB(cfg)
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}

	certificate := genCertificate()
	db.MockCreateCertificateTable()
	err = db.CreateCert(*certificate)
	assert.NoError(t, err)

	resCertificate, err := db.GetCert(certificate.CertId)
	assert.NoError(t, err)
	checkCertificate(t, certificate, resCertificate)

	certificate.ParentId = "new"
	err = db.UpdateCert(*certificate)
	assert.NoError(t, err)
	resCertificate, err = db.GetCert(certificate.CertId)
	assert.NoError(t, err)
	checkCertificate(t, certificate, resCertificate)

	c1, err := db.CountCertByParentId(certificate.ParentId)
	assert.NoError(t, err)
	assert.Equal(t, 1, c1)

	err = db.DeleteCert(certificate.CertId)
	assert.NoError(t, err)

	_, err = db.GetCert(certificate.CertId)
	assert.NoError(t, err)

	// test new
	_, err = NewPVCDatabase(cfg)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func checkCertificate(t *testing.T, expect, actual *models.Cert) {
	assert.Equal(t, expect.CertId, actual.CertId)
	assert.Equal(t, expect.ParentId, actual.ParentId)
	assert.Equal(t, expect.Description, actual.Description)
	assert.Equal(t, expect.Type, actual.Type)
	assert.Equal(t, expect.CommonName, actual.CommonName)
	assert.Equal(t, expect.Content, actual.Content)
	assert.Equal(t, expect.Csr, actual.Csr)
	assert.Equal(t, expect.Priv, actual.Priv)
}
