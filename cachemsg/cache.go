package cachemsg

import "fmt"

const (
	AllShadowReportTimeCache = "shadow-time"
	ShadowReportDataCache    = "shadow-%s-report"
)

var CacheReportSetLock = false

func GetShadowReportCacheKey(nodeName string) string {
	return fmt.Sprintf(ShadowReportDataCache, nodeName)
}
