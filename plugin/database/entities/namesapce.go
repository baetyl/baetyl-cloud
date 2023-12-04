// Package entities 数据库存储基本结构与方法
package entities

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Namespace struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func FromNamespaceModel(namespace *models.Namespace) (*Namespace, error) {
	return &Namespace{
		Name: namespace.Name,
	}, nil
}

func ToNamespaceModel(namespace *Namespace) (*models.Namespace, error) {
	return &models.Namespace{
		Name: namespace.Name,
	}, nil
}
