package models

import (
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type Shadow struct {
	Namespace         string                 `json:"namespace,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Report            v1.Report              `json:"report,omitempty"`
	Desire            v1.Desire              `json:"desire,omitempty"`
	ReportMeta        map[string]interface{} `json:"reportMeta,omitempty"`
	DesireMeta        map[string]interface{} `json:"desireMeta,omitempty"`
	CreationTimestamp time.Time              `json:"createTime,omitempty"`
}

// NodeViewList node view list
type ShadowList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []Shadow `json:"items"`
}

func NewShadow(namespace, name string) *Shadow {
	return &Shadow{
		Name:       name,
		Namespace:  namespace,
		Report:     BuildEmptyApps(),
		Desire:     BuildEmptyApps(),
		ReportMeta: make(map[string]interface{}),
		DesireMeta: make(map[string]interface{}),
	}
}

func NewShadowFromNode(node *v1.Node) *Shadow {
	shadow := &Shadow{
		Name:              node.Name,
		Namespace:         node.Namespace,
		Report:            node.Report,
		Desire:            node.Desire,
		ReportMeta:        make(map[string]interface{}),
		DesireMeta:        make(map[string]interface{}),
		CreationTimestamp: node.CreationTimestamp.UTC(),
	}

	if node.Desire == nil {
		shadow.Desire = BuildEmptyApps()
	}

	if node.Report == nil {
		shadow.Report = BuildEmptyApps()
	}

	return shadow
}

func (s *Shadow) GetDesireString() (string, error) {
	desire, err := json.Marshal(s.Desire)
	if err != nil {
		return "", err
	}

	return string(desire), nil
}

func (s *Shadow) GetDesireMetaString() (string, error) {
	desireMeta, err := json.Marshal(s.DesireMeta)
	if err != nil {
		return "", err
	}
	return string(desireMeta), nil
}

func (s *Shadow) GetReportString() (string, error) {
	report, err := json.Marshal(s.Report)
	if err != nil {
		return "", err
	}

	return string(report), nil
}

func (s *Shadow) GetReportMetaString() (string, error) {
	reportMeta, err := json.Marshal(s.ReportMeta)
	if err != nil {
		return "", err
	}
	return string(reportMeta), nil
}

func BuildEmptyApps() map[string]interface{} {
	return map[string]interface{}{
		common.DesiredApplications:    make([]v1.AppInfo, 0),
		common.DesiredSysApplications: make([]v1.AppInfo, 0),
	}
}
