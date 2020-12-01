package server

import (
	"context"
	"net/http"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
)

type InitServer struct {
	cfg    *config.CloudConfig
	router *gin.Engine
	server *http.Server
	api    *api.InitAPI
}

// NewInitServer new init server
func NewInitServer(config *config.CloudConfig) (*InitServer, error) {
	router := gin.New()
	server := &http.Server{
		Addr:           config.InitServer.Port,
		Handler:        router,
		ReadTimeout:    config.InitServer.ReadTimeout,
		WriteTimeout:   config.InitServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	if config.InitServer.Certificate.Cert != "" &&
		config.InitServer.Certificate.Key != "" {
		t, err := utils.NewTLSConfigServer(utils.Certificate{
			Cert: config.InitServer.Certificate.Cert,
			Key:  config.InitServer.Certificate.Key,
		})
		if err != nil {
			return nil, err
		}
		server.TLSConfig = t
	}

	return &InitServer{
		cfg:    config,
		router: router,
		server: server,
	}, nil
}

// Run run server
func (s *InitServer) Run() {
	if s.server.TLSConfig == nil {
		if err := s.server.ListenAndServe(); err != nil {
			log.L().Info("init server http stopped", log.Error(err))
		}
	} else {
		if err := s.server.ListenAndServeTLS("", ""); err != nil {
			log.L().Info("init server https stopped", log.Error(err))
		}
	}
}

func (s *InitServer) SetAPI(api *api.InitAPI) {
	s.api = api
}

// Close close server
func (s *InitServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.InitServer.ShutdownTime)
	s.server.Shutdown(ctx)
}

// GetRoute get router
func (s *InitServer) GetRoute() *gin.Engine {
	return s.router
}

func (s *InitServer) InitRoute() {
	s.router.NoRoute(NoRouteHandler)
	s.router.NoMethod(NoMethodHandler)
	s.router.GET("/health", Health)

	s.router.Use(RequestIDHandler)
	s.router.Use(LoggerHandler)
	v1 := s.router.Group("v1")
	{
		// TODO: deprecated
		active := v1.Group("/active")
		active.GET("/:resource", common.WrapperRaw(s.api.GetResource, true))
	}
	{
		initz := v1.Group("/init")
		initz.GET("/:resource", common.WrapperRaw(s.api.GetResource, true))
	}
}
