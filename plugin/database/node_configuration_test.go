package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	nodeConfiguration = `

CREATE TABLE baetyl_node_configuration
(
    id                integer       PRIMARY KEY AUTOINCREMENT,
    node_name        varchar(250)  default ''                NOT NULL ,
    namespace        varchar(256)  default ''                NOT NULL ,
    data			 varchar(256)  default '{"enable":0}'    NOT NULL ,
    type			 varchar(10)   default ''                NOT NULL ,
    create_time      timestamp     default CURRENT_TIMESTAMP not null ,
    update_time      timestamp     default CURRENT_TIMESTAMP not null
) ;`
)

func (d *BaetylCloudDB) MockCreateNodeConfigurationTable() {
	_, err := d.Exec(nil, nodeConfiguration)
	if err != nil {
		panic(fmt.Sprintf("create  baetyl node configuration exception: %s", err.Error()))
	}
}

func genNodeConfiguration() *models.NodeConfiguration {
	return &models.NodeConfiguration{
		ConfigurationType: "cache",
		NodeName:          "test",
		Data:              "{\"enable\":false}",
		Namespace:         "baetyl_cloud",
	}
}

func TestGetNodeConfig(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateNodeConfigurationTable()
	nc, err := db.CreateNodeConfig(nil, genNodeConfiguration())
	assert.NoError(t, err)
	assert.NotNil(t, nc)
	nc, err = db.GetNodeConfig(nil, "baetyl_cloud", "test", "cache")
	assert.NoError(t, err)
	assert.NotNil(t, nc)
	nc, err = db.UpdateNodeConfig(nil, genNodeConfiguration())
	assert.NoError(t, err)
	assert.NotNil(t, nc)
}
