package facade

import (
	"strings"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/facade/application.go -package=facade github.com/baetyl/baetyl-cloud/v2/facade ApplicationFacade

const (
	FunctionConfigPrefix        = "baetyl-function-config"
	FunctionProgramConfigPrefix = "baetyl-function-program-config"
)

type ApplicationFacade interface {
	Create(ns string, baseApp, app *specV1.Application, configs []specV1.Configuration) (*specV1.Application, error)
	Update(ns string, oldApp, app *specV1.Application, configs []specV1.Configuration) (*specV1.Application, error)
	Delete(ns, name string, app *specV1.Application) error
}

type applicationFacade struct {
	node      service.NodeService
	app       service.ApplicationService
	config    service.ConfigService
	index     service.IndexService
	txFactory plugin.TransactionFactory
	log       *log.Logger
}

func NewApplicationFacade(config *config.CloudConfig) (ApplicationFacade, error) {
	node, err := service.NewNodeService(config)
	if err != nil {
		return nil, err
	}
	app, err := service.NewApplicationService(config)
	if err != nil {
		return nil, err
	}
	cfg, err := service.NewConfigService(config)
	if err != nil {
		return nil, err
	}
	index, err := service.NewIndexService(config)
	if err != nil {
		return nil, err
	}
	tx, err := plugin.GetPlugin(config.Plugin.Tx)
	if err != nil {
		return nil, err
	}

	return &applicationFacade{
		node:      node,
		app:       app,
		config:    cfg,
		index:     index,
		txFactory: tx.(plugin.TransactionFactory),
		log:       log.L().With(log.Any("facade", "application")),
	}, nil
}

func (a *applicationFacade) Create(ns string, baseApp, app *specV1.Application, configs []specV1.Configuration) (*specV1.Application, error) {
	tx, errTx := a.txFactory.BeginTx()
	if errTx != nil {
		return nil, errTx
	}
	var err error
	defer func() {
		if p := recover(); p != nil {
			a.txFactory.Rollback(tx)
			panic(p)
		} else if err != nil {
			a.txFactory.Rollback(tx)
		} else {
			a.txFactory.Commit(tx)
		}
	}()
	err = a.updateGenConfigsOfFunctionApp(tx, ns, configs)
	if err != nil {
		return nil, err
	}

	app, err = a.app.CreateWithBase(tx, ns, app, baseApp)
	if err != nil {
		return nil, err
	}

	err = a.UpdateNodeAndAppIndex(tx, ns, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *applicationFacade) Update(ns string, oldApp, app *specV1.Application, configs []specV1.Configuration) (*specV1.Application, error) {
	var err error
	err = a.updateGenConfigsOfFunctionApp(nil, ns, configs)
	if err != nil {
		return nil, err
	}

	app, err = a.app.Update(ns, app)
	if err != nil {
		return nil, err
	}

	if oldApp != nil && oldApp.Selector != app.Selector {
		// delete old nodes
		if err = a.DeleteNodeAndAppIndex(nil, ns, oldApp); err != nil {
			return nil, err
		}
	}

	// update nodes
	if err = a.UpdateNodeAndAppIndex(nil, ns, app); err != nil {
		return nil, err
	}

	a.cleanGenConfigsOfFunctionApp(nil, configs, oldApp)
	return app, nil
}

func (a *applicationFacade) Delete(ns, name string, app *specV1.Application) error {
	var err error
	if err = a.app.Delete(ns, name, ""); err != nil {
		return err
	}

	//delete the app from node
	if err = a.DeleteNodeAndAppIndex(nil, ns, app); err != nil {
		return err
	}

	a.cleanGenConfigsOfFunctionApp(nil, nil, app)
	return nil
}

func (a *applicationFacade) DeleteNodeAndAppIndex(tx interface{}, namespace string, app *specV1.Application) error {
	_, err := a.node.DeleteNodeAppVersion(tx, namespace, app)
	if err != nil {
		return err
	}

	return a.index.RefreshNodesIndexByApp(tx, namespace, app.Name, make([]string, 0))
}

func (a *applicationFacade) updateGenConfigsOfFunctionApp(tx interface{}, namespace string, configs []specV1.Configuration) error {
	for _, cfg := range configs {
		_, err := a.config.Upsert(tx, namespace, &cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *applicationFacade) UpdateNodeAndAppIndex(tx interface{}, namespace string, app *specV1.Application) error {
	nodes, err := a.node.UpdateNodeAppVersion(tx, namespace, app)
	if err != nil {
		return err
	}
	return a.index.RefreshNodesIndexByApp(tx, namespace, app.Name, nodes)
}

func (a *applicationFacade) cleanGenConfigsOfFunctionApp(tx interface{}, configs []specV1.Configuration, oldApp *specV1.Application) {
	m := map[string]bool{}
	for _, cfg := range configs {
		m[cfg.Name] = true
	}

	for _, v := range oldApp.Volumes {
		if v.VolumeSource.Config == nil {
			continue
		}
		if _, ok := m[v.VolumeSource.Config.Name]; !ok && (strings.HasPrefix(v.VolumeSource.Config.Name, FunctionConfigPrefix) ||
			strings.HasPrefix(v.VolumeSource.Config.Name, FunctionProgramConfigPrefix)) {
			err := a.config.Delete(tx, oldApp.Namespace, v.VolumeSource.Config.Name)
			if err != nil {
				common.LogDirtyData(err,
					log.Any("type", common.Config),
					log.Any(common.KeyContextNamespace, oldApp.Namespace),
					log.Any("name", v.VolumeSource.Config.Name))
				continue
			}
		}
	}
}
