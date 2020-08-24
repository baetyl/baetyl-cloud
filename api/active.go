package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

const (
	InfoKind      = "k"
	InfoName      = "n"
	InfoNamespace = "ns"
	InfoTimestamp = "ts"
	InfoExpiry    = "e"
)

var (
	ErrInvalidToken = fmt.Errorf("invalid token")
	SystemApps      = []common.SystemApplication{
		common.BaetylCore,
		common.BaetylFunction,
		common.BaetylBroker,
		common.BaetylRule,
	}
)

func (api *API) GetResource(c *common.Context) (interface{}, error) {
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

	switch resourceName {
	case common.ResourceMetrics:
		res, err := api.initService.GetResource(common.ResourceMetrics)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceLocalPathStorage:
		res, err := api.initService.GetResource(common.ResourceLocalPathStorage)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceSetup:
		return api.getSetupScript(query.Token)
	case common.ResourceInitYaml:
		return api.getInitYaml(query.Token, query.Node)
	default:
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "resource"),
			common.Field("name", resourceName))
	}
}

func (api *API) getInitYaml(token, edgeKubeNodeName string) ([]byte, error) {
	info, err := api.checkAndParseToken(token)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	switch common.Resource(info[InfoKind].(string)) {
	case common.Node:
		return api.initService.InitWithNode(info[InfoNamespace].(string), info[InfoName].(string), edgeKubeNodeName)
	default:
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
}

func (api *API) getSetupScript(token string) ([]byte, error) {
	sysConf, err := api.sysConfigService.GetSysConfig("address", common.AddressActive)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"Token":      token,
		"BaetylHost": sysConf.Value,
	}
	return api.ParseTemplate(common.ResourceSetup, params)
}

func (api *API) genCmd(kind, ns, name string) (string, error) {
	info := map[string]interface{}{
		InfoKind:      kind,
		InfoName:      name,
		InfoNamespace: ns,
		InfoExpiry:    CmdExpirationInSeconds,
		InfoTimestamp: time.Now().Unix(),
	}
	token, err := api.authService.GenToken(info)
	if err != nil {
		return "", err
	}
	host, err := api.sysConfigService.GetSysConfig("address", common.AddressActive)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`curl -skfL '%s/v1/active/setup.sh?token=%s' -osetup.sh && sh setup.sh`, host.Value, token), nil
}

func (api *API) checkAndParseToken(token string) (map[string]interface{}, error) {
	// check len
	if len(token) < 10 {
		return nil, ErrInvalidToken
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
	realToken, err := api.authService.GenToken(info)
	if err != nil {
		return nil, err
	}
	if realToken != token {
		return nil, ErrInvalidToken
	}

	expiry, ok := info[InfoExpiry].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	ts, ok := info[InfoTimestamp].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	// check expiration
	timestamp := time.Unix(int64(ts), 0)
	if timestamp.Add(time.Duration(int64(expiry))*time.Second).Unix() < time.Now().Unix() {
		return nil, ErrInvalidToken
	}
	return info, nil
}
