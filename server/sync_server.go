package server

import (
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/config"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type HandlerReport func(msg specV1.Message) (*specV1.Message, error)
type HandlerDesire func(msg specV1.Message) (*specV1.Message, error)

type SyncServer struct {
	links   map[string]plugin.SyncLink
	syncAPI api.SyncAPI
}

func NewSyncServer(cfg *config.CloudConfig) (*SyncServer, error) {
	sync := &SyncServer{
		links: map[string]plugin.SyncLink{},
	}
	for _, l := range cfg.Plugin.SyncLinks {
		link, err := plugin.GetPlugin(l)
		if err != nil {
			return nil, err
		}
		sync.links[l] = link.(plugin.SyncLink)
	}
	return sync, nil
}

func (s *SyncServer) SetSyncAPI(a api.SyncAPI) {
	s.syncAPI = a
}

func (s *SyncServer) InitMsgRouter() {
	for _, v := range s.links {
		v.AddMsgRouter(string(specV1.MessageReport), HandlerReport(s.syncAPI.Report))
		v.AddMsgRouter(string(specV1.MessageDesire), HandlerDesire(s.syncAPI.Desire))
	}
}

func (s *SyncServer) Run() {
	for k, v := range s.links {
		go v.Start()
		log.L().Info("sync server starting", log.Any("name", k))
	}
}

func (s *SyncServer) Close() {
	for _, v := range s.links {
		v.Close()
	}
}
