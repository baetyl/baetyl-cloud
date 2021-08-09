package service

import (
	"github.com/baetyl/baetyl-go/v2/errors"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/cron.go -package=service github.com/baetyl/baetyl-cloud/v2/service CronService

type CronService interface {
	GetCron(name, namespace string) (*models.Cron, error)
	CreateCron(*models.Cron) error
	UpdateCron(*models.Cron) error
	DeleteCron(name, namespace string) error
	ListExpiredApps() ([]models.Cron, error)
	DeleteExpiredApps([]uint64) error
}

type cronService struct {
	plugin.Cron
}

func NewCronService(config *config.CloudConfig) (CronService, error) {
	cron, err := plugin.GetPlugin(config.Plugin.Cron)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &cronService{
		cron.(plugin.Cron),
	}, nil
}
