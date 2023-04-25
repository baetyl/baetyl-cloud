package triggerFunc

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/cacheMsg"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

const (
	ShadowCreateOrUpdateTrigger = "shadowCreateOrUpdateTrigger"
	ShadowDelete                = "shadowDelete"
)

// ShadowCreateOrUpdateCacheSet update shadow cache when shadow update or create
func ShadowCreateOrUpdateCacheSet(cache plugin.DataCache, shadow models.Shadow) {
	err := cache.Set(cacheMsg.GetShadowReportTimeCacheKey(shadow.Name), shadow.Time.Format(time.RFC3339Nano))
	if err != nil {
		log.L().Error("set shadow cache error", log.Error(err))
	}
	err = cache.Set(cacheMsg.GetShadowReportCacheKey(shadow.Name), shadow.ReportStr)
	if err != nil {
		log.L().Error("set shadow cache error", log.Error(err))
	}

}

// ShadowDeleteCache delete shadow cache when shadow delete
func ShadowDeleteCache(cache plugin.DataCache, name string) {
	err := cache.Delete(cacheMsg.GetShadowReportTimeCacheKey(name))
	if err != nil {
		log.L().Error("delete shadow cache error", log.Error(err))
	}
	err = cache.Delete(cacheMsg.GetShadowReportCacheKey(name))
	if err != nil {
		log.L().Error("delete shadow cache error", log.Error(err))
	}
}
