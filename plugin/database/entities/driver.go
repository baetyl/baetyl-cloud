// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/jinzhu/copier"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Driver struct {
	ID            int64     `db:"id"`
	Name          string    `db:"name"`
	Namespace     string    `db:"namespace"`
	Version       string    `db:"version"`
	Type          byte      `db:"type"`
	Mode          string    `db:"mode"`
	Labels        string    `db:"labels"`
	Protocol      string    `db:"protocol"`
	Architecture  string    `db:"arch"`
	Description   string    `db:"description"`
	DefaultConfig string    `db:"default_config"`
	Service       string    `db:"service"`
	Volumes       string    `db:"volumes"`
	Registries    string    `db:"registries"`
	ProgramConfig string    `db:"program_config"`
	CreateTime    time.Time `db:"create_time"`
	UpdateTime    time.Time `db:"update_time"`
}

func FromModelDriver(dv *models.Driver) (*Driver, error) {
	if dv.Volumes == nil {
		dv.Volumes = make([]v1.Volume, 0)
	}
	volumes, err := json.Marshal(dv.Volumes)
	if err != nil {
		return nil, err
	}
	if dv.Registries == nil {
		dv.Registries = make([]models.RegistryView, 0)
	}
	registries, err := json.Marshal(dv.Registries)
	if err != nil {
		return nil, err
	}
	if dv.Service == nil {
		dv.Service = new(models.Service)
	}
	service, err := json.Marshal(dv.Service)
	if err != nil {
		return nil, err
	}
	labels, err := json.Marshal(dv.Labels)
	if err != nil {
		return nil, errors.Trace(err)
	}
	dv.Version = GenResourceVersion()
	driver := &Driver{
		Service:    string(service),
		Volumes:    string(volumes),
		Registries: string(registries),
		Labels:     string(labels),
	}
	if err := copier.Copy(driver, dv); err != nil {
		return nil, err
	}
	return driver, nil
}

func ToModelDriver(dv *Driver) (*models.Driver, error) {
	var service models.Service
	if err := json.Unmarshal([]byte(dv.Service), &service); err != nil {
		return nil, err
	}
	var volumes []v1.Volume
	if err := json.Unmarshal([]byte(dv.Volumes), &volumes); err != nil {
		return nil, err
	}
	var registries []models.RegistryView
	if err := json.Unmarshal([]byte(dv.Registries), &registries); err != nil {
		return nil, err
	}
	labels := map[string]string{}
	err := json.Unmarshal([]byte(dv.Labels), &labels)
	if err != nil {
		return nil, errors.Trace(err)
	}
	driver := &models.Driver{
		Service:    &service,
		Volumes:    volumes,
		Registries: registries,
		Labels:     labels,
	}
	if err := copier.Copy(driver, dv); err != nil {
		return nil, err
	}
	return driver, nil
}
