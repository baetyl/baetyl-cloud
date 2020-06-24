package entities

import (
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
	"time"
)

func genBatch() *models.Batch {
	return &models.Batch{
		Name:            "zx",
		Namespace:       "default",
		Description:     "desc",
		QuotaNum:        20,
		EnableWhitelist: 0,
		SecurityType:    common.None,
		SecurityKey:     "",
		CallbackName:    "test",
		CreateTime:      time.Unix(1000, 10),
		UpdateTime:      time.Unix(1000, 10),
		Labels:          map[string]string{"a": "a"},
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
		QuotaNum:        20,
		EnableWhitelist: 0,
		SecurityType:    common.None,
		SecurityKey:     "",
		CallbackName:    "test",
		CreateTime:      time.Unix(1000, 10),
		UpdateTime:      time.Unix(1000, 10),
		Labels:          "{\"a\":\"a\"}",
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
