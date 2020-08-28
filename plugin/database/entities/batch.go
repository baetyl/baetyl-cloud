package entities

import (
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Batch struct {
	Name            string          `db:"name"`
	Namespace       string          `db:"namespace"`
	Description     string          `db:"description"`
	QuotaNum        int             `db:"quota_num"`
	EnableWhitelist int             `db:"enable_whitelist"`
	SecurityType    common.Security `db:"security_type"`
	SecurityKey     string          `db:"security_key"`
	CallbackName    string          `db:"callback_name"`
	CreateTime      time.Time       `db:"create_time"`
	UpdateTime      time.Time       `db:"update_time"`
	Labels          string          `db:"labels"`
	Fingerprint     string          `db:"fingerprint"`
}

func ToBatchModel(batch *Batch) *models.Batch {
	var labels map[string]string
	if err := json.Unmarshal([]byte(batch.Labels), &labels); err != nil {
		log.L().Error("batch db labels unmarshal error", log.Any("labels", batch.Labels))
	}
	var fp models.Fingerprint
	if err := json.Unmarshal([]byte(batch.Fingerprint), &fp); err != nil {
		log.L().Error("batch db fingerprint unmarshal error", log.Any("fingerprint", batch.Fingerprint))
	}
	res := &models.Batch{
		Name:            batch.Name,
		Namespace:       batch.Namespace,
		Description:     batch.Description,
		QuotaNum:        batch.QuotaNum,
		EnableWhitelist: batch.EnableWhitelist,
		SecurityType:    batch.SecurityType,
		SecurityKey:     batch.SecurityKey,
		CallbackName:    batch.CallbackName,
		CreateTime:      batch.CreateTime,
		UpdateTime:      batch.UpdateTime,
		Labels:          labels,
		Fingerprint:     fp,
	}
	return res
}

func FromBatchModel(batch *models.Batch) *Batch {
	labels, err := json.Marshal(batch.Labels)
	if err != nil {
		log.L().Error("batch labels marshal error", log.Any("labels", batch.Labels))
		labels = []byte("{}")
	}
	fingerprint, err := json.Marshal(batch.Fingerprint)
	if err != nil {
		log.L().Error("batch fingerprint marshal error", log.Any("fingerprint", batch.Fingerprint))
		fingerprint = []byte("{}")
	}
	res := &Batch{
		Name:            batch.Name,
		Namespace:       batch.Namespace,
		Description:     batch.Description,
		QuotaNum:        batch.QuotaNum,
		EnableWhitelist: batch.EnableWhitelist,
		SecurityType:    batch.SecurityType,
		SecurityKey:     batch.SecurityKey,
		CallbackName:    batch.CallbackName,
		CreateTime:      batch.CreateTime,
		UpdateTime:      batch.UpdateTime,
		Labels:          string(labels),
		Fingerprint:     string(fingerprint),
	}
	return res
}
