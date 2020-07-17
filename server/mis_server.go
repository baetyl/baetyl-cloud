package server

import (
	"context"
	"github.com/baetyl/baetyl-go/v2/log"
	"net/http"

	"github.com/baetyl/baetyl-cloud/api"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/service"
	"github.com/gin-gonic/gin"
)

// MisServer admin server
type MisServer struct {
	cfg     *config.CloudConfig
	router  *gin.Engine
	server  *http.Server
	api     *api.API
	auth    service.AuthService
}

// NewMisServer create Mis server
func NewMisServer(config *config.CloudConfig) (*MisServer, error) {
	api, err := api.NewAPI(config)
	if err != nil {
		return nil, err
	}
	auth, err := service.NewAuthService(config)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	server := &http.Server{
		Addr:           config.MisServer.Port,
		Handler:        router,
		ReadTimeout:    config.MisServer.ReadTimeout,
		WriteTimeout:   config.MisServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return &MisServer{
		cfg:     config,
		router:  router,
		server:  server,
		auth:    auth,
		api:     api,
	}, nil
}

// Run run server
func (s *MisServer) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.L().Info("admin server stopped", log.Error(err))
	}
}

// Close close server
func (s *MisServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.MisServer.ShutdownTime)
	s.server.Shutdown(ctx)
}
