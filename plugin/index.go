package plugin

import (
	"database/sql"
	"io"

	"github.com/jmoiron/sqlx"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

//go:generate mockgen -destination=../mock/plugin/index.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Index

// Index interface of Index
type Index interface {
	// index
	CreateIndex(namespace string, keyA, keyB common.Resource, valueA, valueB string) (sql.Result, error)
	ListIndex(namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error)
	DeleteIndex(namespace string, keyA, byKeyB common.Resource, valueB string) (sql.Result, error)
	CreateIndexTx(tx *sqlx.Tx, namespace string, keyA, keyB common.Resource, valueA, valueB string) (sql.Result, error)
	ListIndexTx(tx *sqlx.Tx, namespace string, keyA, byKeyB common.Resource, valueB string) ([]string, error)
	DeleteIndexTx(tx *sqlx.Tx, namespace string, keyA, byKeyB common.Resource, valueB string) (sql.Result, error)
	RefreshIndex(namespace string, keyA, keyB common.Resource, valueA string, valueBs []string) error

	ListResourcesByNamespace(namespace string, keyA, keyB common.Resource) ([]string, error)

	DeleteIndexByNamespace(namespace string, keyA, keyB common.Resource) (sql.Result, error)
	DeleteIndexByNamespaceTx(tx *sqlx.Tx, namespace string, keyA, keyB common.Resource) (sql.Result, error)
	io.Closer
}
