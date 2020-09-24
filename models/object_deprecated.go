package models

type ObjectStorageSourceView struct {
	Sources []ObjectStorageSource `json:"sources"`
}

type ObjectStorageSource struct {
	Name string `json:"name,omitempty"`
}
