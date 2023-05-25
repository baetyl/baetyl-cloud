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

// reportSaveStruct save temp need update data
// @time control the updated content during high concurrency
type reportSaveStruct struct {
	Value string
	Time  time.Time
}

var reportSaveMap = map[string]reportSaveStruct{}

func init() {
	reportSaveMap = map[string]reportSaveStruct{}
	t = time.NewTimer(time.Second)
}

var t *time.Timer

// ShadowCreateOrUpdateCacheSet update shadow cache when shadow update or create
func ShadowCreateOrUpdateCacheSet(cache plugin.DataCache, shadow models.Shadow) {
	reportSaveMap[shadow.Name] = reportSaveStruct{
		Value: shadow.Time.Format(time.RFC3339Nano),
		Time:  time.Now(),
	}
	select {
	case <-t.C:
		saveCache(cache, time.Now())
		t.Reset(time.Second)
	default:
	}
	err := cache.SetByte(cachemsg.GetShadowReportCacheKey(shadow.Name), []byte(shadow.ReportStr))
	if err != nil {
		log.L().Error("update report  err", log.Error(err))
		return
	}
}

// ShadowDeleteCache delete shadow cache when shadow delete
func ShadowDeleteCache(cache plugin.DataCache, name string) {
	reportTimeData, err := cache.GetByte(cachemsg.AllShadowReportTimeCache)
	if err != nil {
		log.L().Error("get shadow cache error", log.Error(err))
		return
	}
	if reportTimeData == nil {
		return
	}
	reportTimeMap := map[string]string{}
	err = json.Unmarshal(reportTimeData, &reportTimeMap)
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
	err = cache.SetByte(cachemsg.AllShadowReportTimeCache, returnTimeData)
	if err != nil {
		log.L().Error("delete report time err", log.Error(err))
		return
	}

	err = cache.Delete(cachemsg.GetShadowReportCacheKey(name))
	if err != nil {
		log.L().Error("delete report  err", log.Error(err))
		return
	}
}

func saveCache(cache plugin.DataCache, timeCheck time.Time) {
	reportTimeData, err := cache.GetByte(cachemsg.AllShadowReportTimeCache)
	if err != nil {
		log.L().Error("get shadow cache error", log.Error(err))
		return
	}
	reportTimeMap := map[string]string{}
	if reportTimeData != nil {
		err = json.Unmarshal(reportTimeData, &reportTimeMap)
		if err != nil {
			log.L().Error("unmarshal err", log.Error(err))
			return
		}
	}
	for key, val := range reportSaveMap {
		//before timeCheck data update
		if val.Time.Before(timeCheck) {
			reportTimeMap[key] = val.Value
			delete(reportSaveMap, key)
		}
	}
	returnTimeData, err := json.Marshal(reportTimeMap)
	if err != nil {
		log.L().Error("Marshal err", log.Error(err))
		return
	}
	err = cache.SetByte(cachemsg.AllShadowReportTimeCache, returnTimeData)
	if err != nil {
		log.L().Error("delete report time err", log.Error(err))
		return
	}
}
