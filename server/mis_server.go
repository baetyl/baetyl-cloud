package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/gin-gonic/gin"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
)

// MisServer mis server
type MisServer struct {
	cfg    *config.CloudConfig
	router *gin.Engine
	server *http.Server
	api    *api.API // TODO: define independent api
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

func (s *MisServer) GetRoute() *gin.Engine {
	return s.router
}

func (s *MisServer) InitRoute() {
	s.router.NoRoute(NoRouteHandler)
	s.router.NoMethod(NoMethodHandler)
	s.router.GET("/health", Health)

	s.router.Use(RequestIDHandler)
	s.router.Use(LoggerHandler)
	s.router.Use(s.authHandler)
	v1 := s.router.Group("v1")
	{
		cache := v1.Group("/properties")

		cache.POST("", common.WrapperMis(s.api.CreateProperty))
		cache.DELETE("/:name", common.WrapperMis(s.api.DeleteProperty))
		cache.GET("", common.WrapperMis(s.api.ListProperty))
		cache.PUT("/:name", common.WrapperMis(s.api.UpdateProperty))
	}
	{
		quota := v1.Group("/quotas")

		quota.POST("", common.WrapperMis(s.api.CreateQuota))
		quota.DELETE("", common.WrapperMis(s.api.DeleteQuota))
		quota.GET("/:namespace", common.WrapperMis(s.api.GetQuotaForMis))
		quota.PUT("", common.WrapperMis(s.api.UpdateQuota))
	}
	{
		module := v1.Group("/modules")

		module.GET("", common.WrapperMis(s.api.ListModules))
		module.GET("/:name", common.WrapperMis(s.api.GetModules))
		module.GET("/:name/version/:version", common.WrapperMis(s.api.GetModuleByVersion))
		module.GET("/:name/latest", common.WrapperMis(s.api.GetLatestModule))
		module.POST("", common.WrapperMis(s.api.CreateModule))
		module.PUT("/:name/version/:version", common.WrapperMis(s.api.UpdateModule))
		module.DELETE("/:name", common.WrapperMis(s.api.DeleteModules))
		module.DELETE("/:name/version/:version", common.WrapperMis(s.api.DeleteModules))
	}
}

// auth handler
func (s *MisServer) authHandler(c *gin.Context) {
	cc := common.NewContext(c)

	token := c.Request.Header.Get(s.cfg.MisServer.TokenHeader)
	if strings.Compare(token, s.cfg.MisServer.AuthToken) == 0 {
		user := c.Request.Header.Get(s.cfg.MisServer.UserHeader)
		if len(user) != 0 {
			log.L().Info("mis server accessed",
				log.Any("user", user),
				log.Any(cc.GetTrace()),
			)
			return
		}
	}
	err := common.Error(common.ErrRequestAccessDenied, common.Field("error", common.Code(common.ErrRequestAccessDenied)))
	log.L().Error(common.Code(common.ErrRequestAccessDenied).String(),
		log.Any(cc.GetTrace()),
		log.Code(err),
		log.Error(err),
	)
	common.PopulateFailedResponse(cc, err, true)
}
