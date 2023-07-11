package models

type YamlResourceList struct {
	Total int           `json:"total"`
	Items []interface{} `json:"items"`
}
