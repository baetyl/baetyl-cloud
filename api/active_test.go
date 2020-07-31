package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initActiveAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		active := v1.Group("/active")
		active.POST("", mockIM, common.Wrapper(api.Active))
		active.GET("/:resource", mockIM, common.WrapperRaw(api.GetResource))
	}
	return api, router, mockCtl
}

func TestAPI_GetResource(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	init := service.NewMockInitializeService(ctl)
	ss := service.NewMockSysConfigService(ctl)
	api.sysConfigService = ss
	api.initService = init

	// good case : metrics
	init.EXPECT().GetResource(common.ResourceMetrics).Return("metrics", nil).Times(1)

	req, _ := http.NewRequest(http.MethodGet, "/v1/active/"+common.ResourceMetrics, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// good case : local_path_storage
	init.EXPECT().GetResource(common.ResourceLocalPathStorage).Return("local-path-storage", nil).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/active/"+common.ResourceLocalPathStorage, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// good case : setup
	sc := &models.SysConfig{
		Type:  "address",
		Key:   common.AddressActive,
		Value: "baetyl.com",
	}
	ss.EXPECT().GetSysConfig(sc.Type, sc.Key).Return(sc, nil).Times(1)
	init.EXPECT().GetResource(common.ResourceSetup).Return("{}", nil).Times(1)

	req, _ = http.NewRequest(http.MethodGet, "/v1/active/"+common.ResourceSetup, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// bad case : init
	req, _ = http.NewRequest(http.MethodGet, "/v1/active/baetyl-init.yml", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// bad case : not found
	req, _ = http.NewRequest(http.MethodGet, "/v1/active/notfound", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestApi_genInitYaml(t *testing.T) {
	// expiry token
	token := "ac40cc632e217d7675abfdfbf64e285f7b22657870697279223a333630302c226b696e64223a226e6f6465222c226e616d65223a22303431353031222c226e616d657370616365223a2264656661756c74222c2274696d657374616d70223a313538363935363931367d"
	kube := "k3s"
	api, _, ctl := initActiveAPI(t)
	auth := service.NewMockAuthService(ctl)
	api.authService = auth
	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	res, err := api.getInitYaml(token, kube)
	assert.Error(t, err, ErrInvalidToken)
	assert.Nil(t, res)
}

func TestApi_genSetupScript(t *testing.T) {
	// bad case : sys config not found
	token := "ac40cc632e217d7675abfdfbf64e285f7b22657870697279223a333630302c226b696e64223a226e6f6465222c226e616d65223a22303431353031222c226e616d657370616365223a2264656661756c74222c2274696d657374616d70223a313538363935363931367d"
	api, _, ctl := initActiveAPI(t)
	ss := service.NewMockSysConfigService(ctl)
	api.sysConfigService = ss

	ss.EXPECT().GetSysConfig("address", common.AddressActive).Return(nil, fmt.Errorf("not found")).Times(1)

	res, err := api.getSetupScript(token)
	assert.Error(t, err, common.Error(common.ErrResourceNotFound,
		common.Field("type", "address"),
		common.Field("name", common.AddressActive)))
	assert.Nil(t, res)
}

func TestApi_genCmd(t *testing.T) {
	token := "ac40cc632e217d7675abfdfbf64e285f7b22657870697279223a333630302c226b696e64223a226e6f6465222c226e616d65223a22303431353031222c226e616d657370616365223a2264656661756c74222c2274696d657374616d70223a313538363935363931367d"
	api, _, ctl := initActiveAPI(t)
	ss := service.NewMockSysConfigService(ctl)
	auth := service.NewMockAuthService(ctl)
	api.authService = auth
	api.sysConfigService = ss

	// bad case 0: gen Token error
	auth.EXPECT().GenToken(gomock.Any()).Return("", fmt.Errorf("gen token err")).Times(1)

	res, err := api.genCmd("batch", "default", "test")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	// bad case 1: get sys config error
	auth.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	ss.EXPECT().GetSysConfig("address", common.AddressActive).Return(nil, fmt.Errorf("not found")).Times(1)

	res, err = api.genCmd("batch", "default", "test")
	assert.Error(t, err, common.Error(common.ErrResourceNotFound,
		common.Field("type", "address"),
		common.Field("name", common.AddressActive)))
	assert.Equal(t, "", res)
}

func TestAPI_Active_ErrBatch(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	rs := service.NewMockRegisterService(ctl)
	api.registerService = rs

	mBatch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		SecurityType: "Token",
		SecurityKey:  "123",
	}
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)

	info := &specV1.ActiveRequest{
		BatchName:        "test",
		Namespace:        "default",
		FingerprintValue: "123123",
		SecurityType:     "Token",
		SecurityValue:    "123",
		PenetrateData:    map[string]string{"a": "b"},
	}
	body, err := json.Marshal(info)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_Active_ErrRecord(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	rs := service.NewMockRegisterService(ctl)
	cs := service.NewMockCallbackService(ctl)
	ns := service.NewMockNodeService(ctl)
	api.callbackService = cs
	api.registerService = rs
	api.nodeService = ns

	mBatch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		SecurityType: "Token",
		SecurityKey:  "123",
	}
	mRecord := &models.Record{
		Name:             "r0",
		Namespace:        mBatch.Namespace,
		BatchName:        mBatch.Name,
		Active:           0,
		FingerprintValue: "123123",
		NodeName:         "123123",
	}
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(nil, nil).Times(2)
	rs.EXPECT().CreateRecord(gomock.Any()).Return(nil, nil).Times(1)

	info := &specV1.ActiveRequest{
		BatchName:        "test",
		Namespace:        "default",
		FingerprintValue: "123123",
		SecurityType:     "Token",
		SecurityValue:    "123",
		PenetrateData:    map[string]string{"a": "b"},
	}
	body, err := json.Marshal(info)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_Active_ErrNode(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	rs := service.NewMockRegisterService(ctl)
	cs := service.NewMockCallbackService(ctl)
	ns := service.NewMockNodeService(ctl)
	api.callbackService = cs
	api.registerService = rs
	api.nodeService = ns

	mBatch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		SecurityType: "Token",
		SecurityKey:  "123",
	}
	mRecord := &models.Record{
		Name:             "r0",
		Namespace:        mBatch.Namespace,
		BatchName:        mBatch.Name,
		FingerprintValue: "123123",
		NodeName:         "123123",
	}
	mNode := &specV1.Node{
		Name:      "123123",
		Namespace: mBatch.Namespace,
		Labels: map[string]string{
			common.LabelBatch:    "test",
			common.LabelNodeName: "123123",
		},
	}

	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(nil, nil).Times(1)
	ns.EXPECT().Create(mNode.Namespace, mNode).Return(nil, common.Error(common.ErrK8S, common.Field("error", "node error"))).Times(1)

	info := &specV1.ActiveRequest{
		BatchName:        "test",
		Namespace:        "default",
		FingerprintValue: "123123",
		SecurityType:     "Token",
		SecurityValue:    "123",
		PenetrateData:    map[string]string{"a": "b"},
	}
	body, err := json.Marshal(info)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAPI_Active_ErrSys(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	rs := service.NewMockRegisterService(ctl)
	cs := service.NewMockCallbackService(ctl)
	ns := service.NewMockNodeService(ctl)
	api.callbackService = cs
	api.registerService = rs
	api.nodeService = ns

	mBatch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		SecurityType: "Token",
		SecurityKey:  "123",
	}
	mRecord := &models.Record{
		Name:             "r0",
		Namespace:        mBatch.Namespace,
		BatchName:        mBatch.Name,
		FingerprintValue: "123123",
		NodeName:         "123123",
	}
	mNode := &specV1.Node{
		Name:      "123123",
		Namespace: mBatch.Namespace,
		Labels: map[string]string{
			common.LabelBatch:    "test",
			common.LabelNodeName: "123123",
		},
	}
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(nil, nil).Times(1)
	ns.EXPECT().Create(mNode.Namespace, mNode).Return(mNode, nil).Times(1)

	info := &specV1.ActiveRequest{
		BatchName:        "test",
		Namespace:        "default",
		FingerprintValue: "123123",
		SecurityType:     "Token",
		SecurityValue:    "123",
		PenetrateData:    map[string]string{"a": "b"},
	}
	body, err := json.Marshal(info)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAPI_Active_ErrSecret(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	rs := service.NewMockRegisterService(ctl)
	as := service.NewMockApplicationService(ctl)
	cs := service.NewMockCallbackService(ctl)
	ns := service.NewMockNodeService(ctl)
	ss := service.NewMockSecretService(ctl)
	ccs := service.NewMockConfigService(ctl)
	scs := service.NewMockSysConfigService(ctl)
	pki := service.NewMockPKIService(ctl)
	is := service.NewMockIndexService(ctl)
	init := service.NewMockInitializeService(ctl)

	api.callbackService = cs
	api.registerService = rs
	api.nodeService = ns
	api.configService = ccs
	api.applicationService = as
	api.secretService = ss
	api.sysConfigService = scs
	api.pkiService = pki
	api.indexService = is
	api.initService = init

	mBatch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		CallbackName: "123",
		SecurityType: "Token",
		SecurityKey:  "123",
	}
	mRecord := &models.Record{
		Name:             "r0",
		Namespace:        mBatch.Namespace,
		BatchName:        mBatch.Name,
		FingerprintValue: "123123",
		NodeName:         "123123",
	}
	mNode := &specV1.Node{
		Name:      "123123",
		Namespace: mBatch.Namespace,
		Labels: map[string]string{
			common.LabelBatch:    "test",
			common.LabelNodeName: "123123",
		},
	}

	conf := &specV1.Configuration{
		Name:      "testConf",
		Namespace: mNode.Namespace,
		Data:      nil,
	}
	app := &specV1.Application{
		Name:      "testApp",
		Namespace: mNode.Namespace,
	}
	nodeList := []string{"s0", "s1", "s2"}
	sysConf := &models.SysConfig{
		Type:  "baetyl-edge",
		Key:   "test",
		Value: "123",
	}
	certPEM := &models.PEMCredential{
		CertPEM: []byte("test"),
		KeyPEM:  []byte("test"),
	}
	certMap := map[string][]byte{
		"client.pem": certPEM.CertPEM,
		"client.key": certPEM.KeyPEM,
		"ca.pem":     []byte("test"),
	}
	secret := &specV1.Secret{
		Name:      "sync-" + mNode.Name + "-core",
		Namespace: mRecord.Namespace,
		Data:      certMap,
		Version:   "123",
	}
	ccs.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(conf, nil).Times(2)
	as.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(app, nil).Times(2)
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	rs.EXPECT().UpdateRecord(mRecord).Return(nil, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(nil, nil).Times(1)
	ns.EXPECT().Create(mNode.Namespace, mNode).Return(mNode, nil).Times(1)
	ss.EXPECT().Get(mRecord.Namespace, gomock.Any(), "").Return(nil, nil).AnyTimes()
	cs.EXPECT().Callback(mBatch.CallbackName, mBatch.Namespace, gomock.Any()).Return(nil, nil).Times(1)
	scs.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysConf, nil).AnyTimes()
	pki.EXPECT().SignClientCertificate(gomock.Any(), gomock.Any()).Return(certPEM, nil).AnyTimes()
	pki.EXPECT().GetCA().Return([]byte("test"), nil).AnyTimes()
	ss.EXPECT().Create(mRecord.Namespace, gomock.Any()).Return(secret, nil).AnyTimes()
	is.EXPECT().RefreshNodesIndexByApp(mRecord.Namespace, gomock.Any(), nodeList).Times(2)
	ns.EXPECT().UpdateNodeAppVersion(mRecord.Namespace, gomock.Any()).Return(nodeList, nil).Times(2)
	init.EXPECT().GetResource(gomock.Any()).Return("{}", nil).AnyTimes()
	init.EXPECT().GetSyncCert(mRecord.Namespace, mRecord.NodeName).Return(nil, fmt.Errorf("resource not found")).Times(1)

	info := &specV1.ActiveRequest{
		BatchName:        "test",
		Namespace:        "default",
		FingerprintValue: "123123",
		SecurityType:     "Token",
		SecurityValue:    "123",
		PenetrateData:    map[string]string{"a": "b"},
	}
	body, err := json.Marshal(info)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAPI_Active(t *testing.T) {
	api, router, ctl := initActiveAPI(t)
	rs := service.NewMockRegisterService(ctl)
	as := service.NewMockApplicationService(ctl)
	cs := service.NewMockCallbackService(ctl)
	ns := service.NewMockNodeService(ctl)
	ss := service.NewMockSecretService(ctl)
	ccs := service.NewMockConfigService(ctl)
	scs := service.NewMockSysConfigService(ctl)
	pki := service.NewMockPKIService(ctl)
	is := service.NewMockIndexService(ctl)
	init := service.NewMockInitializeService(ctl)
	api.callbackService = cs
	api.registerService = rs
	api.nodeService = ns
	api.configService = ccs
	api.applicationService = as
	api.secretService = ss
	api.sysConfigService = scs
	api.pkiService = pki
	api.indexService = is
	api.initService = init

	mBatch := &models.Batch{
		Name:         "test",
		Namespace:    "default",
		CallbackName: "123",
		SecurityType: "Token",
		SecurityKey:  "123",
	}
	mRecord := &models.Record{
		Name:             "r0",
		Namespace:        mBatch.Namespace,
		BatchName:        mBatch.Name,
		FingerprintValue: "123123",
		NodeName:         "123123",
	}
	mNode := &specV1.Node{
		Name:      "123123",
		Namespace: mBatch.Namespace,
		Labels: map[string]string{
			common.LabelBatch:    "test",
			common.LabelNodeName: "123123",
		},
	}

	conf := &specV1.Configuration{
		Name:      "testConf",
		Namespace: mNode.Namespace,
		Data:      nil,
	}
	app := &specV1.Application{
		Name:      "core-" + mNode.Name,
		Namespace: mNode.Namespace,
	}
	nodeList := []string{"s0", "s1", "s2"}
	sysConf := &models.SysConfig{
		Type:  "baetyl-edge",
		Key:   "test",
		Value: "123",
	}
	certPEM := &models.PEMCredential{
		CertPEM: []byte("test"),
		KeyPEM:  []byte("test"),
	}
	certMap := map[string][]byte{
		"client.pem": certPEM.CertPEM,
		"client.key": certPEM.KeyPEM,
		"ca.pem":     []byte("test"),
	}
	secret := &specV1.Secret{
		Name:      "sync-" + mNode.Name + "-core",
		Namespace: mRecord.Namespace,
		Data:      certMap,
		Version:   "123",
	}
	ccs.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(conf, nil).Times(2)
	as.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(app, nil).Times(2)
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	rs.EXPECT().UpdateRecord(mRecord).Return(nil, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(nil, fmt.Errorf("error")).Times(1)

	info := &specV1.ActiveRequest{
		BatchName:        "test",
		Namespace:        "default",
		FingerprintValue: "123123",
		SecurityType:     "Token",
		SecurityValue:    "123",
		PenetrateData:    map[string]string{"a": "b"},
	}

	body, err := json.Marshal(info)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	ccs.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(conf, nil).Times(2)
	as.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(app, nil).Times(2)
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	rs.EXPECT().UpdateRecord(mRecord).Return(nil, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	ns.EXPECT().Create(mNode.Namespace, mNode).Return(nil, fmt.Errorf("error")).Times(1)

	body, err = json.Marshal(info)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	ccs.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(conf, nil).Times(3)
	as.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(app, nil).Times(3)
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	rs.EXPECT().UpdateRecord(mRecord).Return(nil, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(nil, common.Error(common.ErrResourceNotFound)).Times(1)
	ns.EXPECT().Create(mNode.Namespace, mNode).Return(mNode, nil).Times(1)
	ss.EXPECT().Get(mRecord.Namespace, gomock.Any(), "").Return(secret, nil).AnyTimes()
	ss.EXPECT().Create(mRecord.Namespace, gomock.Any()).Return(secret, nil).AnyTimes()
	cs.EXPECT().Callback(mBatch.CallbackName, mBatch.Namespace, gomock.Any()).Return(nil, nil).Times(1)
	scs.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysConf, nil).AnyTimes()
	pki.EXPECT().SignClientCertificate(gomock.Any(), gomock.Any()).Return(certPEM, nil).AnyTimes()
	pki.EXPECT().GetCA().Return([]byte("test"), nil).AnyTimes()
	ns.EXPECT().UpdateNodeAppVersion(mRecord.Namespace, gomock.Any()).Return(nodeList, nil).Times(3)
	is.EXPECT().RefreshNodesIndexByApp(mRecord.Namespace, gomock.Any(), nodeList).Times(3)
	init.EXPECT().GetResource(gomock.Any()).Return("{}", nil).AnyTimes()
	init.EXPECT().GetSyncCert(mRecord.Namespace, mRecord.NodeName).Return(secret, nil).Times(1)

	body, err = json.Marshal(info)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	ccs.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(conf, nil).Times(3)
	as.EXPECT().Create(mNode.Namespace, gomock.Any()).Return(app, nil).Times(3)
	rs.EXPECT().GetBatch(mBatch.Name, mBatch.Namespace).Return(mBatch, nil).Times(1)
	rs.EXPECT().GetRecordByFingerprint(mBatch.Name, mBatch.Namespace, mRecord.FingerprintValue).Return(mRecord, nil).Times(1)
	rs.EXPECT().UpdateRecord(mRecord).Return(nil, nil).Times(1)
	ns.EXPECT().Get(mRecord.Namespace, mRecord.NodeName).Return(mNode, nil).Times(1)
	ns.EXPECT().Create(mNode.Namespace, mNode).Return(mNode, nil).Times(1)
	ss.EXPECT().Get(mRecord.Namespace, gomock.Any(), "").Return(secret, nil).AnyTimes()
	ss.EXPECT().Create(mRecord.Namespace, gomock.Any()).Return(secret, nil).AnyTimes()
	cs.EXPECT().Callback(mBatch.CallbackName, mBatch.Namespace, gomock.Any()).Return(nil, nil).Times(1)
	scs.EXPECT().GetSysConfig(gomock.Any(), gomock.Any()).Return(sysConf, nil).AnyTimes()
	pki.EXPECT().SignClientCertificate(gomock.Any(), gomock.Any()).Return(certPEM, nil).AnyTimes()
	pki.EXPECT().GetCA().Return([]byte("test"), nil).AnyTimes()
	ns.EXPECT().UpdateNodeAppVersion(mRecord.Namespace, gomock.Any()).Return(nodeList, nil).Times(3)
	is.EXPECT().RefreshNodesIndexByApp(mRecord.Namespace, gomock.Any(), nodeList).Times(3)
	init.EXPECT().GetResource(gomock.Any()).Return("{}", nil).AnyTimes()
	init.EXPECT().GetSyncCert(mRecord.Namespace, mRecord.NodeName).Return(secret, nil).Times(1)

	body, err = json.Marshal(info)
	assert.NoError(t, err)
	req, _ = http.NewRequest(http.MethodPost, "/v1/active", bytes.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_getInitYaml(t *testing.T) {
	api, _, ctl := initActiveAPI(t)
	as := service.NewMockAuthService(ctl)
	rs := service.NewMockRegisterService(ctl)
	init := service.NewMockInitializeService(ctl)
	api.authService = as
	api.initService = init
	api.registerService = rs

	info := map[string]interface{}{
		InfoKind:      "node",
		InfoName:      "n0",
		InfoNamespace: "default",
		InfoTimestamp: time.Now().Unix(),
		InfoExpiry:    60 * 60 * 24 * 3650,
	}
	data, err := json.Marshal(info)
	assert.NoError(t, err)
	encode := hex.EncodeToString(data)
	sign := "0123456789"
	token := sign + encode
	kube := "k3s"

	// good case 0
	as.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	init.EXPECT().InitWithNode("default", "n0", kube).Return(nil, nil).Times(1)
	_, err = api.getInitYaml(token, kube)
	assert.NoError(t, err)

	// good case 1
	info[InfoKind] = "batch"
	data, err = json.Marshal(info)
	assert.NoError(t, err)
	encode = hex.EncodeToString(data)
	token = sign + encode

	b := &models.Batch{}

	as.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	rs.EXPECT().GetBatch("n0", "default").Return(b, nil).Times(1)
	init.EXPECT().InitWithBitch(b, kube).Return(nil, nil).Times(1)
	_, err = api.getInitYaml(token, kube)
	assert.NoError(t, err)

	// bad case 0
	info[InfoKind] = "error"
	data, err = json.Marshal(info)
	assert.NoError(t, err)
	encode = hex.EncodeToString(data)
	token = sign + encode
	as.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	_, err = api.getInitYaml(token, kube)
	assert.Error(t, err)
}
