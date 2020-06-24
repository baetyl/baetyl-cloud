package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-go/errors"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"github.com/baetyl/baetyl-go/utils"
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

func (api *API) Active(c *common.Context) (interface{}, error) {
	info := &specV1.ActiveRequest{}
	err := c.LoadBody(info)
	if err != nil {
		err = common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	batch, err := api.checkBatch(info)
	if err != nil {
		return nil, err
	}
	record, err := api.checkRecord(batch, info.FingerprintValue)
	if err != nil {
		return nil, err
	}
	_, err = api.genNodeAndSysApp(record.Namespace, record.BatchName, record.NodeName)
	if err != nil {
		return nil, err
	}
	if err = api.activeAndCallback(record, info.PenetrateData, batch.CallbackName, c.ClientIP()); err != nil {
		return nil, err
	}
	cert, err := api.initService.GetSyncCert(record.Namespace, record.NodeName)
	if err != nil {
		return nil, err
	}

	return specV1.ActiveResponse{
		NodeName:  record.NodeName,
		Namespace: record.Namespace,
		Certificate: utils.Certificate{
			CA:                 string(cert.Data["ca.pem"]),
			Key:                string(cert.Data["client.key"]),
			Cert:               string(cert.Data["client.pem"]),
			InsecureSkipVerify: false,
		},
	}, nil
}

func (api *API) checkBatch(info *specV1.ActiveRequest) (*models.Batch, error) {
	batch, err := api.registerService.GetBatch(info.BatchName, info.Namespace)
	if err != nil {
		return nil, err
	}
	if string(batch.SecurityType) != info.SecurityType {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "SecurityType error"))
	}
	if batch.SecurityType == common.Token && batch.SecurityKey != info.SecurityValue {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", "SecurityValue error"))
	}
	return batch, nil
}

func (api *API) checkRecord(batch *models.Batch, fingerprintValue string) (*models.Record, error) {
	record, err := api.registerService.GetRecordByFingerprint(batch.Name, batch.Namespace, fingerprintValue)
	if err != nil {
		return nil, common.Error(common.ErrDatabase, common.Field("error", err))
	}
	if record != nil {
		if record.Active == common.Activated {
			return record, nil
		}
	} else {
		lowFv := strings.ToLower(fingerprintValue)
		if batch.EnableWhitelist == common.DisableWhitelist {
			r := &models.Record{
				Name:             lowFv,
				Namespace:        batch.Namespace,
				BatchName:        batch.Name,
				NodeName:         lowFv,
				FingerprintValue: fingerprintValue,
				Active:           common.Inactivated,
				ActiveTime:       time.Unix(common.DefaultActiveTime, 0),
			}
			if _, err = api.registerService.CreateRecord(r); err != nil {
				return nil, err
			}
		}
		record, err = api.registerService.GetRecordByFingerprint(batch.Name, batch.Namespace, fingerprintValue)
		if err != nil {
			return nil, err
		}
		if record == nil {
			return nil, common.Error(
				common.ErrResourceNotFound,
				common.Field("type", "record"),
				common.Field("name", lowFv),
				common.Field("namespace", batch.Namespace))
		}
	}
	return record, nil
}

func (api *API) genNodeAndSysApp(ns, batchName, nodeName string) ([]specV1.Application, error) {
	node, err := api.nodeService.Get(ns, nodeName)
	if err != nil {
		if e, ok := err.(errors.Coder); ok && e.Code() == common.ErrResourceNotFound {
			n := &specV1.Node{
				Name:      nodeName,
				Namespace: ns,
				Labels: map[string]string{
					common.LabelBatch:    batchName,
					common.LabelNodeName: nodeName,
				},
			}
			node, err = api.nodeService.Create(ns, n)
			if err != nil {
				return nil, common.Error(common.ErrK8S, common.Field("error", err))
			}
		} else {
			return nil, err
		}

	}

	apps, err := api.GenSysApp(node.Name, ns, SystemApps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (api *API) activeAndCallback(record *models.Record, data map[string]string, callName, ip string) error {
	record.Active = common.Activated
	record.ActiveIP = ip
	record.ActiveTime = time.Now()
	_, err := api.registerService.UpdateRecord(record)
	if err != nil {
		return err
	}
	if callName != "" {
		if _, err = api.callbackService.Callback(callName, record.Namespace, data); err != nil {
			return err
		}
	}
	return nil
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
	case common.Batch:
		batch, err := api.registerService.GetBatch(info[InfoName].(string), info[InfoNamespace].(string))
		if err != nil {
			return nil, err
		}
		return api.initService.InitWithBitch(batch, edgeKubeNodeName)
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
