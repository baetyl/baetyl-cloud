package server

import (
	"context"
	"net/http"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin"
)

type ActiveServer struct {
	cfg    *config.CloudConfig
	router *gin.Engine
	server *http.Server
	api    *api.API
}

// NewActiveServer new active server
func NewActiveServer(config *config.CloudConfig) (*ActiveServer, error) {
	router := gin.New()
	server := &http.Server{
		Addr:           config.ActiveServer.Port,
		Handler:        router,
		ReadTimeout:    config.ActiveServer.ReadTimeout,
		WriteTimeout:   config.ActiveServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	if config.ActiveServer.Certificate.Cert != "" &&
		config.ActiveServer.Certificate.Key != "" {
		t, err := utils.NewTLSConfigServer(utils.Certificate{
			Cert: config.ActiveServer.Certificate.Cert,
			Key:  config.ActiveServer.Certificate.Key,
		})
		if err != nil {
			return nil, err
		}
		server.TLSConfig = t
	}

	api, err := api.NewAPI(config)
	if err != nil {
		return nil, err
	}

	return &ActiveServer{
		cfg:    config,
		router: router,
		server: server,
		api:    api,
	}, nil
}

// Run run server
func (s *ActiveServer) Run() {
	if s.server.TLSConfig == nil {
		if err := s.server.ListenAndServe(); err != nil {
			log.L().Info("active server http stopped", log.Error(err))
		}
	} else {
		if err := s.server.ListenAndServeTLS("", ""); err != nil {
			log.L().Info("active server https stopped", log.Error(err))
		}
	}
}

// Close close server
func (s *ActiveServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.ActiveServer.ShutdownTime)
	s.server.Shutdown(ctx)
}
