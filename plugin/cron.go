package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/cron.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Cron

type Cron interface {
	GetCron(name, namespace string) (*models.Cron, error)
	CreateCron(*models.Cron) error
	UpdateCron(*models.Cron) error
	DeleteCron(name, namespace string) error
	ListExpiredApps() ([]models.Cron, error)
	DeleteExpiredApps([]uint64) error
	io.Closer
}
