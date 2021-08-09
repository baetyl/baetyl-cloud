package database

import (
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func(d *DB) GetCron(name, namespace string) (*models.Cron, error) {
	selectSQL := `SELECT name, namespace, selector, cron_time FROM baetyl_cron_app WHERE name=? AND namespace=?`
	var cronApps []entities.CronApp
	err := d.Query(nil, selectSQL, &cronApps, name, namespace)
	if err != nil {
		return nil, err
	}
	if len(cronApps) > 0 {
		return &models.Cron{
			Name: cronApps[0].Name,
			Namespace: cronApps[0].Namespace,
			Selector: cronApps[0].Selector,
			CronTime: cronApps[0].CronTime.UTC(),
		}, nil
	}
	return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "cronApp"), common.Field("name", name))
}

func(d *DB) CreateCron(cronApp *models.Cron) error {
	insertSQL := `INSERT INTO baetyl_cron_app (name, namespace, selector, cron_time) VALUES (?,?,?,?)`
	_, err := d.Exec(nil, insertSQL, cronApp.Name, cronApp.Namespace, cronApp.Selector, cronApp.CronTime)
	return err
}

func(d *DB) UpdateCron(cronApp *models.Cron) error {
	updateSQL := `UPDATE baetyl_cron_app SET selector=?, cron_time=? WHERE name=? AND namespace=?`
	_, err := d.Exec(nil, updateSQL, cronApp.Selector, cronApp.CronTime ,cronApp.Name, cronApp.Namespace)
	return err
}

func(d *DB) DeleteCron(name, namespace string) error {
	deleteSQL := `DELETE FROM baetyl_cron_app WHERE name=? AND namespace=?`
	_, err := d.Exec(nil, deleteSQL, name, namespace)
	return err
}

func(d *DB) ListExpiredApps() ([]models.Cron, error) {
	var applications []entities.CronApp
	selectSQL := `
SELECT id, name, namespace, selector, cron_time 
FROM baetyl_cron_app WHERE cron_time <= now()
	`
	if err := d.Query(nil, selectSQL, &applications); err != nil {
		return nil, err
	}
	apps := make([]models.Cron, 0)
	for _, application := range applications {
		apps = append(apps, models.Cron{
			Id: application.Id,
			Name: application.Name,
			Namespace: application.Namespace,
			Selector: application.Selector,
			CronTime: application.CronTime.UTC(),
		})
	}
	return apps, nil
}

func(d *DB) DeleteExpiredApps(cronApps []uint64) error {
	deleteSql := `DELETE FROM baetyl_cron_app WHERE id IN (?)`
	dSql, args, err := sqlx.In(deleteSql, cronApps)
	if err != nil {
		return err
	}
	if _, err = d.Exec(nil, dSql, args...); err != nil {
		return err
	}
	return nil
}
