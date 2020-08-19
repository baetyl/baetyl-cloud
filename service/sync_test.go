package service

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func genCallbackTestCase() *models.Callback {
	return &models.Callback{
		Name:        "c",
		Namespace:   "default",
		Method:      "Post",
		Description: "123",
		Params:      map[string]string{"a": "a"},
		Body:        map[string]string{"b": "b"},
	}
}

func TestReport(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	ns := ms.NewMockNodeService(mockObject.ctl)

	namespace := "ns01"
	name := "node01"

	ns.EXPECT().UpdateReport(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	sync := SyncServiceImpl{
		NodeService: ns,
	}
	info := specV1.Report{}
	response, err := sync.Report(namespace, name, info)
	assert.NotNil(t, err)

	shadow := &models.Shadow{
		Desire: specV1.Desire{
			common.DesiredApplications: []specV1.AppInfo{{
				"app",
				"v1",
			}},
		},
	}

	ns.EXPECT().UpdateReport(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	response, err = sync.Report(namespace, name, info)
	assert.Error(t, err)

	shadow = &models.Shadow{
		Desire: specV1.Desire{
			common.DesiredApplications: []specV1.AppInfo{{
				"app",
				"v1",
			}},
			common.DesiredSysApplications: []specV1.AppInfo{
				{
					"sysapp01",
					"v1",
				},
			},
		},
	}
	ns.EXPECT().UpdateReport(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil)
	response, err = sync.Report(namespace, name, info)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestSyncDesire(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()
	cs := ms.NewMockConfigService(mockObject.ctl)
	as := ms.NewMockApplicationService(mockObject.ctl)
	os := ms.NewMockObjectService(mockObject.ctl)
	sync := SyncServiceImpl{
		ConfigService: cs,
		AppService:    as,
		ObjectService: os,
		Hooks:         map[string]interface{}{},
	}
	sync.Hooks[MethodPopulateConfig] = HandlerPopulateConfig(sync.populateConfig)
	reqs := []specV1.ResourceInfo{
		{
			Kind:    specV1.KindApplication,
			Name:    "app",
			Version: "v1",
		},
		{
			Kind:    specV1.KindConfiguration,
			Name:    "config",
			Version: "v1",
		},
	}
	app := &specV1.Application{
		Name:    "app",
		Version: "v1",
	}
	obj1 := &specV1.ConfigurationObject{
		Metadata: map[string]string{
			"source": "source1",
			"bucket": "bucket1",
			"object": "object1",
			"userID": "namespace1",
		},
	}
	obj2 := &specV1.ConfigurationObject{
		MD5: "md52",
		URL: "url2",
	}
	obj1Data, err := json.Marshal(obj1)
	assert.NoError(t, err)
	obj2Data, err := json.Marshal(obj2)
	assert.NoError(t, err)

	data := map[string]string{
		common.ConfigObjectPrefix + "obj1": string(obj1Data),
		common.ConfigObjectPrefix + "obj2": string(obj2Data),
	}

	config := &specV1.Configuration{
		Name:    "config",
		Version: "v1",
		Data:    data,
	}
	objURL := &models.ObjectURL{
		URL:   "url1",
		Token: "token1",
	}

	namespace := "namespace1"
	param := models.ConfigObjectItem{
		Source: "source1",
		Bucket: "bucket1",
		Object: "object1",
	}
	as.EXPECT().Get(namespace, reqs[0].Name, reqs[0].Version).Return(app, nil).Times(1)
	cs.EXPECT().Get(namespace, reqs[1].Name, "").Return(config, nil).Times(1)
	os.EXPECT().GenObjectURL(namespace, param).Return(objURL, nil).Times(1)
	res, err := sync.Desire(namespace, reqs, map[string]string{})
	assert.NoError(t, err)

	resApp := res[0].Value.Value.(*specV1.Application)
	resConfig := res[1].Value.Value.(*specV1.Configuration)
	assert.Equal(t, resApp, app)
	assert.Equal(t, resConfig.Name, config.Name)
	assert.Equal(t, resConfig.Version, config.Version)
	assert.Len(t, resConfig.Data, 2)

	var resObj1 specV1.ConfigurationObject
	err = json.Unmarshal([]byte(resConfig.Data[common.ConfigObjectPrefix+"obj1"]), &resObj1)
	assert.NoError(t, err)
	assert.Equal(t, resObj1.MD5, obj1.MD5)
	assert.Equal(t, resObj1.URL, objURL.URL)

	var resObj2 specV1.ConfigurationObject
	err = json.Unmarshal([]byte(resConfig.Data[common.ConfigObjectPrefix+"obj2"]), &resObj2)
	assert.Equal(t, resObj2.MD5, obj2.MD5)
	assert.Equal(t, resObj2.URL, obj2.URL)
	assert.Empty(t, resObj2.Token)
}

func TestSyncService_Report(t *testing.T) {
	mockObject := InitMockEnvironment(t)
	defer mockObject.Close()

	mockNs := ms.NewMockNodeService(mockObject.ctl)
	ss := &SyncServiceImpl{
		NodeService: mockNs,
	}
	namespace := "namespace01"
	name := "name"
	report := specV1.Report{
		"time": time.Now(),
		"apps": []specV1.AppInfo{
			{
				Name:    "app01",
				Version: "v1",
			}, {
				Name:    "app02",
				Version: "v2",
			},
		},
		"appstats": []specV1.AppStats{
			{
				AppInfo: specV1.AppInfo{
					Name:    "app01",
					Version: "v1",
				},
				Status: "pending",
			},
		},
	}

	shadow := &models.Shadow{
		Desire: specV1.Desire{
			"apps": []interface{}{
				map[string]interface{}{
					"name":    "app01",
					"version": "v1",
				},
				map[string]interface{}{
					"name":    "app02",
					"version": "v2",
				},
				map[string]interface{}{
					"name":    "app03",
					"version": "v3",
				},
			},
		},
	}
	mockNs.EXPECT().UpdateReport(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))

	_, err := ss.Report(namespace, name, report)
	assert.NotNil(t, err)

	mockNs.EXPECT().UpdateReport(gomock.Any(), gomock.Any(), gomock.Any()).Return(shadow, nil).AnyTimes()

	_, err = ss.Report(namespace, name, report)
	assert.NotNil(t, err)
	report[common.DesiredSysApplications] = []specV1.AppInfo{
		{
			"sysapp01",
			"v1",
		},
	}
	_, err = ss.Report(namespace, name, report)
	assert.NotNil(t, err)
}

func TestDesireDiff(t *testing.T) {
	report := specV1.Report{
		"time": time.Now(),
		"apps": []specV1.AppInfo{
			{
				Name:    "app01",
				Version: "v1",
			}, {
				Name:    "app02",
				Version: "v2",
			},
		},
		"appstats": []specV1.AppStats{
			{
				AppInfo: specV1.AppInfo{
					Name:    "app01",
					Version: "v1",
				},
				Status: "pending",
			},
		},
	}

	desire := specV1.Desire{
		"apps": []specV1.AppInfo{
			{
				"app01",
				"v1",
			}, {
				"app02",
				"v2",
			}, {
				"app03",
				"v2",
			},
		},
	}

	isSysApp := false
	delta, _ := desire.Diff(report)
	assert.Equal(t, desire.AppInfos(isSysApp), delta.AppInfos(isSysApp))
}
