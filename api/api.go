package api

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

// API baetyl api server
type API struct {
	namespaceService   service.NamespaceService
	applicationService service.ApplicationService
	nodeService        service.NodeService
	configService      service.ConfigService
	secretService      service.SecretService
	indexService       service.IndexService
	functionService    service.FunctionService
	objectService      service.ObjectService
	sysConfigService   service.SysConfigService
	pkiService         service.PKIService
	authService        service.AuthService
	propertyService    service.PropertyService
	templateService    service.TemplateService
	activeService      service.ActiveService
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
	pkiService, err := service.NewPKIService(config)
	if err != nil {
		return nil, err
	}
	authService, err := service.NewAuthService(config)
	if err != nil {
		return nil, err
	}
	propertyService, err := service.NewPropertyService(config)
	if err != nil {
		return nil, err
	}
	templateService, err := service.NewTemplateService(config)
	if err != nil {
		return nil, err
	}
	activeService, err := service.NewActiveService(config)
	if err != nil {
		return nil, err
	}
	return &API{
		applicationService: applicationService,
		nodeService:        nodeService,
		configService:      configService,
		namespaceService:   namespaceService,
		secretService:      secretService,
		indexService:       indexService,
		functionService:    functionService,
		objectService:      objectService,
		sysConfigService:   sysConfigService,
		pkiService:         pkiService,
		authService:        authService,
		propertyService:    propertyService,
		templateService:    templateService,
		activeService:      activeService,
	}, nil
}
