package service

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/errors"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/config"
	mockPlugin "github.com/baetyl/baetyl-cloud/v2/mock/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func mockTx(mock plugin.TransactionFactory) plugin.Factory {
	factory := func() (plugin.Plugin, error) {
		return mock, nil
	}
	return factory
}

func mockCreateNodeY(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error) {
	return nil, nil
}

func mockCreateNodeN(tx interface{}, namespace string, node *specV1.Node) (*specV1.Node, error) {
	return nil, errors.New("error")
}

func TestCreateWrapper(t *testing.T) {
	cfg := &config.CloudConfig{}
	cfg.Plugin.Tx = common.RandString(9)
	mockCtl := gomock.NewController(t)

	mTx := mockPlugin.NewMockTransactionFactory(mockCtl)
	_, err := NewWrapperService(cfg)
	assert.Error(t, err)

	plugin.RegisterFactory(cfg.Plugin.Tx, mockTx(mTx))
	wrapper, err := NewWrapperService(cfg)
	assert.NoError(t, err)

	mTx.EXPECT().BeginTx().Return(nil, errors.New("error"))
	_, err = wrapper.CreateNodeTx(mockCreateNodeY)(nil, "", nil)
	assert.Error(t, err)

	mTx.EXPECT().BeginTx().Return(nil, nil)
	mTx.EXPECT().Rollback(nil).Return()
	_, err = wrapper.CreateNodeTx(mockCreateNodeN)(nil, "", nil)
	assert.Error(t, err)

	mTx.EXPECT().BeginTx().Return(nil, nil)
	mTx.EXPECT().Commit(nil).Return()
	_, err = wrapper.CreateNodeTx(mockCreateNodeY)(nil, "", nil)
	assert.NoError(t, err)
}
