// Package models 模型定义
package models

import (
	"time"
)

type DeviceUplinkListView struct {
	Items        []DeviceUplinkView `json:"items"`
	*ListOptions `json:",inline"`
	Total        int         `json:"total"`
	TopicInfo    []TopicInfo `json:"topicInfo"`
}

type TopicInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateUplink struct {
	NodeName        string `json:"nodeName" form:"nodeName" `
	Namespace       string `json:"namespace" form:"namespace"`
	Protocol        string `json:"protocol" form:"protocol" `
	Destination     string `json:"destination" form:"destination" binding:"required"`
	DestinationName string `json:"destinationName" form:"destinationName" binding:"required,res_name"`
	Address         string `json:"address" form:"address" `
	MQTTUser        string `json:"mqttUser" form:"mqttUser"`
	MQTTPassword    string `json:"mqttPassword" form:"mqttPassword"`
	HTTPMethod      string `json:"httpMethod" form:"httpMethod"`
	HTTPPath        string `json:"httpPath" form:"httpPath"`
	CA              string `json:"ca" form:"ca" `
	Cert            string `json:"cert" form:"cert"`
	PrivateKey      string `json:"privateKey" form:"privateKey"`
	Passphrase      string `json:"passphrase" form:"passphrase"`
}

type UpdateUplink struct {
	Address      string `json:"address" form:"address" `
	MQTTUser     string `json:"mqttUser" form:"mqttUser"`
	MQTTPassword string `json:"mqttPassword" form:"mqttPassword"`
	HTTPMethod   string `json:"httpMethod" form:"httpMethod"`
	HTTPPath     string `json:"httpPath" form:"httpPath"`
	CA           string `json:"ca" form:"ca" `
	Cert         string `json:"cert" form:"cert"`
	PrivateKey   string `json:"privateKey" form:"privateKey"`
	Passphrase   string `json:"passphrase" form:"passphrase"`
}

type DeviceUplinkView struct {
	NodeName        string    `json:"nodeName"`
	Protocol        string    `json:"protocol"`
	Destination     string    `json:"destination"`
	DestinationName string    `json:"destinationName"`
	Address         string    `json:"address"`
	MQTTUser        string    `json:"mqttUser"`
	MQTTPassword    string    `json:"mqttPassword"`
	HTTPMethod      string    `json:"httpMethod"`
	HTTPPath        string    `json:"httpPath"`
	CA              string    `json:"ca" `
	Cert            string    `json:"cert"`
	PrivateKey      string    `json:"privateKey"`
	Passphrase      string    `json:"passphrase" form:"passphrase"`
	CreateTime      time.Time `json:"createTime"`
	UpdateTime      time.Time `json:"updateTime"`
}
