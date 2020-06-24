package models

// Namespace Namespace
type Namespace struct {
	Name string `json:"name,omitempty" validate:"namespace"`
}
