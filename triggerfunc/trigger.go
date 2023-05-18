package triggerfunc

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/cachemsg"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

const (
	ShadowCreateOrUpdateTrigger = "shadowCreateOrUpdateTrigger"
	ShadowDelete                = "shadowDelete"
)

// ShadowCreateOrUpdateCacheSet update shadow cache when shadow update or create
func ShadowCreateOrUpdateCacheSet(cache plugin.DataCache, shadow models.Shadow) {
	reportTimeData, err := cache.Get(cachemsg.AllShadowReportTimeCache)
	if err != nil {
		log.L().Error("get shadow cache error", log.Error(err))
		return
	}
	reportTimeMap := map[string]string{}
	err = json.Unmarshal([]byte(reportTimeData), &reportTimeMap)
	if err != nil {
		log.L().Error("unmarshal err", log.Error(err))
		return
	}
	reportTimeMap[shadow.Name] = shadow.Time.Format(time.RFC3339Nano)

	returnTimeData, err := json.Marshal(reportTimeMap)
	if err != nil {
		log.L().Error("Marshal err", log.Error(err))
		return
	}
	err = cache.Set(cachemsg.AllShadowReportTimeCache, string(returnTimeData))
	if err != nil {
		log.L().Error("set report time err", log.Error(err))
		return
	}

	err = cachemsg.ShadowReportCache.UpdateShadowReportCache(cache, shadow.Name, shadow.ReportStr)
	if err != nil {
		log.L().Error("update report  err", log.Error(err))
		return
	}
}

// ShadowDeleteCache delete shadow cache when shadow delete
func ShadowDeleteCache(cache plugin.DataCache, name string) {
	reportTimeData, err := cache.Get(cachemsg.AllShadowReportTimeCache)
	if err != nil {
		log.L().Error("get shadow cache error", log.Error(err))
		return
	}
	reportTimeMap := map[string]string{}
	err = json.Unmarshal([]byte(reportTimeData), &reportTimeMap)
	if err != nil {
		log.L().Error("unmarshal err", log.Error(err))
		return
	}
	delete(reportTimeMap, name)
	returnTimeData, err := json.Marshal(reportTimeMap)
	if err != nil {
		log.L().Error("Marshal err", log.Error(err))
		return
	}
	err = cache.Set(cachemsg.AllShadowReportTimeCache, string(returnTimeData))
	if err != nil {
		log.L().Error("set report time err", log.Error(err))
		return
	}

	err = cachemsg.ShadowReportCache.DeleteShadowReportCache(cache, name)
	if err != nil {
		log.L().Error("delete report  err", log.Error(err))
		return
	}

}
