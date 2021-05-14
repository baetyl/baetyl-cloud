package facade

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/golang/mock/gomock"

	mp "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	ms "github.com/baetyl/baetyl-cloud/v2/mock/service"
)

var (
	unknownErr = errors.New("unknown")
)

type MockAppFacade struct {
	sNode     *ms.MockNodeService
	sApp      *ms.MockApplicationService
	sConfig   *ms.MockConfigService
	sIndex    *ms.MockIndexService
	txFactory *mp.MockTransactionFactory
}

func InitMockEnvironment(t *testing.T) (*MockAppFacade, *gomock.Controller) {
	mockCtl := gomock.NewController(t)
	return &MockAppFacade{
		sNode:     ms.NewMockNodeService(mockCtl),
		sApp:      ms.NewMockApplicationService(mockCtl),
		sConfig:   ms.NewMockConfigService(mockCtl),
		sIndex:    ms.NewMockIndexService(mockCtl),
		txFactory: mp.NewMockTransactionFactory(mockCtl),
	}, mockCtl
}
