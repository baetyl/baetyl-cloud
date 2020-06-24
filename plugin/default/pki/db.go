package pki

import (
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type dbStorage struct {
	db *sqlx.DB
}

func NewPVCDatabase(cfg Persistent) (PVC, error) {
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

func (d dbStorage) CreateCert(cert models.Cert) error {
	insertSQL := `
INSERT INTO baetyl_certificate (
cert_id, parent_id, type, common_name, 
description, csr, content, private_key, not_before, not_after) 
VALUES (?,?,?,?,?,?,?,?,?,?)
`
	c := FromCertModel(&cert)
	_, err := d.db.Exec(insertSQL,
		c.CertId, c.ParentId, c.Type,
		c.CommonName, c.Description, c.Csr,
		c.Content, c.PrivateKey, c.NotBefore, c.NotAfter)
	return err
}

func (d dbStorage) DeleteCert(certId string) error {
	deleteSQL := `
DELETE FROM baetyl_certificate where cert_id=?
`
	_, err := d.db.Exec(deleteSQL, certId)
	return err
}

func (d dbStorage) UpdateCert(cert models.Cert) error {
	updateSQL := `
UPDATE baetyl_certificate SET parent_id=?,type=?,
common_name=?,description=?,csr=?,content=?,private_key=?,
not_before=?, not_after=? 
WHERE cert_id=?
`
	c := FromCertModel(&cert)
	_, err := d.db.Exec(updateSQL,
		c.ParentId, c.Type, c.CommonName, c.Description, c.Csr,
		c.Content, c.PrivateKey, c.NotBefore, c.NotAfter, c.CertId)
	return err
}

func (d dbStorage) GetCert(certId string) (*models.Cert, error) {
	selectSQL := `
SELECT cert_id, parent_id, type, common_name, 
description, csr, content, private_key, not_before, not_after,
create_time, update_time 
FROM baetyl_certificate 
WHERE cert_id=? LIMIT 0,1
`
	var cert []Cert
	if err := d.db.Select(&cert, selectSQL, certId); err != nil {
		return nil, err
	}
	if len(cert) > 0 {
		return ToCertModel(&cert[0]), nil
	}
	return nil, nil
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
