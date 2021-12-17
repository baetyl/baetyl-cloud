package api

import (
	"encoding/json"
	"strings"

	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
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
	Node service.NodeService
	log  *log.Logger
}

func NewSyncAPI(cfg *config.CloudConfig) (SyncAPI, error) {
	syncService, err := service.NewSyncService(cfg)
	if err != nil {
		return nil, err
	}
	nodeService, err := service.NewNodeService(cfg)
	if err != nil {
		return nil, err
	}
	return &SyncAPIImpl{
		Sync: syncService,
		Node: nodeService,
		log:  log.L().With(log.Any("api", "sync")),
	}, nil
}

// Report for node report
func (s *SyncAPIImpl) Report(msg specV1.Message) (*specV1.Message, error) {
	var report specV1.Report
	err := msg.Content.Unmarshal(&report)
	if err != nil {
		return nil, err
	}

	setNodeClientIPIfExist(msg, &report)

	// TODO remove the trick. set node prop if source=baetyl-init
	ns, n := msg.Metadata["namespace"], msg.Metadata["name"]
	if msg.Metadata != nil && msg.Metadata["source"] == specV1.BaetylInit {
		nodeInfo, err := s.Node.Get(nil, ns, n)
		if err != nil {
			s.log.Warn("failed to get node info", log.Any("source", specV1.BaetylInit))
		} else {
			s.log.Debug("set init node info", log.Any("source", specV1.BaetylInit))
			report = updateCoreStateOnly(report, nodeInfo)
		}
	}
	delta, err := s.Sync.Report(ns, n, report)
	if err != nil {
		return nil, err
	}

	s.log.Debug("api sync", log.Any("delta", delta), log.Any("report", report))

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

func setNodeClientIPIfExist(msg specV1.Message, report *specV1.Report) {
	if ip, ok := msg.Metadata["clientIP"]; !ok {
		return
	} else {
		nodeVal, ok := (*report)[common.NodeInfo]
		if ok {
			nodes, ok := nodeVal.(map[string]interface{})
			if !ok {
				return
			}
			for k, info := range nodes {
				node, ok := info.(map[string]interface{})
				if !ok {
					continue
				}
				node["clientIP"] = ip
				(*report)[common.NodeInfo].(map[string]interface{})[k] = node
			}
		}
	}
}

func updateCoreStateOnly(report specV1.Report, nodeInfo *specV1.Node) specV1.Report {
	coreApps, err := json.Marshal(nodeInfo.Report["sysapps"])
	if err != nil {
		return nodeInfo.Report
	}
	initApps, err := json.Marshal(report["sysapps"])
	if err != nil {
		return nodeInfo.Report
	}
	var coreSysApps, initSysApps []specV1.AppInfo
	err = json.Unmarshal(coreApps, &coreSysApps)
	if err != nil {
		return nodeInfo.Report
	}
	err = json.Unmarshal(initApps, &initSysApps)
	if err != nil {
		return nodeInfo.Report
	}

	mApps := map[string]specV1.AppInfo{}
	for _, v := range coreSysApps {
		mApps[v.Name] = v
	}
	for _, v := range initSysApps {
		mApps[v.Name] = v
	}
	sysApps := []specV1.AppInfo{}
	for _, v := range mApps {
		sysApps = append(sysApps, v)
	}
	nodeInfo.Report["sysapps"] = sysApps

	coreState, err := json.Marshal(nodeInfo.Report["sysappstats"])
	if err != nil {
		return nodeInfo.Report
	}
	initState, err := json.Marshal(report["sysappstats"])
	if err != nil {
		return nodeInfo.Report
	}
	var coreSysAppState, initSysAppState []specV1.AppStats
	err = json.Unmarshal(coreState, &coreSysAppState)
	if err != nil {
		return nodeInfo.Report
	}
	err = json.Unmarshal(initState, &initSysAppState)
	if err != nil {
		return nodeInfo.Report
	}
	mState := map[string]specV1.AppStats{}
	for _, v := range coreSysAppState {
		mState[v.Name] = v
	}
	for _, v := range initSysAppState {
		if strings.Contains(v.Name, specV1.BaetylCore) {
			mState[v.Name] = v
		}
	}
	sysState := []specV1.AppStats{}
	for _, v := range mState {
		sysState = append(sysState, v)
	}
	nodeInfo.Report["sysappstats"] = sysState

	return nodeInfo.Report
}
