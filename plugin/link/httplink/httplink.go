package httplink

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/server"
)

type HTTPLink struct {
	Cfg       *CloudConfig
	Router    *gin.Engine
	Svr       *http.Server
	MsgRouter map[string]interface{}
}

func init() {
	plugin.RegisterFactory("httplink", NewHTTPLink)
}

func NewHTTPLink() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}

	router := gin.New()
	svr := &http.Server{
		Addr:           cfg.HTTPLink.Port,
		Handler:        router,
		ReadTimeout:    cfg.HTTPLink.ReadTimeout,
		WriteTimeout:   cfg.HTTPLink.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.HTTPLink.Certificate.Cert != "" &&
		cfg.HTTPLink.Certificate.Key != "" &&
		cfg.HTTPLink.Certificate.CA != "" {
		t, err := utils.NewTLSConfigServer(utils.Certificate{
			CA:             cfg.HTTPLink.Certificate.CA,
			Cert:           cfg.HTTPLink.Certificate.Cert,
			Key:            cfg.HTTPLink.Certificate.Key,
			ClientAuthType: tls.RequireAnyClientCert,
		})
		if err != nil {
			return nil, err
		}
		svr.TLSConfig = t
	}

	if svr.TLSConfig == nil {
		server.HeaderCommonName = cfg.HTTPLink.CommonName
		router.Use(server.ExtractNodeCommonNameFromHeader)
	} else {
		router.Use(server.ExtractNodeCommonNameFromCert)
	}

	link := &HTTPLink{
		Cfg:       &cfg,
		Router:    router,
		Svr:       svr,
		MsgRouter: map[string]interface{}{},
	}
	link.initRouter()
	return link, nil
}

func (l *HTTPLink) Start() {
	if l.Svr.TLSConfig == nil {
		if err := l.Svr.ListenAndServe(); err != nil {
			log.L().Info("sync server http stopped", log.Error(err))
		}
	} else {
		if err := l.Svr.ListenAndServeTLS("", ""); err != nil {
			log.L().Info("sync server https stopped", log.Error(err))
		}
	}
}

func (l *HTTPLink) AddMsgRouter(k string, v interface{}) {
	l.MsgRouter[k] = v
}

func (l *HTTPLink) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), l.Cfg.HTTPLink.ShutdownTime)
	return l.Svr.Shutdown(ctx)
}

func (l *HTTPLink) initRouter() {
	l.Router.NoRoute(server.NoRouteHandler)
	l.Router.NoMethod(server.NoMethodHandler)
	l.Router.GET("/health", server.Health)

	l.Router.Use(server.RequestIDHandler)
	l.Router.Use(server.LoggerHandler)
	v1 := l.Router.Group("v1")
	{
		sync := v1.Group("/sync")
		sync.POST("/report", common.Wrapper(l.wrapper(specV1.MessageReport)))
		sync.POST("/desire", common.Wrapper(l.wrapper(specV1.MessageDesire)))
	}
}
