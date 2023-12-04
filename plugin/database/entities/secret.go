// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Secret struct {
	ID          int64     `db:"id"`
	Namespace   string    `db:"namespace"`
	Name        string    `db:"name"`
	Labels      string    `db:"labels"`
	Data        string    `db:"data"`
	Description string    `db:"description"`
	Version     string    `db:"version"`
	System      bool      `db:"is_system"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func ToSecretModel(secret *Secret) (*specV1.Secret, error) {
	labels := map[string]string{}
	err := json.Unmarshal([]byte(secret.Labels), &labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	data := map[string][]byte{}
	err = json.Unmarshal([]byte(secret.Data), &data)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &specV1.Secret{
		Namespace:         secret.Namespace,
		Name:              secret.Name,
		Version:           secret.Version,
		System:            secret.System,
		CreationTimestamp: secret.CreateTime.UTC(),
		UpdateTimestamp:   secret.UpdateTime.UTC(),
		Labels:            labels,
		Data:              data,
		Description:       secret.Description,
	}, nil
}

func FromSecretModel(namespace string, secret *specV1.Secret) (*Secret, error) {
	labels, err := json.Marshal(secret.Labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	data, err := json.Marshal(secret.Data)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if len(data) >= common.MaxSecretDataSize {
		return nil, common.Error(common.ErrDataTooLarge, common.Field("name", secret.Name),
			common.Field("size", len(data)), common.Field("max", common.MaxSecretDataSize-1))
	}

	return &Secret{
		Name:        secret.Name,
		Namespace:   namespace,
		Version:     GenResourceVersion(),
		System:      secret.System,
		CreateTime:  secret.CreationTimestamp.UTC(),
		UpdateTime:  secret.UpdateTimestamp.UTC(),
		Labels:      string(labels),
		Data:        string(data),
		Description: secret.Description,
	}, nil
}
