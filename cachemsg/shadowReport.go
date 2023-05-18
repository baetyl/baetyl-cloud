package cachemsg

import (
	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type shadowReport struct {
}

var ShadowReportCache shadowReport

func init() {
	ShadowReportCache = shadowReport{}
}

// one report  64kb * 5000=  312M < one redis key value 512M
const maxLenCacheSize = 5000

func (s *shadowReport) SetAllShadowReportCache(cache plugin.DataCache, cacheData map[string]string) (err error) {
	log.L().Info("start cache set")
	var keys []string
	temp := map[string]string{}
	i := 0
	keysNum := 0
	//less than maxLenCacheSize
	if len(cacheData) < maxLenCacheSize {
		key := AllShadowReportCacheKey(keysNum)
		keys = append(keys, key)
		err = s.setShadowReportCache(cache, key, cacheData)
		if err != nil {
			return err
		}
	} else {
		// every maxLenCacheSize save one cache
		for key, value := range cacheData {
			temp[key] = value
			i++
			if i > maxLenCacheSize {
				key := AllShadowReportCacheKey(keysNum)
				err = s.setShadowReportCache(cache, key, temp)
				if err != nil {
					return err
				}
				temp = map[string]string{}
				keys = append(keys, key)
				keysNum++
				i = 0
			}
		}
		// save last cache
		if len(temp) > 0 {
			key := AllShadowReportCacheKey(keysNum)
			keys = append(keys, key)
			err = s.setShadowReportCache(cache, key, cacheData)
			if err != nil {
				return err
			}
		}
	}
	//save keys
	err = s.setShadowReportCacheKey(cache, keys)
	if err != nil {
		return err
	}
	return nil
}

func (s *shadowReport) GetAllShadowReportCache(cache plugin.DataCache) (map[string]string, error) {
	var keys []string
	data, err := cache.GetByte(AllShadowReportCacheKeys)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, err
	}
	returnData := map[string]string{}
	for i := range keys {
		tempData, err := cache.GetByte(keys[i])
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(tempData, &returnData)
		if err != nil {
			return nil, err
		}

	}
	return returnData, err

}

func (s *shadowReport) UpdateShadowReportCache(cache plugin.DataCache, node, reportStr string) error {
	data, err := s.GetAllShadowReportCache(cache)
	if err != nil {
		return err
	}
	data[node] = reportStr
	err = s.SetAllShadowReportCache(cache, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *shadowReport) DeleteShadowReportCache(cache plugin.DataCache, node string) error {
	data, err := s.GetAllShadowReportCache(cache)
	if err != nil {
		return err
	}
	delete(data, node)
	err = s.SetAllShadowReportCache(cache, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *shadowReport) setShadowReportCache(cache plugin.DataCache, key string, data map[string]string) error {
	marshalData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = cache.SetByte(key, marshalData)
	if err != nil {
		return err
	}
	return nil
}

func (s *shadowReport) setShadowReportCacheKey(cache plugin.DataCache, keys []string) error {
	marshalData, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	err = cache.SetByte(AllShadowReportCacheKeys, marshalData)
	if err != nil {
		return err
	}
	return nil
}
