package facade

import (
	"strings"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func (a *facade) CreateConfig(ns string, config *specV1.Configuration) (*specV1.Configuration, error) {
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
	config, err = a.config.Create(tx, ns, config)
	return config, err
}

func (a *facade) UpdateConfig(ns string, config *specV1.Configuration) (*specV1.Configuration, error) {
	var res *specV1.Configuration
	var err error
	res, err = a.config.Update(nil, ns, config)
	if err != nil {
		log.L().Error("Update config failed", log.Error(err))
		return nil, err
	}

	var appNames []string
	appNames, err = a.index.ListAppIndexByConfig(ns, res.Name)
	if err != nil {
		log.L().Error("list app index by config failed", log.Error(err))
		return nil, err
	}

	if err = a.updateNodeAndApp(ns, res, appNames); err != nil {
		log.L().Error("update node and app failed", log.Error(err))
		return nil, err
	}
	return res, err
}

func (a *facade) DeleteConfig(ns, name string) error {
	return a.config.Delete(nil, ns, name)
}

func (a *facade) updateNodeAndApp(namespace string, config *specV1.Configuration, appNames []string) error {
	for _, appName := range appNames {
		app, err := a.app.Get(namespace, appName, "")
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			return err
		}

		if !needUpdateApp(config, app) {
			continue
		}
		// Todo remove by list watch
		app, err = a.app.Update(namespace, app)
		if err != nil {
			return err
		}
		_, err = a.node.UpdateNodeAppVersion(nil, namespace, app)
		if err != nil {
			return err
		}
	}
	return nil
}

func needUpdateApp(config *specV1.Configuration, app *specV1.Application) bool {
	appNeedUpdate := false
	for _, volume := range app.Volumes {
		if volume.Config != nil &&
			volume.Config.Name == config.Name &&
			// config's version must increment
			strings.Compare(config.Version, volume.Config.Version) > 0 {
			volume.Config.Version = config.Version
			appNeedUpdate = true
		}
	}
	return appNeedUpdate
}
