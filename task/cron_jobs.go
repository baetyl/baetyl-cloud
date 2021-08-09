package task

import (
	"context"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/robfig/cron/v3"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

const (
	CronApp = "cronAppJob"
)

var (
	ErrCronNotSupport = errors.New("failed to add unsupported cron.")
)

type CronJobs interface {
	LoadCronJobs(cfg *config.CloudConfig) error
	Start()
	Stop()
}

type cronJobs struct {
	cronDb  plugin.Cron
	locker  plugin.Locker
	app     plugin.Application
	node    service.NodeService
	index   service.IndexService
	cron    *cron.Cron
	cronLog *log.Logger
}

func NewCronJobs(cronEntity *cron.Cron, cfg *config.CloudConfig) (CronJobs, error) {
	cronDb, err := plugin.GetPlugin("database")
	if err != nil {
		return nil, errors.Trace(err)
	}
	locker, err := plugin.GetPlugin(cfg.Plugin.Locker)
	app, err := plugin.GetPlugin(cfg.Plugin.Resource)
	node, err := service.NewNodeService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	index, err := service.NewIndexService(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	cronLog := log.L().With(log.Any("task", "cron"))
	return &cronJobs{
		cronDb: cronDb.(plugin.Cron),
		locker: locker.(plugin.Locker),
		app: app.(plugin.Application),
		node: node,
		index: index,
		cron: cronEntity,
		cronLog: cronLog,
	}, nil
}

func(c *cronJobs) CronAppFunc() {
	ctx := context.Background()
	version, err := c.locker.Lock(ctx, CronApp, 0)
	if err != nil {
		c.cronLog.Error("failed to get distributed lock", log.Any("error", err))
		return
	}
	defer c.locker.Unlock(ctx, CronApp, version)
	cronAppList, err := c.cronDb.ListExpiredApps()
	if err != nil {
		c.cronLog.Error("failed to list apps", log.Any("error", err))
		return
	}
	ids := make([]uint64, 0)
	for _, cronApp := range cronAppList {
		application, err := c.app.GetApplication(cronApp.Namespace, cronApp.Name, "")
		if err != nil {
			c.cronLog.Error("failed to get application", log.Any("error", err))
			continue
		}
		application.Selector = cronApp.Selector
		application.CronStatus = specV1.CronFinished
		application, err = c.app.UpdateApplication(cronApp.Namespace, application)
		if err != nil {
			c.cronLog.Error("failed to update application", log.Any("error", err))
			continue
		}
		nodes, err := c.node.UpdateNodeAppVersion(nil, cronApp.Namespace, application)
		if err != nil {
			c.cronLog.Error("failed to update shadow", log.Any("error", err))
			continue
		}
		err = c.index.RefreshNodesIndexByApp(nil, cronApp.Namespace, cronApp.Name, nodes)
		if err != nil {
			c.cronLog.Error("failed to refresh index", log.Any("error", err))
			continue
		}
		c.cronLog.Debug("success set app", log.Any("name", cronApp.Name))
		ids = append(ids, cronApp.Id)
	}
	if len(ids) > 0 {
		err = c.cronDb.DeleteExpiredApps(ids)
		if err != nil {
			c.cronLog.Error("failed to delete expired apps", log.Any("error", err))
			return
		}
	}
}

func(c *cronJobs) LoadCronJobs(cfg *config.CloudConfig) error {
	for _, cronJob := range cfg.CronJobs {
		switch cronJob.CronName {
		case CronApp:
			c.cron.AddFunc("@every " + cronJob.CronGap, c.CronAppFunc)
		default:
			return ErrCronNotSupport
		}
	}
	return nil
}

func(c *cronJobs) Start() {
	c.cron.Start()
}

func(c *cronJobs) Stop()  {
	c.cron.Stop()
}
