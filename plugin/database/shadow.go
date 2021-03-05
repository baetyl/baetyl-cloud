package database

import (
	"database/sql"
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-go/v2/errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

const batchSize = 200

func (d *DB) Get(namespace, name string) (*models.Shadow, error) {
	return d.GetShadowTx(nil, namespace, name)
}

func (d *DB) Create(shadow *models.Shadow) (*models.Shadow, error) {
	var shd *models.Shadow
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateShadowTx(tx, shadow)
		if err != nil {
			return err
		}
		shd, err = d.GetShadowTx(tx, shadow.Namespace, shadow.Name)
		return err
	})

	return shd, err
}

func (d *DB) List(namespace string, nodeList *models.NodeList) (*models.ShadowList, error) {
	names := make([]string, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		names = append(names, node.Name)
	}

	shadows, err := d.ListShadowByNamesTx(nil, namespace, names)
	if err != nil {
		return nil, err
	}

	total := len(shadows)
	items := make([]models.Shadow, 0, total)

	for _, s := range shadows {
		shd, _ := s.ToShadowModel()
		items = append(items, *shd)
	}

	result := &models.ShadowList{
		Total: total,
		Items: items,
	}
	return result, nil
}

func (d *DB) Delete(namespace, name string) error {
	_, err := d.DeleteShadowTx(nil, namespace, name)
	return err
}

func (d *DB) UpdateDesire(shadow *models.Shadow) (*models.Shadow, error) {
	var shd *models.Shadow
	err := d.Transact(func(tx *sqlx.Tx) error {
		res, err := d.UpdateShadowDesireTx(tx, shadow)
		if err != nil {
			return err
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			return errors.New(common.ErrUpdateCas)
		}
		shd, err = d.GetShadowTx(tx, shadow.Namespace, shadow.Name)
		return err
	})

	return shd, err
}
func (d *DB) UpdateReport(shadow *models.Shadow) (*models.Shadow, error) {
	var shd *models.Shadow
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.UpdateShadowReportTx(tx, shadow)
		if err != nil {
			return err
		}
		shd, err = d.GetShadowTx(tx, shadow.Namespace, shadow.Name)
		return err
	})

	return shd, err
}

func (d *DB) GetShadowTx(tx *sqlx.Tx, namespace, name string) (*models.Shadow, error) {
	selectSQL := `
SELECT 
id, name, namespace, report, desire, create_time, update_time, desire_version 
FROM baetyl_node_shadow WHERE namespace=? AND name=?
`
	var shadows []entities.Shadow
	if err := d.Query(tx, selectSQL, &shadows, namespace, name); err != nil {
		return nil, err
	}
	if len(shadows) > 0 {
		return shadows[0].ToShadowModel()
	}
	return nil, nil
}

func (d *DB) CreateShadowTx(tx *sqlx.Tx, shadow *models.Shadow) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_node_shadow (namespace, name, report, desire)
VALUES (?, ?, ?, ?)
`

	shd, err := entities.NewShadowFromShadowModel(shadow)
	if err != nil {
		return nil, err
	}

	return d.Exec(tx, insertSQL, shd.Namespace, shd.Name, shd.Report, shd.Desire)
}

func (d *DB) DeleteShadowTx(tx *sqlx.Tx, namespace, name string) (sql.Result, error) {
	deleteSql := `
DELETE FROM baetyl_node_shadow WHERE namespace=? AND name=?
`
	return d.Exec(tx, deleteSql, namespace, name)
}

func (d *DB) UpdateShadowDesireTx(tx *sqlx.Tx, shadow *models.Shadow) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_node_shadow
SET desire=?, desire_version=?
WHERE namespace=? AND name=? AND desire_version=?
`
	desire, err := shadow.GetDesireString()
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, desire, genResourceVersion(), shadow.Namespace, shadow.Name, shadow.DesireVersion)
}

func (d *DB) UpdateShadowReportTx(tx *sqlx.Tx, shadow *models.Shadow) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_node_shadow
SET report=?
WHERE namespace=? AND name=?
`
	report, err := shadow.GetReportString()
	if err != nil {
		return nil, err
	}

	return d.Exec(tx, updateSQL, report, shadow.Namespace, shadow.Name)
}

func (d *DB) ListShadowByNamesTx(tx *sqlx.Tx, namespace string, names []string) ([]entities.Shadow, error) {
	selectSQL := `
SELECT 
id, name, namespace, report, desire, create_time, update_time, desire_version
FROM baetyl_node_shadow WHERE namespace=? AND name in (?)
`
	result := make([]entities.Shadow, 0)
	length := len(names)

	for start, end := 0, batchSize; start < length; start, end = end, end+batchSize {

		if end > length {
			end = length
		}
		var shadows []entities.Shadow

		sql, args, err := sqlx.In(selectSQL, namespace, names[start:end])
		if err != nil {
			return nil, err
		}

		if err := d.Query(tx, sql, &shadows, args...); err != nil {
			return nil, err
		}
		result = append(result, shadows...)
	}

	return result, nil
}

func genResourceVersion() string {
	return fmt.Sprintf("%d%s", time.Now().UTC().Unix(), common.RandString(6))
}
