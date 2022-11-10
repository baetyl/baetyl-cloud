package models

// Namespace Namespace
type Namespace struct {
	Name string `json:"name,omitempty" binding:"namespace"`
}

// NamespaceList namespace list
type NamespaceList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []Namespace `json:"items"`
}
