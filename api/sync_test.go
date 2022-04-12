package api

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/config"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
)

func TestNewSyncAPI(t *testing.T) {
	// bad case
	_, err := NewSyncAPI(&config.CloudConfig{})
	assert.Error(t, err)
}

func TestSyncAPIImpl_Report(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	sync := &SyncAPIImpl{}
	mSync := ms.NewMockSyncService(mockCtl)
	sync.Sync = mSync
	sync.log = log.L().With(log.Any("test", "sync"))

	// good case 0
	info := specV1.Report{
		"apps": []specV1.AppInfo{
			{
				Name:    "app01",
				Version: "v1",
			},
		},
	}
	msg := specV1.Message{
		Kind:     specV1.MessageReport,
		Metadata: map[string]string{"name": "test", "namespace": "default"},
		Content:  specV1.LazyValue{},
	}
	bt, err := json.Marshal(info)
	assert.NoError(t, err)
	err = msg.Content.UnmarshalJSON(bt)
	assert.NoError(t, err)
	resp := specV1.Delta{}
	expMsg := &specV1.Message{
		Kind:     msg.Kind,
		Metadata: msg.Metadata,
		Content:  specV1.LazyValue{},
	}
	mSync.EXPECT().Report("default", "test", gomock.Any()).Return(resp, nil).Times(1)
	res, err := sync.Report(msg)
	assert.NoError(t, err)
	assert.EqualValues(t, expMsg.Kind, res.Kind)
	assert.EqualValues(t, expMsg.Metadata, res.Metadata)

	// bad case 0
	mSync.EXPECT().Report("default", "test", gomock.Any()).Return(nil, os.ErrInvalid).Times(1)
	_, err = sync.Report(msg)
	assert.Error(t, err)
}

func TestSyncAPIImpl_Desire(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	sync := &SyncAPIImpl{}
	mSync := ms.NewMockSyncService(mockCtl)
	sync.Sync = mSync

	// good case 0
	infos := specV1.DesireRequest{}
	msg := specV1.Message{
		Kind:     specV1.MessageDesire,
		Metadata: map[string]string{"namespace": "default"},
		Content:  specV1.LazyValue{},
	}
	bt, err := json.Marshal(infos)
	assert.NoError(t, err)
	err = msg.Content.UnmarshalJSON(bt)
	assert.NoError(t, err)
	resp := []specV1.ResourceValue{}
	expMsg := &specV1.Message{
		Kind:     msg.Kind,
		Metadata: msg.Metadata,
		Content:  specV1.LazyValue{},
	}
	mSync.EXPECT().Desire("default", nil, msg.Metadata).Return(resp, nil).Times(1)
	res, err := sync.Desire(msg)
	assert.NoError(t, err)
	assert.EqualValues(t, expMsg.Kind, res.Kind)
	assert.EqualValues(t, expMsg.Metadata, res.Metadata)

	// bad case 0
	mSync.EXPECT().Desire("default", nil, msg.Metadata).Return(nil, os.ErrInvalid).Times(1)
	_, err = sync.Desire(msg)
	assert.Error(t, err)
}

func TestSyncAPIImpl_updateAndroidInfo(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	sync := &SyncAPIImpl{}
	mSync := ms.NewMockSyncService(mockCtl)
	mNode := ms.NewMockNodeService(mockCtl)
	sync.Sync = mSync
	sync.Node = mNode
	sync.log = log.L().With(log.Any("test", "sync"))

	// good case 0
	info := specV1.Report{}
	info["node"] = map[string]interface{}{"1001": nil}

	node := &specV1.Node{
		Name:      "n0",
		Namespace: "default",
		Attributes: map[string]interface{}{
			context.RunModeAndroid: "1002",
		},
	}

	mNode.EXPECT().Update(node.Namespace, node).Return(node, nil).Times(1)
	err := sync.updateAndroidInfo(node, &info)
	assert.NoError(t, err)
}
