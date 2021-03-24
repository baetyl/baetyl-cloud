package service

import (
	"context"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

//go:generate mockgen -destination=../mock/service/locker.go -package=service github.com/baetyl/baetyl-cloud/v2/service LockerService

type LockerService interface {
	Lock(ctx context.Context, name string, ttl int64) (string, error)
	Unlock(ctx context.Context, name, value string)
}

// NewModuleService
func NewLockerService(config *config.CloudConfig) (LockerService, error) {
	locker, err := plugin.GetPlugin(config.Plugin.Locker)
	if err != nil {
		return nil, err
	}
	return locker.(plugin.Locker), nil
}
