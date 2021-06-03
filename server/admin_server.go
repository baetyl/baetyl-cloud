package server

import (
	"context"
	"net/http"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/gin-gonic/gin"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

// AdminServer admin server
type AdminServer struct {
	Auth             service.AuthService
	License          service.LicenseService
	ExternalHandlers []gin.HandlerFunc

	cfg    *config.CloudConfig
	router *gin.Engine
	server *http.Server
	api    *api.API
	log    *log.Logger
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
		Auth:    auth,
		License: ls,
		log:     log.L().With(log.Any("server", "AdminServer")),
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

// InitRoute init router
func (s *AdminServer) InitRoute() {
	s.router.NoRoute(NoRouteHandler)
	s.router.NoMethod(NoMethodHandler)
	s.router.GET("/health", Health)

	s.router.Use(RequestIDHandler)
	s.router.Use(LoggerHandler)
	s.router.Use(s.AuthHandler)
	s.router.Use(s.ExternalHandlers...)
	v1 := s.router.Group("v1")
	{
		configs := v1.Group("/configs")
		configs.GET("/:name", common.Wrapper(s.api.GetConfig))
		configs.PUT("/:name", common.WrapperWithLock(s.api.Locker.Lock, s.api.Locker.Unlock), common.Wrapper(s.api.UpdateConfig))
		configs.DELETE("/:name", common.WrapperRaw(s.api.ValidateResourceForDeleting, true), common.Wrapper(s.api.DeleteConfig))
		configs.POST("", common.WrapperRaw(s.api.ValidateResourceForCreating, true), common.Wrapper(s.api.CreateConfig))
		configs.GET("", common.Wrapper(s.api.ListConfig))
		configs.GET("/:name/apps", common.Wrapper(s.api.GetAppByConfig))
	}
	{
		registry := v1.Group("/registries")
		registry.GET("/:name", common.Wrapper(s.api.GetRegistry))
		registry.PUT("/:name", common.Wrapper(s.api.UpdateRegistry))
		registry.POST("/:name/refresh", common.Wrapper(s.api.RefreshRegistryPassword))
		registry.DELETE("/:name", common.WrapperRaw(s.api.ValidateResourceForDeleting, true), common.Wrapper(s.api.DeleteRegistry))
		registry.POST("", common.WrapperRaw(s.api.ValidateResourceForCreating, true), common.Wrapper(s.api.CreateRegistry))
		registry.GET("", common.Wrapper(s.api.ListRegistry))
		registry.GET("/:name/apps", common.Wrapper(s.api.GetAppByRegistry))
	}
	{
		certificate := v1.Group("/certificates")
		certificate.GET("/:name", common.Wrapper(s.api.GetCertificate))
		certificate.PUT("/:name", common.WrapperWithLock(s.api.Locker.Lock, s.api.Locker.Unlock), common.Wrapper(s.api.UpdateCertificate))
		certificate.DELETE("/:name", common.WrapperRaw(s.api.ValidateResourceForDeleting, true), common.Wrapper(s.api.DeleteCertificate))
		certificate.POST("", common.WrapperRaw(s.api.ValidateResourceForCreating, true), common.Wrapper(s.api.CreateCertificate))
		certificate.GET("", common.Wrapper(s.api.ListCertificate))
		certificate.GET("/:name/apps", common.Wrapper(s.api.GetAppByCertificate))
	}
	{
		configs := v1.Group("/secrets")
		configs.GET("/:name", common.Wrapper(s.api.GetSecret))
		configs.PUT("/:name", common.Wrapper(s.api.UpdateSecret))
		configs.DELETE("/:name", common.WrapperRaw(s.api.ValidateResourceForDeleting, true), common.Wrapper(s.api.DeleteSecret))
		configs.POST("", common.WrapperRaw(s.api.ValidateResourceForCreating, true), common.Wrapper(s.api.CreateSecret))
		configs.GET("", common.Wrapper(s.api.ListSecret))
		configs.GET("/:name/apps", common.Wrapper(s.api.GetAppBySecret))
	}
	{
		nodes := v1.Group("/nodes")
		nodes.GET("/:name", common.Wrapper(s.api.GetNode))
		nodes.PUT("", common.Wrapper(s.api.GetNodes))
		nodes.GET("/:name/apps", common.Wrapper(s.api.GetAppByNode))
		nodes.GET("/:name/stats", common.Wrapper(s.api.GetNodeStats))
		nodes.PUT("/:name", common.WrapperWithLock(s.api.Locker.Lock, s.api.Locker.Unlock), common.Wrapper(s.api.UpdateNode))
		nodes.DELETE("/:name", common.Wrapper(s.api.DeleteNode))
		nodes.POST("", s.NodeQuotaHandler, common.Wrapper(s.api.CreateNode))
		nodes.GET("", common.Wrapper(s.api.ListNode))
		nodes.GET("/:name/deploys", common.Wrapper(s.api.GetNodeDeployHistory))
		nodes.GET("/:name/init", common.Wrapper(s.api.GenInitCmdFromNode))
		nodes.PUT("/:name/mode", common.Wrapper(s.api.UpdateNodeMode))
		nodes.PUT("/:name/properties", common.Wrapper(s.api.UpdateNodeProperties))
		nodes.GET("/:name/properties", common.Wrapper(s.api.GetNodeProperties))
		nodes.PUT("/:name/core/configs", common.Wrapper(s.api.UpdateCoreApp))
		nodes.GET("/:name/core/configs", common.Wrapper(s.api.GetCoreAppConfigs))
		nodes.GET("/:name/core/versions", common.Wrapper(s.api.GetCoreAppVersions))
	}
	{
		apps := v1.Group("/apps")
		apps.GET("/:name", common.Wrapper(s.api.GetApplication))
		apps.GET("/:name/configs", common.Wrapper(s.api.GetSysAppConfigs))
		apps.GET("/:name/secrets", common.Wrapper(s.api.GetSysAppSecrets))
		apps.GET("/:name/certificates", common.Wrapper(s.api.GetSysAppCertificates))
		apps.GET("/:name/registries", common.Wrapper(s.api.GetSysAppRegistries))
		apps.PUT("/:name", common.WrapperWithLock(s.api.Locker.Lock, s.api.Locker.Unlock), common.Wrapper(s.api.UpdateApplication))
		apps.DELETE("/:name", common.WrapperRaw(s.api.ValidateResourceForDeleting, true), common.Wrapper(s.api.DeleteApplication))
		apps.POST("", common.WrapperRaw(s.api.ValidateResourceForCreating, true), common.WrapperWithLock(s.api.Locker.Lock, s.api.Locker.Unlock), common.Wrapper(s.api.CreateApplication))
		apps.GET("", common.Wrapper(s.api.ListApplication))
	}
	{
		namespace := v1.Group("/namespace")
		namespace.POST("", common.Wrapper(s.api.CreateNamespace))
		namespace.GET("", common.Wrapper(s.api.GetNamespace))
		namespace.DELETE("", common.Wrapper(s.api.DeleteNamespace))
	}
	{
		function := v1.Group("/functions")
		function.GET("", common.Wrapper(s.api.ListFunctionSources))
		if len(s.cfg.Plugin.Functions) != 0 {
			function.GET("/:source/functions", common.Wrapper(s.api.ListFunctions))
			function.GET("/:source/functions/:name/versions", common.Wrapper(s.api.ListFunctionVersions))
			function.POST("/:source/functions/:name/versions/:version", common.Wrapper(s.api.ImportFunction))
		}
	}
	{
		// Deprecated
		objects := v1.Group("/objects")
		objects.GET("", common.Wrapper(s.api.ListObjectSources))
		if len(s.cfg.Plugin.Objects) != 0 {
			objects.GET("/:source/buckets", common.Wrapper(s.api.ListBuckets))
			objects.GET("/:source/buckets/:bucket/objects", common.Wrapper(s.api.ListBucketObjects))
		}
	}

	{
		properties := v1.Group("properties")
		properties.GET("/:name", common.Wrapper(s.api.GetProperty))

		// TODO: deprecated, to use property api
		sysconfig := v1.Group("sysconfig")
		sysconfig.GET("/baetyl_version/latest", common.Wrapper(func(c *common.Context) (interface{}, error) {
			res, err := s.api.Module.GetLatestModule("baetyl")
			if err != nil {
				return nil, err
			}
			return map[string]string{
				"type":  "baetyl_version",
				"key":   "latest",
				"value": res.Version,
			}, nil
		}))
		sysconfig.GET("/baetyl-function-runtime", common.Wrapper(func(c *common.Context) (interface{}, error) {
			runtimes, err := s.api.Func.ListRuntimes()
			if err != nil {
				return nil, errors.Trace(err)
			}
			var runtimesView []map[string]string
			for k, v := range runtimes {
				runtimesView = append(runtimesView, map[string]string{
					"type":  "baetyl-function-runtime",
					"key":   k,
					"value": v,
				})
			}
			// {"sysconfigs":[{"type":"baetyl-function-runtime","key":"nodejs10","value":"hub.baidubce.com/baetyl/function-node:10.19-v2.0.0","createTime":"2020-08-20T05:16:27Z","updateTime":"2020-08-20T05:16:27Z"},{"type":"baetyl-function-runtime","key":"python3","value":"hub.baidubce.com/baetyl/function-python:3.6-v2.0.0","createTime":"2020-08-20T05:16:27Z","updateTime":"2020-08-20T05:16:27Z"},{"type":"baetyl-function-runtime","key":"python3-opencv","value":"hub.baidubce.com/baetyl/function-python-opencv:3.6","createTime":"2020-04-26T06:39:32Z","updateTime":"2020-04-26T06:39:32Z"},{"type":"baetyl-function-runtime","key":"sql","value":"hub.baidubce.com/baetyl-sandbox/function-sql:git-4a62dfc","createTime":"2020-08-20T05:16:27Z","updateTime":"2020-08-25T03:16:39Z"}]}
			return map[string]interface{}{
				"sysconfigs": runtimesView,
			}, nil
		}))
	}
	{
		module := v1.Group("modules")
		module.GET("", common.Wrapper(s.api.ListModules))
		module.GET("/:name", common.Wrapper(s.api.GetModules))
		module.GET("/:name/version/:version", common.Wrapper(s.api.GetModuleByVersion))
		module.GET("/:name/latest", common.Wrapper(s.api.GetLatestModule))
		module.POST("", common.Wrapper(s.api.CreateModule))
		module.PUT("/:name/version/:version", common.Wrapper(s.api.UpdateModule))
		module.DELETE("/:name", common.Wrapper(s.api.DeleteModules))
		module.DELETE("/:name/version/:version", common.Wrapper(s.api.DeleteModules))
	}
	{
		quotas := v1.Group("/quotas")
		quotas.GET("", common.Wrapper(s.api.GetQuota))
	}

	v2 := s.router.Group("v2")
	{
		objects := v2.Group("/objects")
		objects.GET("", common.Wrapper(s.api.ListObjectSourcesV2))
		if len(s.cfg.Plugin.Objects) != 0 {
			objects.GET("/:source/buckets", common.Wrapper(s.api.ListBucketsV2))
			objects.GET("/:source/buckets/:bucket/objects", common.Wrapper(s.api.ListBucketObjectsV2))
		}
	}
}

// GetRoute get router
func (s *AdminServer) GetRoute() *gin.Engine {
	return s.router
}

// auth handler
func (s *AdminServer) AuthHandler(c *gin.Context) {
	cc := common.NewContext(c)
	err := s.Auth.Authenticate(cc)
	if err != nil {
		s.log.Error("request authenticate failed",
			log.Any(cc.GetTrace()),
			log.Any("namespace", cc.GetNamespace()),
			log.Any("authorization", c.Request.Header.Get("Authorization")),
			log.Error(err))
		common.PopulateFailedResponse(cc, common.Error(common.ErrRequestAccessDenied), true)
	}
}

func (s *AdminServer) NodeQuotaHandler(c *gin.Context) {
	cc := common.NewContext(c)
	namespace := cc.GetNamespace()
	if err := s.api.License.CheckQuota(namespace, s.api.NodeNumberCollector); err != nil {
		s.log.Error("quota out of limit",
			log.Any(cc.GetTrace()),
			log.Any("namespace", cc.GetNamespace()),
			log.Error(err))
		common.PopulateFailedResponse(cc, err, true)
	}
}
