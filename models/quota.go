package models

// Quota
type Quota struct {
	Namespace string `json:"namespace"`
	QuotaName string `json:"quotaName"`
	Quota     int    `json:"quota"`
	UsedNum   int    `json:"usedNum"`
}
