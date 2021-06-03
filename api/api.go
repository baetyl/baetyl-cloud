package api

import (
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/facade"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

// API baetyl api server
type API struct {
	NS       service.NamespaceService
	Node     service.NodeService
	Index    service.IndexService
	Func     service.FunctionService
	Obj      service.ObjectService
	PKI      service.PKIService
	Auth     service.AuthService
	Prop     service.PropertyService
	Module   service.ModuleService
	Init     service.InitService
	License  service.LicenseService
	Template service.TemplateService
	Task     service.TaskService
	Locker   service.LockerService
	SysApp   service.SystemAppService
	Sign     service.SignService
	Wrapper  service.WrapperService
	Facade   facade.Facade
	*service.AppCombinedService
	log *log.Logger
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
	signService, err := service.NewSignService(config)
	if err != nil {
		return nil, err
	}
	propertyService, err := service.NewPropertyService(config)
	if err != nil {
		return nil, err
	}
	moduleService, err := service.NewModuleService(config)
	if err != nil {
		return nil, err
	}
	initService, err := service.NewInitService(config)
	if err != nil {
		return nil, err
	}
	licenseService, err := service.NewLicenseService(config)
	if err != nil {
		return nil, err
	}
	templateService, err := service.NewTemplateService(config, map[string]interface{}{
		"GetProperty":      propertyService.GetPropertyValue,
		"RandString":       common.RandString,
		"GetModuleImage":   moduleService.GetLatestModuleImage,
		"GetModuleProgram": moduleService.GetLatestModuleProgram,
	})
	if err != nil {
		return nil, err
	}
	taskService, err := service.NewTaskService(config)
	if err != nil {
		return nil, err
	}
	lockerService, err := service.NewLockerService(config)
	if err != nil {
		return nil, err
	}
	sysApp, err := service.NewSystemAppService(config)
	if err != nil {
		return nil, err
	}
	wrapper, err := service.NewWrapperService(config)
	if err != nil {
		return nil, err
	}
	appFacade, err := facade.NewFacade(config)
	if err != nil {
		return nil, err
	}
	return &API{
		NS:                 namespaceService,
		Node:               nodeService,
		Index:              indexService,
		Obj:                objectService,
		Func:               functionService,
		PKI:                pkiService,
		Auth:               authService,
		Sign:               signService,
		Prop:               propertyService,
		Module:             moduleService,
		Init:               initService,
		License:            licenseService,
		Template:           templateService,
		Task:               taskService,
		Locker:             lockerService,
		SysApp:             sysApp,
		Wrapper:            wrapper,
		AppCombinedService: acs,
		Facade:             appFacade,
		log:                log.L().With(log.Any("api", "admin")),
	}, nil
}
