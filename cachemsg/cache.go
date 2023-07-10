package cachemsg

import "fmt"

const (
	// AllShadowReportTimeCache set report time , type is map. ex :shadow-time : "{ "aaa": "0001-01-01T00:00:00Z", "d-33949349": "0001-01-01T00:00:00Z"}"
	AllShadowReportTimeCache = "shadow-%s-time"
	// ShadowReportDataCache set report cache  ex: shadow-aaa-report : "{"apps": []}"
	ShadowReportDataCache = "shadow-%s-%s-report"
	// CacheReportSetLock set report cache running flag key
	CacheReportSetLock = "cache-report-lock"
	// CacheUpdateReportTimeLock set update report time cache running flag key
	CacheUpdateReportTimeLock = "cache-report-time-lock"
)

// GetShadowReportTimeCacheKey GetShadowReportTime get namesapce report time key
func GetShadowReportTimeCacheKey(namespace string) string {
	return fmt.Sprintf(AllShadowReportTimeCache, namespace)
}

// GetShadowReportCacheKey get node report key
func GetShadowReportCacheKey(namespace, nodeName string) string {
	return fmt.Sprintf(ShadowReportDataCache, namespace, nodeName)
}
