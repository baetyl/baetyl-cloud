package main

import (
	"runtime"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/server"
	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"

	_ "github.com/baetyl/baetyl-cloud/v2/plugin/awss3"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/database"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/auth"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/license"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/default/pki"
	_ "github.com/baetyl/baetyl-cloud/v2/plugin/kube"
	_ "github.com/go-sql-driver/mysql"
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
		config.SetPortFromEnv(&cfg)

		a, err := api.NewAPI(&cfg)
		if err != nil {
			return err
		}
		sa, err := api.NewSyncAPI(&cfg)
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

		ts, err := server.NewNodeServer(&cfg)
		if err != nil {
			return err
		}
		ts.SetSyncAPI(sa)
		ts.InitRoute()
		go ts.Run()
		defer ts.Close()
		ctx.Log().Info("node server starting")

		ctx.Wait()
		return nil
	})
}
