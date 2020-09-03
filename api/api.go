package api

import (
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

// API baetyl api server
type API struct {
	namespaceService service.NamespaceService
	nodeService      service.NodeService
	indexService     service.IndexService
	functionService  service.FunctionService
	objectService    service.ObjectService
	pkiService       service.PKIService
	authService      service.AuthService
	propertyService  service.PropertyService
	initService      service.InitService
	*service.AppCombinedService
}

// NewAPI NewAPI
func NewAPI(config *config.CloudConfig) (*API, error) {
	acs, err := service.NewAppCombinedService(config)
	if err != nil {
		return nil, err
	}
	nodeService, err := service.NewNodeService(config)
	if err != nil {
		return nil, err
	}
	namespaceService, err := service.NewNamespaceService(config)
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
	initService, err := service.NewInitService(config)
	if err != nil {
		return nil, err
	}
	return &API{
		nodeService:        nodeService,
		namespaceService:   namespaceService,
		indexService:       indexService,
		functionService:    functionService,
		objectService:      objectService,
		pkiService:         pkiService,
		authService:        authService,
		propertyService:    propertyService,
		initService:        initService,
		AppCombinedService: acs,
	}, nil
}