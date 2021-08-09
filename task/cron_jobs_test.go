package task

import (
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/config"
	mp "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestLoadCronJobs(t *testing.T) {
	c := &config.CloudConfig{}
	c.CronJobs = make([]config.CronJob, 1)
	mCronJobs, mCtl := InitMockCronEnvironment(t)
	defer mCtl.Finish()
	cronEntity := cron.New()
	cronJob := &cronJobs{
		cronDb:  mCronJobs.sCronDb,
		locker:  mCronJobs.sLocker,
		node:    mCronJobs.sNode,
		app:     mCronJobs.sApp,
		index:   mCronJobs.sIndex,
		cron:    cronEntity,
		cronLog: log.L().With(log.Any("task", "cron")),
	}
	cronEntity.Start()
	defer cronEntity.Stop()
	c.CronJobs[0] = config.CronJob{
		CronName: "test",
	}
	err := cronJob.LoadCronJobs(c)
	assert.Error(t, err, ErrCronNotSupport)
	c.CronJobs[0] = config.CronJob{
		CronName: CronApp,
		CronGap: "10s",
	}
	err = cronJob.LoadCronJobs(c)
	assert.NoError(t, err)
}

func TestCronAppFunc(t *testing.T) {
	c := &config.CloudConfig{}
	c.CronJobs = make([]config.CronJob, 1)
	mCronJobs, mCtl := InitMockCronEnvironment(t)
	defer mCtl.Finish()
	cronEntity := cron.New()
	cronJob := &cronJobs{
		cronDb:  mCronJobs.sCronDb,
		locker:  mCronJobs.sLocker,
		node:    mCronJobs.sNode,
		app:     mCronJobs.sApp,
		index:   mCronJobs.sIndex,
		cron:    cronEntity,
		cronLog: log.L().With(log.Any("task", "cron")),
	}
	mCronJobs.sLocker.EXPECT().Lock(gomock.Any(), CronApp, gomock.Any()).Return("", nil).AnyTimes()
	mCronJobs.sLocker.EXPECT().Unlock(gomock.Any(), CronApp, "").Return().AnyTimes()
	mCronJobs.sCronDb.EXPECT().DeleteExpiredApps(gomock.Any()).Return(nil).AnyTimes()

	mCronJobs.sCronDb.EXPECT().ListExpiredApps().Return(nil, ErrCronNotSupport).Times(1)
	cronJob.CronAppFunc()

	name, ns, selector := "baetyl", "cloud", "baetyl-node-name=node1"

	cronApps := []models.Cron{
		{
			Namespace: ns,
			Name:      name,
			Selector:  selector,
			CronTime:  time.Now(),
		},
	}
	mCronJobs.sCronDb.EXPECT().ListExpiredApps().Return(cronApps, nil).AnyTimes()
	mCronJobs.sApp.EXPECT().GetApplication(ns, name, "").Return(nil, ErrCronNotSupport).Times(1)
	cronJob.CronAppFunc()

	app := &specV1.Application{
		Selector: "baetyl-node-name=node2",
		CronStatus: specV1.CronWait,
	}
	mCronJobs.sApp.EXPECT().GetApplication(ns, name, "").Return(app, nil).AnyTimes()
	mCronJobs.sApp.EXPECT().UpdateApplication(ns, app).Return(nil, ErrCronNotSupport).Times(1)
	cronJob.CronAppFunc()

	mCronJobs.sApp.EXPECT().UpdateApplication(ns, app).Return(app, nil).AnyTimes()
	mCronJobs.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, app).Return(nil, ErrCronNotSupport).Times(1)
	cronJob.CronAppFunc()

	mCronJobs.sNode.EXPECT().UpdateNodeAppVersion(nil, ns, app).Return(nil, nil).AnyTimes()
	mCronJobs.sIndex.EXPECT().RefreshNodesIndexByApp(nil, ns, name, gomock.Any()).Return(ErrCronNotSupport).Times(1)
	cronJob.CronAppFunc()

	mCronJobs.sIndex.EXPECT().RefreshNodesIndexByApp(nil, ns, name, gomock.Any()).Return(nil)
	cronJob.CronAppFunc()

	assert.Equal(t, app.CronStatus, specV1.CronFinished)
	assert.Equal(t, app.Selector, selector)
}

type MockCronJobs struct {
	sCronDb *mp.MockCron
	sLocker *mp.MockLocker
	sApp    *mp.MockApplication
	sNode   *ms.MockNodeService
	sIndex  *ms.MockIndexService
	Cron    *cron.Cron
	CronLog *log.Logger
}

func InitMockCronEnvironment(t *testing.T) (*MockCronJobs, *gomock.Controller) {
	mockCtl := gomock.NewController(t)
	return &MockCronJobs{
		sCronDb: mp.NewMockCron(mockCtl),
		sLocker: mp.NewMockLocker(mockCtl),
		sApp:    mp.NewMockApplication(mockCtl),
		sNode:   ms.NewMockNodeService(mockCtl),
		sIndex:  ms.NewMockIndexService(mockCtl),
	}, mockCtl
}
