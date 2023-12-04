package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/database/entities"
)

var (
	deviceUplinkTables = []string{
		`
CREATE TABLE baetyl_device_uplink
(
	id                integer       PRIMARY KEY AUTOINCREMENT,
    node_name        varchar(250)  default ''                not null ,
    namespace        varchar(256)  default ''                not null ,
    protocol         varchar(50)   default ''                not null ,
    destination      varchar(50)   default ''                not null,
    destination_name varchar(255)  default ''                not null ,
    address          varchar(2048) default ''                not null ,
    mqtt_user        varchar(255)  default ''                not null ,
    mqtt_password    varchar(255)                            not null ,
    http_method      varchar(50)   default ''                not null ,
    http_path        varchar(2048) default ''                not null ,
    ca               text                                    null ,
    cert             text                                    null ,
    private_key      text                                    null ,
    passphrase       varchar(128)  default ''                not null ,
    create_time      timestamp     default CURRENT_TIMESTAMP not null ,
    update_time      timestamp     default CURRENT_TIMESTAMP not null 
);
`,
	}
)

func (d *BaetylCloudDB) MockCreateDeviceUplinkTable() {
	for _, sql := range deviceUplinkTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create device model exception: %s", err.Error()))
		}
	}
}

func TestDeviceUplink(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateDeviceUplinkTable()

	mode, err := db.CreateDeviceUplink(genDeviceUplinkMode())
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", mode.Address)

	mode, err = db.GetDeviceUplink("baetyl", "node", "http")
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", mode.Address)

	data := genDeviceUplinkMode()
	data.Cert = "ccccccc"
	mode, err = db.UpdateDeviceUplink(data)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", mode.Address)
	assert.Equal(t, "ccccccc", mode.Cert)

	list, err := db.ListDeviceUplink("baetyl", "node", &models.ListOptions{Filter: models.Filter{
		PageNo:   1,
		PageSize: 10,
		Name:     "",
	}})
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", list[0].Address)

	err = db.DeleteDeviceUplink("baetyl", "node", "http")
	assert.NoError(t, err)

	count, err := db.GetDeviceUplinkCount("baetyl", "node", &models.ListOptions{})
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	_, err = db.GetDeviceUplink("baetyl", "node", "http")
	assert.Error(t, err)

}

func genDeviceUplinkMode() *entities.DeviceUplink {
	return &entities.DeviceUplink{
		ID:              1,
		NodeName:        "node",
		Namespace:       "baetyl",
		Protocol:        "http",
		Destination:     "custom",
		DestinationName: "http",
		Address:         "127.0.0.1",
		MQTTUser:        "",
		MQTTPassword:    "",
		HTTPMethod:      "",
		HTTPPath:        "",
		CA:              "",
		Cert:            "",
		PrivateKey:      "",
		Passphrase:      "",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
}
