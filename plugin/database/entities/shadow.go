package entities

import (
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Shadow struct {
	Id         int64     `db:"id"`
	Namespace  string    `db:"namespace"`
	Name       string    `db:"name"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
	Report     string    `db:"report"`
	Desire     string    `db:"desire"`
	ReportMeta string    `db:"report_meta"`
	DesireMeta string    `db:"desire_meta"`
}

func (s *Shadow) ToShadowModel() (*models.Shadow, error) {
	shadow := &models.Shadow{
		Namespace:         s.Namespace,
		Name:              s.Name,
		CreationTimestamp: s.CreateTime.UTC(),
		Report:            models.BuildEmptyApps(),
		Desire:            models.BuildEmptyApps(),
	}

	report := v1.Report{}
	err := json.Unmarshal([]byte(s.Report), &report)
	if err != nil {
		return nil, err
	}
	shadow.Report = report

	desire := v1.Desire{}
	err = json.Unmarshal([]byte(s.Desire), &desire)
	if err != nil {
		return nil, err
	}
	shadow.Desire = desire

	var reportMeta map[string]interface{}
	if s.ReportMeta != "" {
		if err := json.Unmarshal([]byte(s.ReportMeta), &reportMeta); err != nil {
			return nil, err
		}
	}
	shadow.ReportMeta = reportMeta
	var desireMeta  map[string]interface{}
	if s.DesireMeta != "" {
		if err := json.Unmarshal([]byte(s.DesireMeta), &desireMeta); err != nil {
			return nil, err
		}
	}
	shadow.DesireMeta = desireMeta
	return shadow, nil
}

func NewShadowFromShadowModel(shadow *models.Shadow) (*Shadow, error) {
	shd := new(Shadow)
	shd.Name = shadow.Name
	shd.Namespace = shadow.Namespace

	report, err := json.Marshal(shadow.Report)
	if err != nil {
		return nil, err
	}
	shd.Report = string(report)

	desire, err := json.Marshal(shadow.Desire)
	if err != nil {
		return nil, err
	}
	shd.Desire = string(desire)

	reportMeta, err := json.Marshal(shadow.ReportMeta)
	if err != nil {
		return nil, err
	}
	shd.ReportMeta = string(reportMeta)

	desireMeta, err := json.Marshal(shadow.DesireMeta)
	if err != nil {
		return nil, err
	}
	shd.DesireMeta = string(desireMeta)
	return shd, nil
}
