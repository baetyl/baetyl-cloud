// Package models 模型定义
package models

import specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

type APPDuplicate struct {
	specV1.Application `json:",inline"`
	Configs            map[string]*specV1.Configuration `json:"configs,omitempty"`
	Secrets            map[string]*specV1.Secret        `json:"secrets,omitempty"`
}
