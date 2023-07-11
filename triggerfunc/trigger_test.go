package triggerfunc

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/cachemsg"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/plugin/cache/localcache"
)

func TestCache(t *testing.T) {
	conf := `
freeCacheConfig:
 maxBytes: 1024

`
	filename := "cloud.yml"
	err := ioutil.WriteFile(filename, []byte(conf), 0644)
	assert.NoError(t, err)
	defer os.Remove(filename)
	common.SetConfFile(filename)

	p, err := localcache.New()
	assert.NoError(t, err)
	assert.NotNil(t, p)

	cache := p.(plugin.DataCache)

	ShadowCreateOrUpdateCacheSet(cache, models.Shadow{
		Namespace: "default",
		Name:      "aaa",
		Time:      time.Now(),
	})
	ShadowCreateOrUpdateCacheSet(cache, models.Shadow{
		Namespace: "default",
		Name:      "bbb",
		Time:      time.Now(),
	})
	time.Sleep(time.Second * 1)
	ShadowCreateOrUpdateCacheSet(cache, models.Shadow{
		Namespace: "default",
		Name:      "aaa",
		Time:      time.Now(),
	})
	time.Sleep(time.Second * 1)
	ShadowCreateOrUpdateCacheSet(cache, models.Shadow{
		Namespace: "default",
		Name:      "bbb",
		Time:      time.Now(),
	})

	dataReportTime, err := cache.GetByte(cachemsg.GetShadowReportTimeCacheKey("default"))
	assert.NoError(t, err)
	if dataReportTime != nil {
		reportTimeMap := map[string]string{}
		err = json.Unmarshal(dataReportTime, &reportTimeMap)
		assert.NotNil(t, reportTimeMap["aaa"])
		assert.NotNil(t, reportTimeMap["bbb"])
	}

	ShadowDeleteCache(cache, "bbb", "default")

	dataReportTime, err = cache.GetByte(cachemsg.GetShadowReportTimeCacheKey("default"))
	assert.NoError(t, err)
	if dataReportTime != nil {
		reportTimeMap := map[string]string{}
		err = json.Unmarshal(dataReportTime, &reportTimeMap)
		assert.NotNil(t, reportTimeMap["aaa"])
		assert.Equal(t, len(reportTimeMap), 1)
	}
	cache.Close()

}
