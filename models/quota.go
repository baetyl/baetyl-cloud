package models

// Quota
type Quota struct {
	Namespace string `json:"namespace" binding:"required"`
	QuotaName string `json:"quotaName,omitempty"`
	Quota     int    `json:"quota" default:"0"`
	UsedNum   int    `json:"usedNum" default:"0"`
}
