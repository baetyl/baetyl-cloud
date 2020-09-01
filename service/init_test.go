package service

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/mock/service"
)

func TestAPI_GetResource(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	tp := service.NewMockTemplateService(mockCtl)
	ns := service.NewMockNodeService(mockCtl)
	aus := service.NewMockAuthService(mockCtl)
	as.TemplateService = tp
	as.NodeService = ns
	as.AuthService = aus

	// good case : metrics
	tp.EXPECT().GetTemplate(common.ResourceMetrics).Return("metrics", nil).Times(1)
	res, _ := as.GetResource(common.ResourceMetrics, "", "")
	assert.Equal(t, res, []byte("metrics"))

	// good case : local_path_storage
	tp.EXPECT().GetTemplate(common.ResourceLocalPathStorage).Return("local-path-storage", nil).Times(1)
	res, _ = as.GetResource(common.ResourceLocalPathStorage, "", "")
	assert.Equal(t, res, []byte("local-path-storage"))

	// good case : setup
	tp.EXPECT().GenSetupShell(gomock.Any()).Return([]byte("shell"), nil).Times(1)
	res, _ = as.GetResource(common.ResourceSetup, "", "")
	assert.Equal(t, res, []byte("shell"))

	// bad case : not found
	_, err := as.GetResource("", "", "")
	assert.Error(t, err)
}

func TestApi_genInitYaml(t *testing.T) {
	// expiry token
	token := "ac40cc632e217d7675abfdfbf64e285f7b22657870697279223a333630302c226b696e64223a226e6f6465222c226e616d65223a22303431353031222c226e616d657370616365223a2264656661756c74222c2274696d657374616d70223a313538363935363931367d"
	kube := "k3s"
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	auth := service.NewMockAuthService(mockCtl)
	as.AuthService = auth
	res, err := as.getInitYaml(token, kube)
	assert.Error(t, err, ErrInvalidToken)
	assert.Nil(t, res)
}

func TestAPI_getInitYaml(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	aus := service.NewMockAuthService(mockCtl)
	as.AuthService = aus

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


	// bad case 0
	info[InfoKind] = "error"
	data, err = json.Marshal(info)
	assert.NoError(t, err)
	encode = hex.EncodeToString(data)
	token = sign + encode
	aus.EXPECT().GenToken(gomock.Any()).Return(token, nil).Times(1)
	_, err = as.getInitYaml(token, kube)
	assert.Error(t, err)
}
