package models

type DesireRequest struct {
	Resources []*Resource `yaml:"resources,omitempty" json:"resources,omitempty"`
}

type DesireResponse struct {
	Resources []*Resource `yaml:"resources,omitempty" json:"resources,omitempty"`
}

type ResponseInfo struct {
	Delta    interface{}            `yaml:"delta,omitempty" json:"delta,omitempty"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty" default:"{}"`
}

type Resource struct {
	Type    string      `yaml:"type,omitempty" json:"type,omitempty"`
	Name    string      `yaml:"name,omitempty" json:"name,omitempty"`
	Version string      `yaml:"version,omitempty" json:"version,omitempty"`
	Value   interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}
