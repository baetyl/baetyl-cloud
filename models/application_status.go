package models

import (
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
)

type AppStatusParams struct {
	App []AppStatusParamsItem `json:"appParams" binding:"required"`
}

type AppStatusParamsItem struct {
	Name  string   `json:"name" binding:"required"`
	Nodes []string `json:"nodes" binding:"required"`
}

type AppNodeStatusReturn struct {
	Items map[string]AppNodeStatusItem `json:"items"`
}

type AppNodeStatusItem struct {
	AppInfo    *AppItem                    `json:"app_info"`
	NodeReport map[string]*specV1.AppStats `json:"node_report"`
}
