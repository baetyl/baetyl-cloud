package entities

import "time"

type CronApp struct {
	Id         uint64    `db:"id"`
	Namespace  string    `db:"namespace"`
	Name       string    `db:"name"`
	Selector   string    `db:"selector"`
	CronTime   time.Time `db:"cron_time"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}
