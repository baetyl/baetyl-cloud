package cachemsg

import "fmt"

const (
	AllShadowReportTimeCache = "shadow-time"
	AllShadowReportCache     = "shadow-report-%d"
	AllShadowReportCacheKeys = "shadow-report-keys"
)

func AllShadowReportCacheKey(num int) string {
	return fmt.Sprintf(AllShadowReportCache, num)
}
