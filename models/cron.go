package models

import "time"

type Cron struct {
	Id        uint64    `json:"id,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	Name      string    `json:"name,omitempty" validate:"resourceName"`
	Selector  string    `json:"selector,omitempty"`
	CronTime  time.Time `json:"cronTime,omitempty"`
}
