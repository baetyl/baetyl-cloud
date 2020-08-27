package service

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"

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
)
// ActiveService
type ActiveService interface {
	GetResource(c *common.Context) (interface{}, error)
}

type ActiveServiceImpl struct {
	InitService   InitializeService
	SysCfgService SysConfigService
	AuthService   AuthService
	NodeService   NodeService
}

// NewSyncService new SyncService
func NewActiveService(config *config.CloudConfig) (ActiveService, error) {
	initService, err := NewInitializeService(config)
	if err != nil {
		return nil, err
	}
	sysCfgService, err := NewSysConfigService(config)
	if err != nil {
		return nil, err
	}
	authService, err := NewAuthService(config)
	if err != nil {
		return nil, err
	}
	nodeService, err := NewNodeService(config)
	if err != nil {
		return nil, err
	}
	return &ActiveServiceImpl{
		InitService:	initService,
		SysCfgService:	sysCfgService,
		AuthService:	authService,
		NodeService:    nodeService,
	}, nil
}

func (as *ActiveServiceImpl) GetResource(c *common.Context) (interface{}, error) {
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
		res, err := as.InitService.GetResource(common.ResourceMetrics)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceLocalPathStorage:
		res, err := as.InitService.GetResource(common.ResourceLocalPathStorage)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case common.ResourceSetup:
		return as.GetSetupScript(query.Token)
	case common.ResourceInitYaml:
		return as.getInitYaml(query.Token, query.Node)
	default:
		return nil, common.Error(
			common.ErrResourceNotFound,
			common.Field("type", "resource"),
			common.Field("name", resourceName))
	}
}

func (as *ActiveServiceImpl) getInitYaml(token, edgeKubeNodeName string) ([]byte, error) {
	info, err := as.CheckAndParseToken(token)
	if err != nil {
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
	switch common.Resource(info[InfoKind].(string)) {
	case common.Node:
		return as.InitService.InitWithNode(info[InfoNamespace].(string), info[InfoName].(string), edgeKubeNodeName)
	default:
		return nil, common.Error(
			common.ErrRequestParamInvalid,
			common.Field("error", err))
	}
}

func (as *ActiveServiceImpl) GetSetupScript(token string) ([]byte, error) {
	sysConf, err := as.SysCfgService.GetSysConfig("address", common.AddressActive)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"Token":      token,
		"BaetylHost": sysConf.Value,
	}
	return as.ParseTemplate(common.ResourceSetup, params)
}

func (as *ActiveServiceImpl) CheckAndParseToken(token string) (map[string]interface{}, error) {
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
	realToken, err := as.AuthService.GenToken(info)
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

func (as *ActiveServiceImpl) ParseTemplate(key string, data map[string]string) ([]byte, error) {
	tl, err := as.InitService.GetResource(key)
	if err != nil {
		return nil, err
	}
	t, err := template.New(key).Parse(tl)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
