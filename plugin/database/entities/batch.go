// Package entities 数据库存储基本结构与方法
package entities

import (
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/common"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Batch struct {
	Name            string          `db:"name"`
	Namespace       string          `db:"namespace"`
	Description     string          `db:"description"`
	Accelerator     string          `db:"accelerator"`
	SysApps         string          `db:"sys_apps"`
	QuotaNum        string          `db:"quota_num"`
	EnableWhitelist int             `db:"enable_whitelist"`
	Cluster         int             `db:"cluster"`
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
	var sysApps []string
	if err := json.Unmarshal([]byte(batch.SysApps), &sysApps); err != nil {
		log.L().Error("batch db sysApps unmarshal error", log.Any("sysApps", batch.SysApps))
	}
	res := &models.Batch{
		Name:            batch.Name,
		Namespace:       batch.Namespace,
		Description:     batch.Description,
		QuotaNum:        batch.QuotaNum,
		EnableWhitelist: batch.EnableWhitelist,
		Cluster:         batch.Cluster,
		SecurityType:    batch.SecurityType,
		SecurityKey:     batch.SecurityKey,
		CallbackName:    batch.CallbackName,
		CreateTime:      batch.CreateTime,
		UpdateTime:      batch.UpdateTime,
		Labels:          labels,
		Fingerprint:     fp,
		Accelerator:     batch.Accelerator,
		SysApps:         sysApps,
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
	sysApps, err := json.Marshal(batch.SysApps)
	if err != nil {
		log.L().Error("batch sysApps marshal error", log.Any("sysApps", batch.SysApps))
		sysApps = []byte("[]")
	}
	res := &Batch{
		Name:            batch.Name,
		Namespace:       batch.Namespace,
		Description:     batch.Description,
		QuotaNum:        batch.QuotaNum,
		EnableWhitelist: batch.EnableWhitelist,
		Cluster:         batch.Cluster,
		SecurityType:    batch.SecurityType,
		SecurityKey:     batch.SecurityKey,
		CallbackName:    batch.CallbackName,
		CreateTime:      batch.CreateTime,
		UpdateTime:      batch.UpdateTime,
		Labels:          string(labels),
		Fingerprint:     string(fingerprint),
		Accelerator:     batch.Accelerator,
		SysApps:         string(sysApps),
	}
	return res
}
