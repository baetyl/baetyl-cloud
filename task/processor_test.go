package task

import (
	"fmt"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessor(t *testing.T) {
	services := InitMockEnvironment(t)
	defer services.Close()

	err := RegisterNamespaceProcessor(services.conf)
	assert.NoError(t, err)

	np, err := NewNamespaceProcessor(services.conf)
	assert.NoError(t, err)

	task := gentTask()
	services.namespace.EXPECT().DeleteNamespace(&models.Namespace{Name: task.ResourceName}).Return(nil)
	err = np.DeleteNamespace(task)
	assert.NoError(t, err)

	services.namespace.EXPECT().DeleteNamespace(&models.Namespace{Name: task.ResourceName}).Return(fmt.Errorf("namespace_delete_error"))
	err = np.DeleteNamespace(task)
	assert.Error(t, err)

	services.license.EXPECT().DeleteQuotaByNamespace(task.ResourceName).Return(nil)
	err = np.DeleteQuotaByNamespace(task)
	assert.NoError(t, err)

	services.license.EXPECT().DeleteQuotaByNamespace(task.ResourceName).Return(fmt.Errorf("quota_delete_error"))
	err = np.DeleteQuotaByNamespace(task)
	assert.Error(t, err)

	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Config, common.Application).Return(nil, nil)
	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Secret, common.Application).Return(nil, nil)
	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Node, common.Application).Return(nil, nil)
	err = np.DeleteIndexByNamespace(task)
	assert.NoError(t, err)

	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Config, common.Application).Return(nil, fmt.Errorf("err_delete_index_app"))
	err = np.DeleteIndexByNamespace(task)
	assert.Error(t, err)

	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Config, common.Application).Return(nil, nil)
	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Node, common.Application).Return(nil, fmt.Errorf("err_delete_index_node"))
	err = np.DeleteIndexByNamespace(task)
	assert.Error(t, err)

	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Config, common.Application).Return(nil, nil)
	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Node, common.Application).Return(nil, nil)
	services.index.EXPECT().DeleteIndexByNamespace(task.Namespace, common.Secret, common.Application).Return(nil, fmt.Errorf("err_delete_index_secret"))
	err = np.DeleteIndexByNamespace(task)
	assert.Error(t, err)
}

func gentTask() *models.Task {
	return &models.Task{
		Name:             "task_processor01",
		Namespace:        "default",
		RegistrationName: "namespace_delete",
		ResourceType:     "namespace",
		ResourceName:     "default",
	}
}
