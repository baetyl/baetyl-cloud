package models

type Function struct {
	Name    string       `yaml:"name,omitempty" json:"name,omitempty" binding:"omitempty,res_name,nonbaetyl"`
	Handler string       `yaml:"handler,omitempty" json:"handler,omitempty"`
	Version string       `yaml:"version,omitempty" json:"version,omitempty"`
	Runtime string       `yaml:"runtime,omitempty" json:"runtime,omitempty"`
	Code    FunctionCode `yaml:"code,omitempty" json:"code,omitempty"`
}

type FunctionView struct {
	Functions []Function `json:"functions"`
}

type FunctionSourceView struct {
	Sources  []FunctionSource  `json:"sources"`
	Runtimes map[string]string `json:"runtimes"`
}

type FunctionSource struct {
	Name string `json:"name,omitempty"`
}

type FunctionCode struct {
	Size     int32  `yaml:"size,omitempty" json:"size,omitempty"`
	Sha256   string `yaml:"sha256,omitempty" json:"sha256,omitempty"`
	Location string `yaml:"location,omitempty" json:"location,omitempty"`
}
