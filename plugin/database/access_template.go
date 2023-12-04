// Package database 数据库存储实现
package database

import (
	"database/sql"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

func (d *BaetylCloudDB) GetAccessTemplate(ns, name string) (*models.AccessTemplate, error) {
	return d.GetAccessTemplateTx(nil, ns, name)
}

func (d *BaetylCloudDB) ListAccessTemplate(ns string, listOptions *models.ListOptions) (*models.AccessTemplateList, error) {
	templates, count, err := d.ListAccessTemplateTx(nil, ns, listOptions)
	if err != nil {
		return nil, err
	}
	res := &models.AccessTemplateList{
		Items:       templates,
		Total:       count,
		ListOptions: listOptions,
	}
	return res, nil
}

func (d *BaetylCloudDB) ListAccessTemplateByModelAndProtocol(ns, model, protocol string, listOptions *models.ListOptions) (*models.AccessTemplateList, error) {
	templates, count, err := d.ListAccessTemplateByModelAndProtocolTx(nil, ns, model, protocol, listOptions)
	if err != nil {
		return nil, err
	}
	res := &models.AccessTemplateList{
		Items:       templates,
		Total:       count,
		ListOptions: listOptions,
	}
	return res, nil
}

func (d *BaetylCloudDB) CreateAccessTemplate(template *models.AccessTemplate) (*models.AccessTemplate, error) {
	var res *models.AccessTemplate
	err := d.Transact(func(tx *sqlx.Tx) error {
		_, err := d.CreateAccessTemplateTx(tx, template)
		if err != nil {
			return err
		}
		res, err = d.GetAccessTemplateTx(tx, template.Namespace, template.Name)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) UpdateAccessTemplate(template *models.AccessTemplate) (*models.AccessTemplate, error) {
	var res *models.AccessTemplate
	err := d.Transact(func(tx *sqlx.Tx) error {
		old, err := d.GetAccessTemplateTx(tx, template.Namespace, template.Name)
		if err != nil {
			return err
		}
		if models.EqualAccessTemplate(old, template) {
			res = old
			return nil
		}
		_, err = d.UpdateAccessTemplateTx(tx, template)
		if err != nil {
			return err
		}
		res, err = d.GetAccessTemplateTx(tx, template.Namespace, template.Name)
		return err
	})
	return res, err
}

func (d *BaetylCloudDB) DeleteAccessTemplate(ns, name string) error {
	_, err := d.DeleteAccessTemplateTx(nil, ns, name)
	return err
}

func (d *BaetylCloudDB) GetAccessTemplateTx(tx *sqlx.Tx, ns, name string) (*models.AccessTemplate, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol, labels, device_model,
mappings, properties, create_time, update_time
FROM baetyl_access_template WHERE namespace=? AND name=? LIMIT 0,1
`
	var templates []entities.AccessTemplate
	if err := d.Query(tx, selectSQL, &templates, ns, name); err != nil {
		return nil, err
	}
	if len(templates) > 0 {
		return entities.ToAccessTemplate(&templates[0])
	}
	return nil, common.Error(
		common.ErrResourceNotFound,
		common.Field("type", "AccessTemplate"),
		common.Field("name", name))
}

func (d *BaetylCloudDB) ListAccessTemplateTx(tx *sqlx.Tx, ns string, listOptions *models.ListOptions) ([]models.AccessTemplate, int, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol, labels, device_model,
mappings, properties, create_time, update_time
FROM baetyl_access_template WHERE namespace=? AND name LIKE ? ORDER BY create_time DESC 
`
	var templates []entities.AccessTemplate
	if err := d.Query(tx, selectSQL, &templates, ns, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	result := make([]models.AccessTemplate, 0)
	for _, template := range templates {
		labels := map[string]string{}
		if err := json.Unmarshal([]byte(template.Labels), &labels); err != nil {
			return nil, 0, errors.Trace(err)
		}
		if ok, err := utils.IsLabelMatch(listOptions.LabelSelector, labels); err != nil || !ok {
			continue
		}
		temp, err := entities.ToAccessTemplate(&template)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, *temp)
	}
	start, end := models.GetPagingParam(listOptions, len(result))
	return result[start:end], len(result), nil
}

func (d *BaetylCloudDB) ListAccessTemplateByModelAndProtocolTx(tx *sqlx.Tx, ns, model, protocol string, listOptions *models.ListOptions) ([]models.AccessTemplate, int, error) {
	selectSQL := `
SELECT  
name, namespace, version, description, protocol, labels, device_model,
mappings, properties, create_time, update_time
FROM baetyl_access_template WHERE namespace=? AND device_model=? AND protocol=? AND name LIKE ? ORDER BY create_time DESC 
`
	var templates []entities.AccessTemplate
	if err := d.Query(tx, selectSQL, &templates, ns, model, protocol, listOptions.GetFuzzyName()); err != nil {
		return nil, 0, err
	}
	var res []models.AccessTemplate
	for _, temp := range templates {
		accessTemplate, err := entities.ToAccessTemplate(&temp)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, *accessTemplate)
	}
	start, end := models.GetPagingParam(listOptions, len(res))
	return res[start:end], len(res), nil
}

func (d *BaetylCloudDB) CreateAccessTemplateTx(tx *sqlx.Tx, template *models.AccessTemplate) (sql.Result, error) {
	insertSQL := `
INSERT INTO baetyl_access_template
(name, namespace, version, description, protocol, labels, device_model,
mappings, properties)
VALUES 
(?,?,?,?,?,?,?,?,?)
`
	accessTemplate, err := entities.FromAccessTemplate(template)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, insertSQL, accessTemplate.Name, accessTemplate.Namespace, accessTemplate.Version, accessTemplate.Description,
		accessTemplate.Protocol, accessTemplate.Labels, accessTemplate.DeviceModel, accessTemplate.Mappings, accessTemplate.Properties)
}

func (d *BaetylCloudDB) UpdateAccessTemplateTx(tx *sqlx.Tx, template *models.AccessTemplate) (sql.Result, error) {
	updateSQL := `
UPDATE baetyl_access_template SET version=?, description=?, labels=?,
mappings=?, properties=? WHERE namespace=? AND name=?
`
	devModel, err := entities.FromAccessTemplate(template)
	if err != nil {
		return nil, err
	}
	return d.Exec(tx, updateSQL, devModel.Version, devModel.Description, devModel.Labels,
		devModel.Mappings, devModel.Properties, devModel.Namespace, devModel.Name)
}

func (d *BaetylCloudDB) DeleteAccessTemplateTx(tx *sqlx.Tx, ns, name string) (sql.Result, error) {
	deleteSQL := `
DELETE FROM baetyl_access_template where namespace=? AND name=?
`
	return d.Exec(tx, deleteSQL, ns, name)
}
