package models

type ObjectStorageSourceViewV1 struct {
	Sources []ObjectStorageSourceV1 `json:"sources"`
}

type ObjectStorageSourceV1 struct {
	Name string `json:"name,omitempty"`
}
