package server

import (
	"context"
	"net/http"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/gin-gonic/gin"
)

// AdminServer admin server
type AdminServer struct {
	cfg     *config.CloudConfig
	router  *gin.Engine
	server  *http.Server
	api     *api.API
	auth    service.AuthService
	license service.LicenseService
}

// NewAdminServer create admin server
func NewAdminServer(config *config.CloudConfig) (*AdminServer, error) {
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
		Addr:           config.AdminServer.Port,
		Handler:        router,
		ReadTimeout:    config.AdminServer.ReadTimeout,
		WriteTimeout:   config.AdminServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return &AdminServer{
		cfg:     config,
		router:  router,
		server:  server,
		auth:    auth,
		license: ls,
	}, nil
}

func (s *AdminServer) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.L().Info("admin server stopped", log.Error(err))
	}
}

func (s *AdminServer) SetAPI(api *api.API) {
	s.api = api
}

// Close close server
func (s *AdminServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.AdminServer.ShutdownTime)
	s.server.Shutdown(ctx)
}
