package service

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

//go:generate mockgen -destination=../mock/service/wrapper.go -package=service github.com/baetyl/baetyl-cloud/v2/service WrapperService

type CreateNodeFunc func(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error)

type WrapperService interface {
	CreateNodeTx(CreateNodeFunc) CreateNodeFunc
}

type WrapperServiceImpl struct {
	plugin.TransactionFactory
}

func NewWrapperService(config *config.CloudConfig) (WrapperService, error) {
	tx, err := plugin.GetPlugin(config.Plugin.Tx)
	if err != nil {
		return nil, err
	}
	return &WrapperServiceImpl{
		tx.(plugin.TransactionFactory),
	}, nil
}

func (w *WrapperServiceImpl) Close(tx interface{}) {
	if p := recover(); p != nil {
		w.Rollback(tx)
		panic(p)
	}
}

func (w *WrapperServiceImpl) FinishTx(err error, tx interface{}) {
	if err != nil {
		w.Rollback(tx)
	} else {
		w.Commit(tx)
	}
}

func (w *WrapperServiceImpl) CreateNodeTx(function CreateNodeFunc) CreateNodeFunc {
	return func(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error) {
		transaction, err := w.BeginTx()
		if err != nil {
			return nil, err
		}
		defer w.Close(transaction)
		result, err := function(transaction, namespace, node)
		w.FinishTx(err, transaction)
		return result, err
	}
}
