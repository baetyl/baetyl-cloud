package api

import (
	"encoding/hex"
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/api/init.go -package=api github.com/baetyl/baetyl-cloud/v2/api InitAPI

type InitAPI struct {
	Init service.InitService
	Sign service.SignService
}

func NewInitAPI(cfg *config.CloudConfig) (*InitAPI, error) {
	initService, err := service.NewInitService(cfg)
	if err != nil {
		return nil, err
	}
	signService, err := service.NewSignService(cfg)
	if err != nil {
		return nil, err
	}
	return &InitAPI{
		Init: initService,
		Sign: signService,
	}, nil
}

func (api *InitAPI) GetResource(c *common.Context) (interface{}, error) {
	resourceName := c.Param("resource")
	query := &struct {
		Token         string `form:"token,omitempty"`
		Node          string `form:"node,omitempty"`
		InitApplyYaml string `form:"initApplyYaml,omitempty"`
		Mode          string `form:"mode,omitempty"`
		Path          string `form:"path,omitempty"`
	}{}
	err := c.Bind(query)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	data, err := CheckAndParseToken(query.Token, api.Sign.GenToken)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	return api.Init.GetResource(data[service.InfoNamespace].(string), data[service.InfoName].(string), resourceName, map[string]interface{}{
		"Token":          query.Token,
		"KubeNodeName":   query.Node,
		"InitApplyYaml":  query.InitApplyYaml,
		"Mode":           query.Mode,
		"BaetylHostPath": query.Path,
	})
}

func CheckAndParseToken(token string, genToken func(map[string]interface{}) (string, error)) (map[string]interface{}, error) {
	// check len
	if len(token) < 10 {
		log.L().Info("invalid token length")
		return nil, common.Error(common.ErrInvalidToken)
	}
	// check sign
	data, err := hex.DecodeString(token[10:])
	if err != nil {
		log.L().Info("invalid token string hex", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}
	info := map[string]interface{}{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.L().Info("invalid token string json", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}
	realToken, err := genToken(info)
	if err != nil {
		log.L().Info("invalid token struct", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}
	if realToken != token {
		log.L().Info("token not match", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}

	_, ok := info[service.InfoNamespace].(string)
	if !ok {
		log.L().Info("invalid token no namespace", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}

	_, ok = info[service.InfoName].(string)
	if !ok {
		log.L().Info("invalid token no node name", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}

	expiry, ok := info[service.InfoExpiry].(float64)
	if !ok {
		log.L().Info("invalid token no expiry", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}

	// check expiration
	timestamp := time.Unix(int64(expiry), 0)
	if timestamp.Unix() < time.Now().Unix() {
		log.L().Info("token expired", log.Error(err))
		return nil, common.Error(common.ErrInvalidToken)
	}
	return info, nil
}
