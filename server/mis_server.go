package server

import (
	"context"
	"github.com/baetyl/baetyl-go/v2/log"
	"net/http"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/gin-gonic/gin"
)

// MisServer mis server
type MisServer struct {
	cfg    *config.CloudConfig
	router *gin.Engine
	server *http.Server
	api    *api.API
}

// NewMisServer create Mis server
func NewMisServer(config *config.CloudConfig) (*MisServer, error) {
	router := gin.New()
	server := &http.Server{
		Addr:           config.MisServer.Port,
		Handler:        router,
		ReadTimeout:    config.MisServer.ReadTimeout,
		WriteTimeout:   config.MisServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return &MisServer{
		cfg:    config,
		router: router,
		server: server,
	}, nil
}

// Run run server
func (s *MisServer) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.L().Info("mis server stopped", log.Error(err))
	}
}

func (s *MisServer) SetAPI(api *api.API) {
	s.api = api
}

// Close close server
func (s *MisServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.MisServer.ShutdownTime)
	s.server.Shutdown(ctx)
}
