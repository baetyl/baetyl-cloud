package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func mockCron(mock plugin.Cron) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func TestCronService(t *testing.T) {
	conf := &config.CloudConfig{}
	conf.Plugin.Cron = common.RandString(9)
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mCron := mockPlugin.NewMockCron(mockCtl)
	plugin.RegisterFactory(conf.Plugin.Cron, mockCron(mCron))

	cs, err := NewCronService(conf)
	assert.NoError(t, err)

	n, ns := "baetyl", "cloud"
	cronEntity := &models.Cron{
		Name: n,
		Namespace: ns,
		Selector: "baetyl-node-name=node1",
		CronTime: time.Now(),
	}

	mCron.EXPECT().GetCron(n, ns).Return(cronEntity, nil)
	_, err = cs.GetCron(n, ns)
	assert.NoError(t, err)

	mCron.EXPECT().CreateCron(cronEntity).Return(nil)
	err = cs.CreateCron(cronEntity)
	assert.NoError(t, err)

	mCron.EXPECT().UpdateCron(cronEntity).Return(nil)
	err = cs.UpdateCron(cronEntity)
	assert.NoError(t, err)

	mCron.EXPECT().DeleteCron(n, ns).Return(nil)
	err = cs.DeleteCron(n, ns)
	assert.NoError(t, err)

	mCron.EXPECT().ListExpiredApps().Return(nil, nil)
	_, err = cs.ListExpiredApps()
	assert.NoError(t, err)

	mCron.EXPECT().DeleteExpiredApps(nil).Return(nil)
	err = cs.DeleteExpiredApps(nil)
	assert.NoError(t, err)
}
