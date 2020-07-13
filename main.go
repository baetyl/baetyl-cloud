package main

import (
	"runtime"

	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/plugin"
	_ "github.com/baetyl/baetyl-cloud/plugin/awss3"
	_ "github.com/baetyl/baetyl-cloud/plugin/database"
	_ "github.com/baetyl/baetyl-cloud/plugin/default/auth"
	_ "github.com/baetyl/baetyl-cloud/plugin/default/license"
	_ "github.com/baetyl/baetyl-cloud/plugin/default/pki"
	_ "github.com/baetyl/baetyl-cloud/plugin/kube"
	"github.com/baetyl/baetyl-cloud/server"
	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/log"
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

		s, err := server.NewAdminServer(&cfg)
		if err != nil {
			return err
		}
		s.InitRoute()
		go s.Run()
		defer s.Close()
		ctx.Log().Info("admin server starting")

		ts, err := server.NewNodeServer(&cfg)
		if err != nil {
			return err
		}
		ts.InitRoute()
		go ts.Run()
		defer ts.Close()
		ctx.Log().Info("node server starting")

		as, err := server.NewActiveServer(&cfg)
		if err != nil {
			return err
		}
		as.InitRoute()
		go as.Run()
		defer as.Close()
		ctx.Log().Info("active server starting")

		ctx.Wait()
		return nil
	})
}
