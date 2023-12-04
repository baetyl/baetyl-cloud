// Package entities 数据库存储基本结构与方法
package entities

import (
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/models"

	"github.com/stretchr/testify/assert"
)

func TestFromNamespaceModel(t *testing.T) {
	ns := &Namespace{
		Name: "testNs",
	}

	mNs := &models.Namespace{
		Name: "testNs",
	}
	res, err := FromNamespaceModel(mNs)

	assert.NoError(t, err)
	assert.Equal(t, ns.Name, res.Name)
}

func TestToNamespaceModel(t *testing.T) {
	ns := &Namespace{
		ID:   123,
		Name: "testNs",
	}

	mNs := &models.Namespace{
		Name: "testNs",
	}

	res, err := ToNamespaceModel(ns)
	assert.NoError(t, err)
	assert.Equal(t, mNs.Name, res.Name)
}
