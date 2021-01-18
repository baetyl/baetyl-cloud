package models

import (
	"reflect"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

// NodeViewList node view list
type NodeViewList struct {
	Total       int               `json:"total"`
	ListOptions *ListOptions      `json:"listOptions"`
	Items       []specV1.NodeView `json:"items"`
}

// NodeList node list
type NodeList struct {
	Total       int           `json:"total"`
	ListOptions *ListOptions  `json:"listOptions"`
	Items       []specV1.Node `json:"items"`
}

type NodeNames struct {
	Names []string `json:"names,"validate:"maxLength=20"`
}

type NodeProperties struct {
	State NodePropertiesState    `yaml:"state,omitempty" json:"state,omitempty"`
	Meta  NodePropertiesMetadata `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type NodePropertiesState struct {
	Report map[string]interface{} `yaml:"report,omitempty" json:"report,omitempty"`
	Desire map[string]interface{} `yaml:"desire,omitempty" json:"desire,omitempty"`
}

type NodeProperty struct {
	Name    string        `yaml:"name,omitempty" json:"name,omitempty"`
	Type    string        `yaml:"type,omitempty" json:"type,omitempty"`
	Current PropertyValue `yaml:"current,omitempty" json:"current,omitempty"`
	Expect  PropertyValue `yaml:"expect,omitempty" json:"expect,omitempty"`
}

type PropertyValue struct {
	Value      string `yaml:"value,omitempty" json:"value,omitempty"`
	UpdateTime string `yaml:"updateTime,omitempty" json:"updateTime,omitempty"`
}

type NodeMode struct {
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
}

type NodePropertiesMetadata struct {
	ReportMeta map[string]interface{} `yaml:"report,omitempty" json:"report,omitempty"`
	DesireMeta map[string]interface{} `yaml:"desire,omitempty" json:"desire,omitempty"`
}

func EqualNode(node1, node2 *specV1.Node) bool {
	return reflect.DeepEqual(node1.Labels, node2.Labels) &&
		reflect.DeepEqual(node1.Description, node2.Description) &&
		reflect.DeepEqual(node1.Annotations, node2.Annotations) &&
		reflect.DeepEqual(node1.SysApps, node2.SysApps)
}

func PopulateNode(node *specV1.Node) error {
	val, ok := node.Attributes[specV1.KeyOptionalSysApps]
	if !ok || val == nil {
		return nil
	}

	ss, ok := val.([]interface{})
	if !ok {
		return common.Error(common.ErrConvertConflict, common.Field("name", specV1.KeyOptionalSysApps), common.Field("error", "failed to interface{} to []interface{}`"))
	}

	for _, d := range ss {
		s, ok := d.(string)
		if !ok {
			return common.Error(common.ErrConvertConflict, common.Field("name", specV1.KeyOptionalSysApps), common.Field("error", "failed to interface{} to string`"))
		}
		node.SysApps = append(node.SysApps, s)
	}
	return nil
}

type NodeCoreConfigs struct {
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
	// unit: seconds
	Frequency int `yaml:"frequency,omitempty" json:"frequency,omitempty"`
	APIPort   int `yaml:"apiport,omitempty" json:"apiport,omitempty"`
}

type NodeCoreVersions struct {
	Versions []string `yaml:"versions,omitempty" json:"versions,omitempty"`
}

type NodeSysAppView struct {
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type NodeOptionalSysApps struct {
	Apps []NodeSysAppView `yaml:"apps,omitempty" json:"apps,omitempty"`
}

type NodeSysAppInfo struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Programs    map[string]string `yaml:"programs,omitempty" json:"programs,omitempty"`
}
