package models

type MisData struct {
	Count int         `json:"count,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
}
