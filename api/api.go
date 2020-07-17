package api

import (
	"github.com/baetyl/baetyl-cloud/config"
	"github.com/baetyl/baetyl-cloud/service"
)

// API baetyl api server
type API struct {
	namespaceService   service.NamespaceService
	applicationService service.ApplicationService
	nodeService        service.NodeService
	configService      service.ConfigService
	syncService        service.SyncService
	registerService    service.RegisterService
	secretService      service.SecretService
	indexService       service.IndexService
	functionService    service.FunctionService
	objectService      service.ObjectService
	sysConfigService   service.SysConfigService
	callbackService    service.CallbackService
	pkiService         service.PKIService
	initService        service.InitializeService
	authService        service.AuthService
	cacheService 	   service.CacheService
}

// NewAPI NewAPI
func NewAPI(config *config.CloudConfig) (*API, error) {
	applicationService, err := service.NewApplicationService(config)
	if err != nil {
		return nil, err
	}
	nodeService, err := service.NewNodeService(config)
	if err != nil {
		return nil, err
	}
	configService, err := service.NewConfigService(config)
	if err != nil {
		return nil, err
	}
	syncService, err := service.NewSyncService(config)
	if err != nil {
		return nil, err
	}
	registerService, err := service.NewRegisterService(config)
	if err != nil {
		return nil, err
	}
	namespaceService, err := service.NewNamespaceService(config)
	if err != nil {
		return nil, err
	}
	secretService, err := service.NewSecretService(config)
	if err != nil {
		return nil, err
	}
	indexService, err := service.NewIndexService(config)
	if err != nil {
		return nil, err
	}
	functionService, err := service.NewFunctionService(config)
	if err != nil {
		return nil, err
	}
	objectService, err := service.NewObjectService(config)
	if err != nil {
		return nil, err
	}
	sysConfigService, err := service.NewSysConfigService(config)
	if err != nil {
		return nil, err
	}
	callbackService, err := service.NewCallbackService(config)
	if err != nil {
		return nil, err
	}
	pkiService, err := service.NewPKIService(config)
	if err != nil {
		return nil, err
	}
	initService, err := service.NewInitializeService(config)
	if err != nil {
		return nil, err
	}
	authService, err := service.NewAuthService(config)
	if err != nil {
		return nil, err
	}
	cacheService, err := service.NewCacheService(config)
	if err != nil {
		return nil, err
	}

	return &API{
		applicationService: applicationService,
		nodeService:        nodeService,
		configService:      configService,
		syncService:        syncService,
		registerService:    registerService,
		namespaceService:   namespaceService,
		secretService:      secretService,
		indexService:       indexService,
		functionService:    functionService,
		objectService:      objectService,
		sysConfigService:   sysConfigService,
		callbackService:    callbackService,
		pkiService:         pkiService,
		initService:        initService,
		authService:        authService,
		cacheService: 		cacheService,
	}, nil
}
