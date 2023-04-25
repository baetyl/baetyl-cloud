package cacheMsg

import "fmt"

const (
	ShadowReportTImeCache = "shadow-%s-time"
	ShadowReportCache     = "shadow-%s-report"
)

func GetShadowReportTimeCacheKey(nodeName string) string {
	return fmt.Sprintf(ShadowReportTImeCache, nodeName)
}

func GetShadowReportCacheKey(nodeName string) string {
	return fmt.Sprintf(ShadowReportCache, nodeName)
}
