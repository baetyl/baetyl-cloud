package api

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/api/init.go -package=api github.com/baetyl/baetyl-cloud/v2/api InitAPI

type InitAPI interface {
	GetResource(c *common.Context) (interface{}, error)
}

type InitAPIImpl struct {
	initService service.InitService
	auth        service.AuthService
}

func NewInitAPI(cfg *config.CloudConfig) (InitAPI, error) {
	initService, err := service.NewInitService(cfg)
	if err != nil {
		return nil, err
	}
	authService, err := service.NewAuthService(cfg)
	if err != nil {
		return nil, err
	}
	return &InitAPIImpl{
		initService: initService,
		auth:        authService,
	}, nil
}

func (api *InitAPIImpl) GetResource(c *common.Context) (interface{}, error) {
	resourceName := c.Param("resource")
	query := &struct {
		Token string `form:"token,omitempty"`
		Node  string `form:"node,omitempty"`
	}{}
	err := c.Bind(query)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	info, err := api.CheckAndParseToken(query.Token, resourceName)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	return api.initService.GetResource(resourceName, query.Node, query.Token, info)
}

func (a *InitAPIImpl) CheckAndParseToken(token, resourceName string) (map[string]interface{}, error) {
	if resourceName != common.ResourceInitYaml {
		return nil, nil
	}
	// check len
	if len(token) < 10 {
		return nil, common.Error(
			common.ErrInvalidToken,
			common.Field("error", "invalid token length"))
	}
	// check sign
	data, err := hex.DecodeString(token[10:])
	if err != nil {
		return nil, err
	}
	info := map[string]interface{}{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	realToken, err := a.auth.GenToken(info)
	if err != nil {
		return nil, err
	}
	if realToken != token {
		return nil, common.Error(
			common.ErrInvalidToken,
			common.Field("error", "token check fail"))
	}

	expiry, ok := info[service.InfoExpiry].(float64)
	if !ok {
		return nil, common.Error(
			common.ErrInvalidToken,
			common.Field("error", "expiry error"))
	}

	ts, ok := info[service.InfoTimestamp].(float64)
	if !ok {
		return nil, common.Error(
			common.ErrInvalidToken,
			common.Field("error", "infoTimestamp error"))
	}
	// check expiration
	timestamp := time.Unix(int64(ts), 0)
	if timestamp.Add(time.Duration(int64(expiry))*time.Second).Unix() < time.Now().Unix() {
		return nil, common.Error(
			common.ErrInvalidToken,
			common.Field("error", "timestamp check error"))
	}
	return info, nil
}
