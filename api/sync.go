package api

import (
	"encoding/json"
	"strings"

	"github.com/baetyl/baetyl-go/v2/context"
	"github.com/baetyl/baetyl-go/v2/errors"
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
	if msg.Metadata != nil {
		switch msg.Metadata["source"] {
		case specV1.BaetylInit:
			props, err := s.Node.GetNodeProperties(ns, n)
			if err != nil {
				s.log.Warn("failed to get node properties", log.Any("source", specV1.BaetylInit))
			} else {
				s.log.Debug("set init node properties", log.Any("source", specV1.BaetylInit))
				if props != nil {
					report[common.NodeProps] = props.State.Report
				}
			}
		case specV1.BaetylCore:
			nodeInfo, err := s.Node.Get(nil, ns, n)
			if err != nil {
				s.log.Warn("failed to get node properties", log.Any("source", specV1.BaetylInit))
			} else {
				report = keepCoreStateOnly(report, nodeInfo)
			}
		case specV1.BaetylCoreAndroid:
			nodeInfo, err := s.Node.Get(nil, ns, n)
			if err != nil {
				s.log.Warn("failed to get node", log.Any("source", specV1.BaetylCoreAndroid))
			} else {
				err = s.updateAndroidInfo(nodeInfo, &report)
				if err != nil {
					s.log.Warn("failed to update node android info", log.Any("source", specV1.BaetylCoreAndroid))
				}
			}
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

func (s *SyncAPIImpl) updateAndroidInfo(node *specV1.Node, report *specV1.Report) error {
	nodeVal, ok := (*report)[common.NodeInfo]
	if !ok {
		return nil
	}
	nodes, ok := nodeVal.(map[string]interface{})
	if !ok {
		return nil
	}
	deviceId := ""
	// android only ONE key
	for k, _ := range nodes {
		deviceId = k
		break
	}
	if deviceId == "" {
		return nil
	}
	if node.Attributes == nil {
		node.Attributes = map[string]interface{}{}
	}
	if node.Attributes[context.RunModeAndroid] != deviceId {
		node.Attributes[context.RunModeAndroid] = deviceId
		_, err := s.Node.Update(node.Namespace, node)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
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

func keepCoreStateOnly(report specV1.Report, nodeInfo *specV1.Node) specV1.Report {
	oldApps, err := json.Marshal(nodeInfo.Report["sysapps"])
	if err != nil {
		return report
	}
	newApps, err := json.Marshal(report["sysapps"])
	if err != nil {
		return report
	}
	var oldSysApps, newSysApps []specV1.AppInfo
	err = json.Unmarshal(oldApps, &oldSysApps)
	if err != nil {
		return report
	}
	err = json.Unmarshal(newApps, &newSysApps)
	if err != nil {
		return report
	}

	coreFlag := true
	for _, v := range newSysApps {
		if strings.Contains(v.Name, specV1.BaetylCore) {
			coreFlag = false
			break
		}
	}
	if coreFlag {
		for _, v := range oldSysApps {
			if strings.Contains(v.Name, specV1.BaetylCore) {
				newSysApps = append(newSysApps, v)
				break
			}
		}
	}
	report["sysapps"] = newSysApps

	oldState, err := json.Marshal(nodeInfo.Report["sysappstats"])
	if err != nil {
		return report
	}
	newState, err := json.Marshal(report["sysappstats"])
	if err != nil {
		return report
	}
	var oldSysAppState, newSysAppState []specV1.AppStats
	err = json.Unmarshal(oldState, &oldSysAppState)
	if err != nil {
		return report
	}
	err = json.Unmarshal(newState, &newSysAppState)
	if err != nil {
		return report
	}
	coreStatFlag := true
	for _, v := range newSysAppState {
		if strings.Contains(v.Name, specV1.BaetylCore) {
			coreStatFlag = false
			break
		}
	}
	if coreStatFlag {
		for _, v := range oldSysAppState {
			if strings.Contains(v.Name, specV1.BaetylCore) {
				newSysAppState = append(newSysAppState, v)
				break
			}
		}
	}
	report["sysappstats"] = newSysAppState

	return report
}
