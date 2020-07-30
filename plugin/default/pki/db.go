package pki

import (
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/jmoiron/sqlx"
)

type dbStorage struct {
	db *sqlx.DB
}

func NewStorageDatabase(cfg Persistent) (Storage, error) {
	db, err := sqlx.Open(cfg.Database.Type, cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &dbStorage{
		db: db,
	}, nil
}

func (d dbStorage) CreateCert(cert plugin.Cert) error {
	insertSQL := `
INSERT INTO baetyl_certificate (
cert_id, parent_id, type, common_name, 
description, csr, content, private_key, not_before, not_after) 
VALUES (?,?,?,?,?,?,?,?,?,?)
`
	_, err := d.db.Exec(insertSQL,
		cert.CertId, cert.ParentId, cert.Type,
		cert.CommonName, cert.Description, cert.Csr,
		cert.Content, cert.PrivateKey, cert.NotBefore, cert.NotAfter)
	return err
}

func (d dbStorage) DeleteCert(certId string) error {
	deleteSQL := `
DELETE FROM baetyl_certificate where cert_id=?
`
	_, err := d.db.Exec(deleteSQL, certId)
	return err
}

func (d dbStorage) UpdateCert(cert plugin.Cert) error {
	updateSQL := `
UPDATE baetyl_certificate SET parent_id=?,type=?,
common_name=?,description=?,csr=?,content=?,private_key=?,
not_before=?, not_after=? 
WHERE cert_id=?
`
	_, err := d.db.Exec(updateSQL,
		cert.ParentId, cert.Type, cert.CommonName, cert.Description, cert.Csr,
		cert.Content, cert.PrivateKey, cert.NotBefore, cert.NotAfter, cert.CertId)
	return err
}

func (d dbStorage) GetCert(certId string) (*plugin.Cert, error) {
	selectSQL := `
SELECT cert_id, parent_id, type, common_name, 
description, csr, content, private_key, not_before, not_after
FROM baetyl_certificate 
WHERE cert_id=? LIMIT 0,1
`
	var cert []plugin.Cert
	if err := d.db.Select(&cert, selectSQL, certId); err != nil {
		return nil, err
	}
	if len(cert) > 0 {
		return &cert[0], nil
	}
	return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "certificate"), common.Field("name", certId))
}

func (d dbStorage) CountCertByParentId(parentId string) (int, error) {
	selectSQL := `
SELECT count(cert_id) AS count 
FROM baetyl_certificate 
WHERE parent_id=?
`
	var res []struct {
		Count int `db:"count"`
	}
	if err := d.db.Select(&res, selectSQL, parentId); err != nil {
		return 0, err
	}
	return res[0].Count, nil
}

func (d dbStorage) Close() error {
	return d.db.Close()
}
