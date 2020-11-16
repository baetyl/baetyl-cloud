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
	var report specV1.Report
	err := msg.Content.Unmarshal(&report)
	if err != nil {
		return nil, err
	}

	setNodeAddressIfExist(msg, &report)
	delta, err := s.Sync.Report(msg.Metadata["namespace"], msg.Metadata["name"], report)
	if err != nil {
		return nil, err
	}
	return &specV1.Message{
		Kind:     specV1.MessageReport,
		Metadata: msg.Metadata,
		Content:  specV1.LazyValue{Value: delta},
	}, nil
}

// Desire for node synchronize desire info
func (s *SyncAPIImpl) Desire(msg specV1.Message) (*specV1.Message, error) {
	var desireRes specV1.DesireRequest
	err := msg.Content.Unmarshal(&desireRes)
	if err != nil {
		return nil, err
	}

	res, err := s.Sync.Desire(msg.Metadata["namespace"], desireRes.Infos, msg.Metadata)
	if err != nil {
		return nil, err
	}
	return &specV1.Message{
		Kind:     specV1.MessageDesire,
		Metadata: msg.Metadata,
		Content:  specV1.LazyValue{Value: specV1.DesireResponse{Values: res}},
	}, nil
}

func setNodeAddressIfExist(msg specV1.Message, report *specV1.Report) {
	if addr, ok := msg.Metadata["address"]; !ok {
		return
	} else {
		nodeVal, ok := (*report)["node"]
		if ok {
			node, ok := nodeVal.(map[string]interface{})
			if ok {
				node["address"] = addr
				(*report)["node"] = node
			}
		}
	}
}
