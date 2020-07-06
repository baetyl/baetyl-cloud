package server

import (
	"context"
	"github.com/baetyl/baetyl-go/log"
	"net/http"

	"github.com/baetyl/baetyl-cloud/api"
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/service"
	"github.com/gin-gonic/gin"
)

// AmisServer admin server
type AmisServer struct {
	cfg     *config.CloudConfig
	router  *gin.Engine
	server  *http.Server
	api     *api.API
	auth    service.AuthService
	license service.LicenseService
}


// NewAmisServer create amis server
func NewAmisServer(config *config.CloudConfig) (*AmisServer, error) {
	api, err := api.NewAPI(config)
	if err != nil {
		return nil, err
	}
	auth, err := service.NewAuthService(config)
	if err != nil {
		return nil, err
	}

	ls, err := service.NewLicenseService(config)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	server := &http.Server{
		Addr:           config.AmisServer.Port,
		Handler:        router,
		ReadTimeout:    config.AmisServer.ReadTimeout,
		WriteTimeout:   config.AmisServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return &AmisServer{
		cfg:     config,
		router:  router,
		server:  server,
		auth:    auth,
		api:     api,
		license: ls,
	}, nil
}

// Run run server
func (s *AmisServer) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.L().Info("admin server stopped", log.Error(err))
	}
}

// Close close server
func (s *AmisServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.AmisServer.ShutdownTime)
	s.server.Shutdown(ctx)
}
