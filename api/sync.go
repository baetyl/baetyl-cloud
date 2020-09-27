package api

import (
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/service"
)

//go:generate mockgen -destination=../mock/api/sync.go -package=api github.com/baetyl/baetyl-cloud/v2/api SyncAPI

type SyncAPI interface {
	Report(msg specV1.Message) (*specV1.Message, error)
	Desire(msg specV1.Message) (*specV1.Message, error)
}

type SyncAPIImpl struct {
	Sync service.SyncService
}

func NewSyncAPI(cfg *config.CloudConfig) (SyncAPI, error) {
	syncService, err := service.NewSyncService(cfg)
	if err != nil {
		return nil, err
	}
	return &SyncAPIImpl{
		Sync: syncService,
	}, nil
}

// Report for node report
func (s *SyncAPIImpl) Report(msg specV1.Message) (*specV1.Message, error) {
	if msg.Content.Value == nil {
		msg.Content.Value = &specV1.Report{}
		if err := msg.Content.Unmarshal(msg.Content.Value); err != nil {
			return nil, err
		}
	}

	desire, err := s.Sync.Report(msg.Metadata["namespace"], msg.Metadata["name"], *msg.Content.Value.(*specV1.Report))
	if err != nil {
		return nil, err
	}
	return &specV1.Message{
		Kind:     specV1.MessageReport,
		Metadata: msg.Metadata,
		Content:  specV1.VariableValue{Value: desire},
	}, nil
}

// Desire for node synchronize desire info
func (s *SyncAPIImpl) Desire(msg specV1.Message) (*specV1.Message, error) {
	if msg.Content.Value == nil {
		msg.Content.Value = &specV1.DesireRequest{}
		if err := msg.Content.Unmarshal(msg.Content.Value); err != nil {
			return nil, err
		}
	}

	res, err := s.Sync.Desire(msg.Metadata["namespace"], msg.Content.Value.(*specV1.DesireRequest).Infos, msg.Metadata)
	if err != nil {
		return nil, err
	}
	return &specV1.Message{
		Kind:     specV1.MessageDesire,
		Metadata: msg.Metadata,
		Content:  specV1.VariableValue{Value: specV1.DesireResponse{Values: res}},
	}, nil
}
