package facade

import (
	"strings"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

func (a *facade) CreateSecret(ns string, secret *specV1.Secret) (*specV1.Secret, error) {
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
	secret, err = a.secret.Create(tx, ns, secret)
	return secret, err
}

func (a *facade) UpdateSecret(ns string, secret *specV1.Secret) (*specV1.Secret, error) {
	secret, err := a.secret.Update(ns, secret)
	if err != nil {
		return nil, err
	}
	err = a.updateAppSecret(ns, secret)
	if err != nil {
		return nil, err
	}
	return secret, err
}

func (a *facade) DeleteSecret(ns, name string) error {
	return a.secret.Delete(nil, ns, name)
}

func (a *facade) updateAppSecret(namespace string, secret *specV1.Secret) error {
	appNames, err := a.index.ListAppIndexBySecret(namespace, secret.Name)
	if err != nil {
		return err
	}
	for _, appName := range appNames {
		app, err := a.app.Get(namespace, appName, "")
		if err != nil {
			if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
				continue
			}
			return err
		}
		if !needUpdateAppSecret(secret, app) {
			continue
		}
		app, err = a.app.Update(nil, namespace, app)
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

func needUpdateAppSecret(secret *specV1.Secret, app *specV1.Application) bool {
	appNeedUpdate := false
	for _, volume := range app.Volumes {
		if volume.Secret != nil &&
			volume.Secret.Name == secret.Name &&
			// secret's version must increment
			strings.Compare(secret.Version, volume.Secret.Version) > 0 {
			volume.Secret.Version = secret.Version
			appNeedUpdate = true
		}
	}
	return appNeedUpdate
}
