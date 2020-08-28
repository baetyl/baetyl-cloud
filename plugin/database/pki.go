package database

import (
	"os"

	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

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
	return nil, os.ErrNotExist
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
