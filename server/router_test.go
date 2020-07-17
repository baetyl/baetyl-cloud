package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/mock/plugin"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-cloud/plugin"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func InitMockEnvironment(t *testing.T) (*AdminServer, *NodeServer, *ActiveServer,
	*mockPlugin.MockAuth, *mockPlugin.MockLicense, *gomock.Controller, *config.CloudConfig, *MisServer) {
	c := &config.CloudConfig{}
	c.Plugin.Auth = common.RandString(9)
	c.Plugin.ModelStorage = common.RandString(9)
	c.Plugin.DatabaseStorage = common.RandString(9)
	c.Plugin.Shadow = c.Plugin.DatabaseStorage
	c.Plugin.Objects = []string{common.RandString(9)}
	c.Plugin.PKI = common.RandString(9)
	c.Plugin.Functions = []string{common.RandString(9)}
	c.Plugin.License = common.RandString(9)
	c.Plugin.CacheStorage = common.RandString(9)
	c.NodeServer.Certificate.CA = ""
	c.NodeServer.Certificate.Cert = ""
	c.NodeServer.Certificate.Key = ""
	c.ActiveServer.Certificate.CA = "../test/cloud/ca.pem"
	c.ActiveServer.Certificate.Cert = "../test/cloud/server.pem"
	c.ActiveServer.Certificate.Key = "../test/cloud/server.key"
	mockCtl := gomock.NewController(t)

	mockModelStorage := mockPlugin.NewMockModelStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.ModelStorage, func() (plugin.Plugin, error) {
		return mockModelStorage, nil
	})
	res := &models.ConfigurationList{}
	mockModelStorage.EXPECT().ListConfig(gomock.Any(), gomock.Any()).Return(res, nil).AnyTimes()
	mockDBStorage := mockPlugin.NewMockDBStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.DatabaseStorage, func() (plugin.Plugin, error) {
		return mockDBStorage, nil
	})
	mockObjectStorage := mockPlugin.NewMockObject(mockCtl)
	for _, v := range c.Plugin.Objects {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockObjectStorage, nil
		})
	}
	mockFunction := mockPlugin.NewMockFunction(mockCtl)
	for _, v := range c.Plugin.Functions {
		plugin.RegisterFactory(v, func() (plugin.Plugin, error) {
			return mockFunction, nil
		})
	}

	mockAuth := mockPlugin.NewMockAuth(mockCtl)
	plugin.RegisterFactory(c.Plugin.Auth, func() (plugin.Plugin, error) {
		return mockAuth, nil
	})

	mPKI := mockPlugin.NewMockPKI(mockCtl)
	plugin.RegisterFactory(c.Plugin.PKI, func() (plugin.Plugin, error) {
		return mPKI, nil
	})

	mLicense := mockPlugin.NewMockLicense(mockCtl)
	plugin.RegisterFactory(c.Plugin.License, func() (plugin.Plugin, error) {
		return mLicense, nil
	})
	mockCacheStorage := mockPlugin.NewMockCacheStorage(mockCtl)
	plugin.RegisterFactory(c.Plugin.CacheStorage, func()(plugin.Plugin, error){
		return mockCacheStorage, nil
	})
	s, _ := NewAdminServer(c)
	n, _ := NewNodeServer(c)
	a, _ := NewActiveServer(c)
	m, _ := NewMisServer(c)

	return s, n, a, mockAuth, mLicense, mockCtl, c, m
}

func TestHandler(t *testing.T) {
	s, _, _, mkAuth, mkLicense, mockCtl, _, _ := InitMockEnvironment(t)
	defer mockCtl.Finish()

	s.InitRoute()
	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	// 200
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	// 200
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs", nil)
	w2 := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)

	//mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	// 400
	mConf := &specV1.Configuration{
		Namespace: "default",
		Name:      "abc",
		Labels:    make(map[string]string),
		Data:      make(map[string]string),
	}
	body, _ := json.Marshal(mConf)
	assert.Equal(t, http.StatusOK, w.Code)
	req, _ = http.NewRequest(http.MethodHead, "/zzz", bytes.NewReader(body))
	w3 := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w3, req)
	assert.Equal(t, http.StatusNotFound, w3.Code)

	// 401
	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(fmt.Errorf("err"))
	req, _ = http.NewRequest(http.MethodGet, "/v1/configs", nil)
	w4 := httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w4, req)
	assert.Equal(t, http.StatusUnauthorized, w4.Code)

	// 401
	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil)
	mkLicense.EXPECT().CheckQuota(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	req, _ = http.NewRequest(http.MethodPost, "/v1/nodes", nil)
	w4 = httptest.NewRecorder()
	s.GetRoute().ServeHTTP(w4, req)
	assert.Equal(t, http.StatusInternalServerError, w4.Code)

	go s.Run()
	defer s.Close()
}

func TestHandler_Node(t *testing.T) {
	t.Skip()
	_, n, _, mkAuth, _, mockCtl, c, _ := InitMockEnvironment(t)
	defer mockCtl.Finish()

	n.InitRoute()
	n.GetRoute().GET("/device", func(c *gin.Context) {
		cc := common.NewContext(c)
		c.JSON(common.PackageResponse(&struct {
			Namespace string
			Name      string
		}{
			Namespace: cc.GetNamespace(),
			Name:      cc.GetName(),
		}))
	})

	mkAuth.EXPECT().Authenticate(gomock.Any()).Return(nil).AnyTimes()

	// https 200
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	n.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go n.Run()
	defer n.Close()

	s2 := httptest.NewUnstartedServer(n.GetRoute())
	pool := x509.NewCertPool()
	caCrt, _ := ioutil.ReadFile("../test/cloud/ca.pem")
	pool.AppendCertsFromPEM(caCrt)
	serverCert, _ := ioutil.ReadFile("../test/cloud/server.pem")
	serverKey, _ := ioutil.ReadFile("../test/cloud/server.key")
	cert, _ := tls.X509KeyPair(serverCert, serverKey)

	s2.TLS = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientCAs:          pool,
		InsecureSkipVerify: true,
		ClientAuth:         tls.VerifyClientCertIfGiven, // IfGiven -> report else -> active
	}
	s2.StartTLS()

	poolca := x509.NewCertPool()
	caCrtca, _ := ioutil.ReadFile("../test/node/ca.pem")
	poolca.AppendCertsFromPEM(caCrtca)

	clientCertPEM, _ := ioutil.ReadFile("../test/node/client.pem")
	clientKeyPEM, _ := ioutil.ReadFile("../test/node/client.key")
	clientTLSCert, _ := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            poolca,
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{clientTLSCert}, //提供客户端的证书
			},
		},
	}
	_, err := client.Get(s2.URL + "/device")
	assert.NoError(t, err)

	// http 200
	c.NodeServer.Certificate.Key = ""
	c.NodeServer.Certificate.Cert = ""
	nHttp, err := NewNodeServer(c)
	assert.NoError(t, err)
	nHttp.InitRoute()
	req, err = http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	nHttp.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go nHttp.Run()
	defer nHttp.Close()
}

func TestHandler_Active(t *testing.T) {
	t.Skip()
	_, _, a, _, _, mockCtl, c, _ := InitMockEnvironment(t)
	defer mockCtl.Finish()

	// https 200
	a.InitRoute()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	a.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go a.Run()
	defer a.Close()
	// http 200
	c.ActiveServer.Certificate.Key = ""
	c.ActiveServer.Certificate.Cert = ""
	aHttp, err := NewActiveServer(c)
	assert.NoError(t, err)
	aHttp.InitRoute()
	req, err = http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	aHttp.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go aHttp.Run()
	defer aHttp.Close()
}

func TestHandler_Mis(t *testing.T){
	t.Skip()
	_, _, _, _, _, mockCtl, c, m := InitMockEnvironment(t)
	defer mockCtl.Finish()

	// https 200
	m.InitRoute()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	m.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go m.Run()
	defer m.Close()
	// http 200
	c.ActiveServer.Certificate.Key = ""
	c.ActiveServer.Certificate.Cert = ""
	mHttp, err := NewMisServer(c)
	assert.NoError(t, err)
	mHttp.InitRoute()
	req, err = http.NewRequest(http.MethodGet, "/health", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	mHttp.GetRoute().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	go mHttp.Run()
	defer mHttp.Close()
}
