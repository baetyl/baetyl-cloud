// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"
)

const (
	// 北向协议目前支持协议
	UplinkProtocolMQTT      = "mqtt"
	UplinkProtocolHTTP      = "http"
	UplinkProtocolWebsocket = "websocket"

	DestinationBaetylBroker = "baetyl-broker"
	DestinationBaetylDmp    = "baetyl-dmp"
	DestinationCustom       = "custom"
)

type DeviceUplink struct {
	ID              int64     `db:"id"`
	NodeName        string    `db:"node_name"`
	Namespace       string    `db:"namespace"`
	Protocol        string    `db:"protocol"`
	Destination     string    `db:"destination"`
	DestinationName string    `db:"destination_name"`
	Address         string    `db:"address"`
	MQTTUser        string    `db:"mqtt_user"`
	MQTTPassword    string    `db:"mqtt_password"`
	HTTPMethod      string    `db:"http_method"`
	HTTPPath        string    `db:"http_path"`
	CA              string    `db:"ca"`
	Cert            string    `db:"cert"`
	PrivateKey      string    `db:"private_key"`
	Passphrase      string    `db:"passphrase"`
	CreateTime      time.Time `db:"create_time"`
	UpdateTime      time.Time `db:"update_time"`
}
