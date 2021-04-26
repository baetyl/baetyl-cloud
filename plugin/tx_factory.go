package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/tx_factory.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin TransactionFactory

type TransactionFactory interface {
	BeginTx() (interface{}, error)
	Commit(interface{})
	Rollback(interface{})
	io.Closer
}
