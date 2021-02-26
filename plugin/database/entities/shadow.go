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

	return shd, nil
}
