package api

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/api/init.go -package=api github.com/baetyl/baetyl-cloud/v2/api InitAPI

type InitAPI struct {
	Init service.InitService
	Auth service.AuthService
}

func NewInitAPI(cfg *config.CloudConfig) (*InitAPI, error) {
	initService, err := service.NewInitService(cfg)
	if err != nil {
		return nil, err
	}
	authService, err := service.NewAuthService(cfg)
	if err != nil {
		return nil, err
	}
	return &InitAPI{
		Init: initService,
		Auth: authService,
	}, nil
}

func (api *InitAPI) GetResource(c *common.Context) (interface{}, error) {
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
	ns, nodeName, err := CheckAndParseToken(query.Token, api.Auth.GenToken)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	return api.Init.GetResource(ns, nodeName, resourceName, map[string]interface{}{
		"Token":        query.Token,
		"KubeNodeName": query.Node,
	})
}

func CheckAndParseToken(token string, genToken func(map[string]interface{}) (string, error)) (ns, nodeName string, err error) {
	// check len
	if len(token) < 10 {
		log.L().Info("invalid token length")
		err = common.Error(common.ErrInvalidToken)
		return
	}
	// check sign
	data, err := hex.DecodeString(token[10:])
	if err != nil {
		log.L().Info("invalid token string hex", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}
	info := map[string]interface{}{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.L().Info("invalid token string json", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}
	realToken, err := genToken(info)
	if err != nil {
		log.L().Info("invalid token struct", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}
	if realToken != token {
		log.L().Info("token not match", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}

	ns, ok := info[service.InfoNamespace].(string)
	if !ok {
		log.L().Info("invalid token no namespace", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}

	nodeName, ok = info[service.InfoName].(string)
	if !ok {
		log.L().Info("invalid token no node name", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}

	expiry, ok := info[service.InfoExpiry].(float64)
	if !ok {
		log.L().Info("invalid token no expiry", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}

	ts, ok := info[service.InfoTimestamp].(float64)
	if !ok {
		log.L().Info("invalid token no timestamp", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}

	// check expiration
	timestamp := time.Unix(int64(ts), 0)
	if timestamp.Add(time.Duration(int64(expiry))*time.Second).Unix() < time.Now().Unix() {
		log.L().Info("token expired", log.Error(err))
		err = common.Error(common.ErrInvalidToken)
		return
	}
	return
}
