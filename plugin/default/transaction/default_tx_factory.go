package transaction

import (
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func init() {
	plugin.RegisterFactory("defaulttx", New)
}

type defaultTxFactory struct{}

func New() (plugin.Plugin, error) {
	return &defaultTxFactory{}, nil
}

func (t *defaultTxFactory) BeginTx() (interface{}, error) {
	return nil, nil
}

func (t *defaultTxFactory) Commit(tx interface{}) {}

func (t *defaultTxFactory) Rollback(tx interface{}) {}

func (t *defaultTxFactory) Close() error {
	return nil
}
