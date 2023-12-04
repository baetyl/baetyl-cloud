// Package entities 数据库存储基本结构与方法
package entities

import (
	"reflect"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Configuration struct {
	ID          int64     `db:"id"`
	Namespace   string    `db:"namespace"`
	Name        string    `db:"name"`
	Labels      string    `db:"labels"`
	Data        string    `db:"data"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	Description string    `db:"description"`
	Version     string    `db:"version"`
	System      bool      `db:"is_system"`
}

func ToConfigModel(config *Configuration) (*specV1.Configuration, error) {
	labels := map[string]string{}
	err := json.Unmarshal([]byte(config.Labels), &labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	data := map[string]string{}
	err = json.Unmarshal([]byte(config.Data), &data)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &specV1.Configuration{
		Namespace:         config.Namespace,
		Name:              config.Name,
		Version:           config.Version,
		System:            config.System,
		CreationTimestamp: config.CreateTime.UTC(),
		UpdateTimestamp:   config.UpdateTime.UTC(),
		Labels:            labels,
		Data:              data,
		Description:       config.Description,
	}, nil
}

func EqualConfig(oldCfg, newCfg *specV1.Configuration) bool {
	if oldCfg.Name != newCfg.Name || oldCfg.Namespace != newCfg.Namespace ||
		oldCfg.Description != newCfg.Description || oldCfg.System != newCfg.System {
		return false
	}
	return len(oldCfg.Labels) == len(newCfg.Labels) &&
		reflect.DeepEqual(oldCfg.Labels, newCfg.Labels) && reflect.DeepEqual(oldCfg.Data, newCfg.Data)
}

func FromConfigModel(namespace string, config *specV1.Configuration) (*Configuration, error) {
	labels, err := json.Marshal(config.Labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	data, err := json.Marshal(config.Data)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if len(data) >= common.MaxConfigDataSize {
		return nil, common.Error(common.ErrDataTooLarge, common.Field("name", config.Name),
			common.Field("size", len(data)), common.Field("max", common.MaxConfigDataSize-1))
	}

	return &Configuration{
		Name:        config.Name,
		Namespace:   namespace,
		Version:     GenResourceVersion(),
		System:      config.System,
		CreateTime:  config.CreationTimestamp.UTC(),
		UpdateTime:  config.UpdateTimestamp.UTC(),
		Labels:      string(labels),
		Data:        string(data),
		Description: config.Description,
	}, nil
}
