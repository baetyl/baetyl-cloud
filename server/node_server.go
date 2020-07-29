package server

import (
	"context"
	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NodeServer struct {
	cfg    *config.CloudConfig
	router *gin.Engine
	server *http.Server
	api    *api.API
}

// NewNodeServer new server
func NewNodeServer(config *config.CloudConfig) (*NodeServer, error) {
	router := gin.New()
	server := &http.Server{
		Addr:           config.NodeServer.Port,
		Handler:        router,
		ReadTimeout:    config.NodeServer.ReadTimeout,
		WriteTimeout:   config.NodeServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	if config.NodeServer.Certificate.Cert != "" &&
		config.NodeServer.Certificate.Key != "" &&
		config.NodeServer.Certificate.CA != "" {
		t, err := utils.NewTLSConfigServer(utils.Certificate{
			CA:   config.NodeServer.Certificate.CA,
			Cert: config.NodeServer.Certificate.Cert,
			Key:  config.NodeServer.Certificate.Key,
		})
		if err != nil {
			return nil, err
		}
		server.TLSConfig = t
	}

	newAPI, err := api.NewAPI(config)
	if err != nil {
		return nil, err
	}

	return &NodeServer{
		cfg:    config,
		router: router,
		server: server,
		api:    newAPI,
	}, nil
}

// Run run server
func (s *NodeServer) Run() {
	if s.server.TLSConfig == nil {
		if err := s.server.ListenAndServe(); err != nil {
			log.L().Info("node server http stopped", log.Error(err))
		}
	} else {
		if err := s.server.ListenAndServeTLS("", ""); err != nil {
			log.L().Info("node server https stopped", log.Error(err))
		}
	}
}

// Close close server
func (s *NodeServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), s.cfg.NodeServer.ShutdownTime)
	s.server.Shutdown(ctx)
}
