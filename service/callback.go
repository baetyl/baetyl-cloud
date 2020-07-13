package service

import (
	"bytes"
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
	"github.com/baetyl/baetyl-go/v2/http"
	"github.com/baetyl/baetyl-go/v2/log"
)

//go:generate mockgen -destination=../mock/service/callback.go -package=plugin github.com/baetyl/baetyl-cloud/service CallbackService

type CallbackService interface {
	Create(callback *models.Callback) (*models.Callback, error)
	Delete(name, ns string) error
	Update(callback *models.Callback) (*models.Callback, error)
	Get(name, ns string) (*models.Callback, error)
	Callback(name, ns string, data map[string]string) ([]byte, error)
}

type callbackService struct {
	cfg          *config.CloudConfig
	modelStorage plugin.ModelStorage
	dbStorage    plugin.DBStorage
	http         *http.Client
}

// NewCallbackService New Callback Service
func NewCallbackService(config *config.CloudConfig) (CallbackService, error) {
	ms, err := plugin.GetPlugin(config.Plugin.ModelStorage)
	if err != nil {
		return nil, err
	}
	ds, err := plugin.GetPlugin(config.Plugin.DatabaseStorage)
	if err != nil {
		return nil, err
	}
	return &callbackService{
		cfg:          config,
		modelStorage: ms.(plugin.ModelStorage),
		dbStorage:    ds.(plugin.DBStorage),
		http:         http.NewClient(http.NewClientOptions()),
	}, nil
}

func (c *callbackService) Create(callback *models.Callback) (*models.Callback, error) {
	_, err := c.dbStorage.CreateCallback(callback)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	res, err := c.dbStorage.GetCallback(callback.Name, callback.Namespace)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	return res, nil
}

func (c *callbackService) Delete(name, ns string) error {
	count, err := c.dbStorage.CountBatchByCallback(name, ns)
	if err != nil {
		return common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if count > 0 {
		return common.Error(common.ErrRegisterDeleteCallback, common.Field("name", name))
	}
	res, err := c.dbStorage.DeleteCallback(name, ns)
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return common.Error(common.ErrResourceNotFound, common.Field("type", "callback"), common.Field("name", name))
	}
	return nil
}

func (c *callbackService) Update(callback *models.Callback) (*models.Callback, error) {
	_, err := c.dbStorage.UpdateCallback(callback)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	res, err := c.dbStorage.GetCallback(callback.Name, callback.Namespace)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *callbackService) Get(name, ns string) (*models.Callback, error) {
	call, err := c.dbStorage.GetCallback(name, ns)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if call == nil {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "callback"), common.Field("name", name))
	}
	return call, nil
}

func (c *callbackService) Callback(name, ns string, data map[string]string) ([]byte, error) {
	call, err := c.dbStorage.GetCallback(name, ns)
	if err != nil {
		return nil, err
	}
	queryParams := ""
	for k, v := range call.Params {
		queryParams += "&" + k + "=" + v
	}
	if len(queryParams) > 0 {
		queryParams = "?" + queryParams[1:]
	}
	for k, v := range data {
		call.Body[k] = v
	}

	payload, err := json.Marshal(call.Body)
	if err != nil {
		return nil, err
	}
	res, err := c.http.SendUrl(call.Method, call.Url+queryParams, bytes.NewReader(payload), call.Header)
	if err != nil {
		return nil, err
	}
	buf, err := http.HandleResponse(res)
	if err != nil {
		return nil, err
	}
	log.L().Info("callback success", log.Any("name", call.Name), log.Any("res", string(buf)))
	return buf, nil
}
