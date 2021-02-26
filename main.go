package main

import (
	"runtime"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/server"
	"github.com/baetyl/baetyl-cloud/v2/task"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/baetyl/baetyl-cloud/v2/plugin/awss3"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/database"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/auth"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/license"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/lock"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/pki"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/kube"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/link/httplink"
)

func main() {
	defer plugin.ClosePlugins()
	runtime.GOMAXPROCS(runtime.NumCPU())

	context.Run(func(ctx context.Context) error {
		var cfg config.CloudConfig
		err := ctx.LoadCustomConfig(&cfg)
		if err != nil {
			return err
		}

		ctx.Log().Debug("cloud config", log.Any("cfg", cfg))

		common.SetConfFile(ctx.ConfFile())

		a, err := api.NewAPI(&cfg)
		if err != nil {
			return err
		}
		sa, err := api.NewSyncAPI(&cfg)
		if err != nil {
			return err
		}
		ia, err := api.NewInitAPI(&cfg)
		if err != nil {
			return err
		}
		s, err := server.NewAdminServer(&cfg)
		if err != nil {
			return err
		}
		s.SetAPI(a)
		s.InitRoute()
		go s.Run()
		defer s.Close()
		ctx.Log().Info("admin server starting")

		ss, err := server.NewSyncServer(&cfg)
		if err != nil {
			return err
		}
		ss.SetSyncAPI(sa)
		ss.InitMsgRouter()
		ss.Run()
		defer ss.Close()

		as, err := server.NewInitServer(&cfg)
		if err != nil {
			return err
		}
		as.SetAPI(ia)
		as.InitRoute()
		go as.Run()
		defer as.Close()
		ctx.Log().Info("init  server starting")
		tm, err := task.NewTaskManager(&cfg)
		if err != nil {
			return err
		}
		task.RegisterNamespaceProcessor(&cfg)
		tm.Start()
		ctx.Wait()
		return nil
	})
}
