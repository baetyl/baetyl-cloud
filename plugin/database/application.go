// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetApplication(tx interface{}, namespace, name, _ string) (*specV1.Application, error) {
	defer utils.Trace(d.Log.Debug, "GetApplication")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	return d.GetApplicationTx(transaction, namespace, name)
}

func (d *BaetylCloudDB) CreateApplication(tx interface{}, namespace string, application *specV1.Application) (*specV1.Application, error) {
	var app *specV1.Application
	var err error
	defer utils.Trace(d.Log.Debug, "CreateApplication")()
	if tx == nil {
		err = d.Transact(func(tx *sqlx.Tx) error {
			return d.createAndGetApplication(tx, namespace, application, &app)
		})
	} else {
		err = d.createAndGetApplication(tx.(*sqlx.Tx), namespace, application, &app)
	}
	return app, err
}

func (d *BaetylCloudDB) createAndGetApplication(tx *sqlx.Tx, ns string, in *specV1.Application, out **specV1.Application) error {
	_, err := d.CreateApplicationTx(tx, ns, in)
	if err != nil {
		return err
	}
	*out, err = d.GetApplicationTx(tx, ns, in.Name)
	return err
}

func (d *BaetylCloudDB) UpdateApplication(tx interface{}, namespace string, application *specV1.Application) (*specV1.Application, error) {
	defer utils.Trace(d.Log.Debug, "UpdateApplication")()
	var err error
	if tx == nil {
		err = d.Transact(func(tx *sqlx.Tx) error {
			return d.updateAndGetApplication(tx, namespace, application)
		})
	} else {
		err = d.updateAndGetApplication(tx.(*sqlx.Tx), namespace, application)
	}
	return application, err
}

func (d *BaetylCloudDB) updateAndGetApplication(tx *sqlx.Tx, ns string, in *specV1.Application) error {
	oldApp, err := d.GetApplicationTx(tx, ns, in.Name)
	if err != nil {
		return err
	}
	if entities.EqualApp(in, oldApp) {
		in.Version = oldApp.Version
		return nil
	}
	_, err = d.UpdateApplicationTx(tx, ns, in)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaetylCloudDB) DeleteApplication(tx interface{}, namespace, name string) error {
	defer utils.Trace(d.Log.Debug, "DeleteApplication")()
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return err
	}
	_, err = d.DeleteApplicationTx(transaction, namespace, name)
	return err
}

func (d *BaetylCloudDB) ListApplication(tx interface{}, namespace string, listOptions *models.ListOptions) (*models.ApplicationList, error) {
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(d.Log.Debug, "ListApplication")()
	apps, resLen, err := d.ListApplicationTx(transaction, namespace, listOptions)
	if err != nil {
		return nil, err
	}

	result := &models.ApplicationList{
		Total:       resLen,
		ListOptions: listOptions,
		Items:       apps,
	}
	return result, nil
}

func (d *BaetylCloudDB) ListApplicationsByNames(tx interface{}, ns string, names []string) ([]models.AppItem, int, error) {
	transaction, err := d.InterfaceToTx(tx)
	if err != nil {
		return nil, 0, errors.Trace(err)
	}
	defer utils.Trace(d.Log.Debug, "ListApplicationsByNames")()
	apps, resLen, err := d.ListApplicationsByNamesTx(transaction, ns, names)
	if err != nil {
		return nil, 0, errors.Trace(err)
	}
	return apps, resLen, nil
}

func (d *BaetylCloudDB) GetApplicationTx(tx *sqlx.Tx, namespace, name string) (*specV1.Application, error) {
	selectSQL := `
SELECT 
id, namespace, name, version, type, mode, is_system, create_time, labels, 
selector, node_selector, description, services, init_services, volumes, cron_status, 
update_time, cron_time, workload, host_network, dns_policy, replica, job_config ,ota, autoScaleCfg,preserve_updates
FROM baetyl_application WHERE namespace=? AND name=?
`
	var apps []entities.Application
	if err := d.Query(tx, selectSQL, &apps, namespace, name); err != nil {
		return nil, err
	}
	if len(apps) > 0 {
		return entities.ToAppModel(&apps[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "application"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) CreateApplicationTx(tx *sqlx.Tx, namespace string, application *specV1.Application) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_application (namespace, name, version, type, mode, 
is_system, labels, selector, node_selector, description, services, init_services, volumes, 
cron_status, cron_time, workload, host_network, dns_policy, replica, job_config, ota, autoScaleCfg, preserve_updates)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	app, err := entities.FromAppModel(namespace, application)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return d.Exec(tx, insertSQL, app.Namespace, app.Name, app.Version,
		app.Type, app.Mode, app.System, app.Labels, app.Selector,
		app.NodeSelector, app.Description, app.Services, app.InitService, app.Volumes,
		app.CronStatus, app.CronTime, app.Workload, app.HostNetwork, app.DNSPolicy,
		app.Replica, app.JobConfig, app.Ota, app.AutoScaleCfg, app.PreserveUpdates)
}

func (d *BaetylCloudDB) UpdateApplicationTx(tx *sqlx.Tx, namespace string, application *specV1.Application) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_application
SET version = ?, description = ?, labels = ?, selector = ?, 
node_selector = ?, services = ?, init_services = ?, volumes = ?, cron_status=?, 
cron_time=?, workload=?, host_network=?, dns_policy=?, replica=?, job_config=?, ota=?, autoScaleCfg=?, preserve_updates=?
WHERE namespace=? AND name=?
`
	app, err := entities.FromAppModel(namespace, application)
	if err != nil {
		return nil, err
	}
	application.Version = app.Version
	return d.Exec(tx, updateSQL, app.Version, app.Description, app.Labels,
		app.Selector, app.NodeSelector, app.Services, app.InitService, app.Volumes,
		app.CronStatus, app.CronTime, app.Workload, app.HostNetwork, app.DNSPolicy,
		app.Replica, app.JobConfig, app.Ota, app.AutoScaleCfg, app.PreserveUpdates, app.Namespace, app.Name)
}

func (d *BaetylCloudDB) ListApplicationTx(tx *sqlx.Tx, namespace string, listOptions *models.ListOptions) ([]models.AppItem, int, error) {
	selectSQL := `
SELECT 
id, namespace, name, version, type, mode, is_system, labels, 
selector, node_selector, description, services, init_services, volumes, 
create_time, cron_status, update_time, cron_time, 
workload, host_network, replica, job_config , ota, autoScaleCfg, preserve_updates
FROM baetyl_application WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC
`
	var applications []entities.Application
	if err := d.Query(tx, selectSQL, &applications, namespace, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	result := make([]models.AppItem, 0)
	for _, application := range applications {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(application.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		app := entities.ToAppListModel(&application)
		result = append(result, *app)
	}
	start, end := models.GetPagingParam(listOptions, len(result))
	return result[start:end], len(result), nil
}

func (d *BaetylCloudDB) DeleteApplicationTx(tx *sqlx.Tx, namespace, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_application WHERE namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, namespace, name)
}

func (d *BaetylCloudDB) ListApplicationsByNamesTx(tx *sqlx.Tx, ns string, names []string) ([]models.AppItem, int, error) {
	selectSQL := `
SELECT 
id, namespace, name, version, type, mode, is_system, labels, 
selector, node_selector, description, services, init_services, volumes, 
create_time, cron_status, update_time, cron_time, 
workload, host_network, replica, job_config , ota, autoScaleCfg, preserve_updates
FROM baetyl_application WHERE namespace=? AND name IN (?)
`
	qry, args, err := sqlx.In(selectSQL, ns, names)
	if err != nil {
		return nil, 0, err
	}
	var applications []entities.Application
	if err = d.Query(tx, qry, &applications, args...); err != nil {
		return nil, 0, err
	}
	result := make([]models.AppItem, 0)
	for _, application := range applications {
		app := entities.ToAppListModel(&application)
		result = append(result, *app)
	}
	return result, len(result), nil
}
