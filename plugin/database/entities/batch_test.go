// Package entities 数据库存储基本结构与方法
package entities

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func genBatch() *models.Batch {
	return &models.Batch{
		Name:            "zx",
		Namespace:       "default",
		Description:     "desc",
		QuotaNum:        "20",
		EnableWhitelist: 0,
		SecurityType:    common.None,
		SecurityKey:     "",
		CallbackName:    "test",
		CreateTime:      time.Unix(1000, 10),
		UpdateTime:      time.Unix(1000, 10),
		Labels:          map[string]string{"a": "a"},
		Accelerator:     "nvidia",
		SysApps:         []string{"a", "b"},
		Fingerprint: models.Fingerprint{
			Type:   common.FingerprintSN,
			SnPath: path.Join(common.DefaultSNPath, common.DefaultSNFile),
		},
	}
}

func genBatchDB() *Batch {
	return &Batch{
		Name:            "zx",
		Namespace:       "default",
		Description:     "desc",
		QuotaNum:        "20",
		EnableWhitelist: 0,
		SecurityType:    common.None,
		SecurityKey:     "",
		CallbackName:    "test",
		CreateTime:      time.Unix(1000, 10),
		UpdateTime:      time.Unix(1000, 10),
		Labels:          "{\"a\":\"a\"}",
		Accelerator:     "nvidia",
		SysApps:         "[\"a\",\"b\"]",
		Fingerprint:     "{\"type\":1,\"snPath\":\"/var/lib/baetyl/sn/fingerprint.txt\"}",
	}
}

func TestConvertBatch(t *testing.T) {
	batch := genBatch()
	batchDB := genBatchDB()
	resBatch := ToBatchModel(batchDB)
	assert.EqualValues(t, batch, resBatch)
	resBatchDB := FromBatchModel(batch)
	assert.EqualValues(t, batchDB, resBatchDB)
}
