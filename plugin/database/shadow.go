package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/trigger"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
	"github.com/baetyl/baetyl-cloud/v2/triggerfunc"
)

const batchSize = 200

func (d *DB) Get(tx interface{}, namespace, name string) (*models.Shadow, error) {
	var transaction *sqlx.Tx
	if tx != nil {
		transaction = tx.(*sqlx.Tx)
	}
	return d.GetShadowTx(transaction, namespace, name)
}

func (d *DB) Create(tx interface{}, shadow *models.Shadow) (*models.Shadow, error) {
	var shd *models.Shadow
	var err error
	if tx == nil {
		err = d.Transact(func(transaction *sqlx.Tx) error {
			return d.createAndGetShadow(transaction, shadow, &shd)
		})
	} else {
		transaction := tx.(*sqlx.Tx)
		err = d.createAndGetShadow(transaction, shadow, &shd)
		if err != nil {
			return shd, err
		}
	}
	//exec triggerfunc trigger.go  ShadowCreateOrUpdateCacheSet
	_, err = trigger.Exec(triggerfunc.ShadowCreateOrUpdateTrigger, *shd)
	return shd, err
}

func (d *DB) createAndGetShadow(tx *sqlx.Tx, input *models.Shadow, output **models.Shadow) error {
	_, err := d.CreateShadowTx(tx, input)
	if err != nil {
		return err
	}
	*output, err = d.GetShadowTx(tx, input.Namespace, input.Name)
	return err
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

func (d *DB) ListShadowByNames(tx interface{}, namespace string, names []string) ([]*models.Shadow, error) {
	if names == nil || len(names) < 1 {
		return nil, nil
	}
	var transaction *sqlx.Tx
	if tx != nil {
		transaction = tx.(*sqlx.Tx)
	}
	return d.listShadowByNamesTx(transaction, namespace, names)
}

func (d *DB) Delete(namespace, name string) error {
	_, err := d.DeleteShadowTx(nil, namespace, name)

	if err != nil {
		return err
	}
	//exec common trigger.go  ShadowDelete
	_, err = trigger.Exec(triggerfunc.ShadowDelete, name)
	return err
}

func (d *DB) UpdateDesire(tx interface{}, shadow *models.Shadow) error {
	var transaction *sqlx.Tx
	if tx != nil {
		transaction = tx.(*sqlx.Tx)
	}
	_, err := d.UpdateShadowDesireTx(transaction, shadow)
	return err
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
	//exec common trigger.go  ShadowCreateOrUpdateCacheSet
	_, err = trigger.Exec(triggerfunc.ShadowCreateOrUpdateTrigger, *shd)
	return shd, err
}

func (d *DB) UpdateDesires(tx interface{}, shadows []*models.Shadow) error {
	if shadows == nil || len(shadows) < 1 {
		return nil
	}
	var err error
	if tx == nil {
		err = d.Transact(func(transaction *sqlx.Tx) error {
			_, updateErr := d.updateDesiresTx(transaction, shadows)
			return updateErr
		})
	} else {
		transaction := tx.(*sqlx.Tx)
		_, err = d.updateDesiresTx(transaction, shadows)
	}
	return err
}

func (d *DB) updateDesiresTx(tx *sqlx.Tx, shadows []*models.Shadow) (sql.Result, error) {
	insertSQL := `INSERT INTO baetyl_node_shadow(name,namespace,desire,desire_version,desire_meta) VALUES `
	params := []interface{}{}
	for _, shadow := range shadows {
		insertSQL += `(?,?,?,?,?),`
		desire, err := shadow.GetDesireString()
		if err != nil {
			return nil, err
		}
		desireMeta, err := shadow.GetDesireMetaString()
		if err != nil {
			return nil, err
		}
		params = append(params, shadow.Name, shadow.Namespace, desire, genResourceVersion(), desireMeta)
	}
	insertSQL = strings.TrimRight(insertSQL, ",")
	insertSQL += `
ON DUPLICATE KEY UPDATE 
name=VALUES(name),namespace=VALUES(namespace),
desire=VALUES(desire),desire_version=VALUES(desire_version),desire_meta=VALUES(desire_meta)
`
	return d.Exec(tx, insertSQL, params...)
}

func (d *DB) listShadowByNamesTx(tx *sqlx.Tx, namespace string, names []string) ([]*models.Shadow, error) {
	selectSQL := `
SELECT id, name, namespace, report, desire, report_meta, desire_meta, create_time, update_time, desire_version 
FROM baetyl_node_shadow WHERE namespace=? AND name IN (?)`
	qry, args, err := sqlx.In(selectSQL, namespace, names)
	if err != nil {
		return nil, err
	}
	var shadows []entities.Shadow
	if err = d.Query(tx, qry, &shadows, args...); err != nil {
		return nil, err
	}
	var res []*models.Shadow
	for _, shadow := range shadows {
		s, transErr := shadow.ToShadowModel()
		if err != nil {
			return nil, transErr
		}
		res = append(res, s)
	}
	return res, nil
}

func (d *DB) GetShadowTx(tx *sqlx.Tx, namespace, name string) (*models.Shadow, error) {
	selectSQL := `
SELECT 
id, name, namespace, report, desire, report_meta, desire_meta, create_time, update_time, desire_version 
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
INSERT INTO baetyl_node_shadow (namespace, name, report, desire, report_meta, desire_meta)
VALUES (?, ?, ?, ?, ?, ?)
`

	shd, err := entities.NewShadowFromShadowModel(shadow)
	if err != nil {
		return nil, err
	}

	return d.Exec(tx, insertSQL, shd.Namespace, shd.Name, shd.Report, shd.Desire, shd.ReportMeta, shd.DesireMeta)
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
SET desire=?, desire_version=?, desire_meta=?
WHERE namespace=? AND name=? AND desire_version=?
`
	desire, err := shadow.GetDesireString()
	if err != nil {
		return nil, err
	}
	desireMeta, err := shadow.GetDesireMetaString()
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, desire, genResourceVersion(), desireMeta, shadow.Namespace, shadow.Name, shadow.DesireVersion)
}

func (d *DB) UpdateShadowReportTx(tx *sqlx.Tx, shadow *models.Shadow) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_node_shadow
SET report=?, report_meta=?
WHERE namespace=? AND name=?
`
	report, err := shadow.GetReportString()
	if err != nil {
		return nil, err
	}
	reportMeta, err := shadow.GetReportMetaString()
	if err != nil {
		return nil, err
	}

	return d.Exec(tx, updateSQL, report, reportMeta, shadow.Namespace, shadow.Name)
}

func (d *DB) ListShadowByNamesTx(tx *sqlx.Tx, namespace string, names []string) ([]entities.Shadow, error) {
	selectSQL := `
SELECT 
id, name, namespace, report, desire, report_meta, desire_meta, create_time, update_time, desire_version
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

func (d *DB) ListShadowTx(tx *sqlx.Tx, namespace string) ([]entities.Shadow, error) {
	selectSQL := `
SELECT 
id, name, namespace, report, desire, report_meta, desire_meta, create_time, update_time, desire_version
FROM baetyl_node_shadow WHERE namespace=?
`
	var shadows []entities.Shadow

	if err := d.Query(tx, selectSQL, &shadows, namespace); err != nil {
		return nil, err
	}
	return shadows, nil
}

func genResourceVersion() string {
	return fmt.Sprintf("%d%s", time.Now().UTC().Unix(), common.RandString(6))
}

func (d *DB) ListAll(namespace string) (*models.ShadowList, error) {

	shadows, err := d.ListShadowTx(nil, namespace)
	if err != nil {
		return nil, err
	}

	total := len(shadows)
	items := make([]models.Shadow, 0, total)

	for _, s := range shadows {
		shd, _ := s.ToReportShadow()
		items = append(items, *shd)
	}
	result := &models.ShadowList{
		Total: total,
		Items: items,
	}
	return result, nil
}
