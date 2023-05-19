package cachemsg

import "fmt"

const (
	// AllShadowReportTimeCache set report time , type is map. ex :shadow-time : "{ "aaa": "0001-01-01T00:00:00Z", "d-33949349": "0001-01-01T00:00:00Z"}"
	AllShadowReportTimeCache = "shadow-time"
	// ShadowReportDataCache set report cache  ex: shadow-aaa-report : "{"apps": []}"
	ShadowReportDataCache = "shadow-%s-report"
	// CacheReportSetLock set report cache running flag key
	CacheReportSetLock = "cache-report-lock"
)

func GetShadowReportCacheKey(nodeName string) string {
	return fmt.Sprintf(ShadowReportDataCache, nodeName)
}
