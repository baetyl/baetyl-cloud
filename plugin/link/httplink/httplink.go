package httplink

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/server"
)

type httpLink struct {
	cfg       *CloudConfig
	router    *gin.Engine
	svr       *http.Server
	msgRouter map[string]interface{}
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

	link := &httpLink{
		cfg:       &cfg,
		router:    router,
		svr:       svr,
		msgRouter: map[string]interface{}{},
	}
	link.initRouter()
	link.setPortFromEnv()
	return link, nil
}

func (l *httpLink) Start() {
	if l.svr.TLSConfig == nil {
		if err := l.svr.ListenAndServe(); err != nil {
			log.L().Info("sync server http stopped", log.Error(err))
		}
	} else {
		if err := l.svr.ListenAndServeTLS("", ""); err != nil {
			log.L().Info("sync server https stopped", log.Error(err))
		}
	}
}

func (l *httpLink) AddMsgRouter(k string, v interface{}) {
	l.msgRouter[k] = v
}

func (l *httpLink) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), l.cfg.HTTPLink.ShutdownTime)
	return l.svr.Shutdown(ctx)
}

func (l *httpLink) initRouter() {
	l.router.NoRoute(server.NoRouteHandler)
	l.router.NoMethod(server.NoMethodHandler)
	l.router.GET("/health", server.Health)

	l.router.Use(server.RequestIDHandler)
	l.router.Use(server.LoggerHandler)
	v1 := l.router.Group("v1")
	{
		sync := v1.Group("/sync")
		sync.POST("/report", common.Wrapper(l.wrapper(specV1.MessageReport)))
		sync.POST("/desire", common.Wrapper(l.wrapper(specV1.MessageDesire)))
	}
}

func (l *httpLink) setPortFromEnv() {
	nodePort := os.Getenv(config.NodeServerPort)
	if nodePort != "" {
		l.cfg.HTTPLink.Port = ":" + nodePort
	}
}
