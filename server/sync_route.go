package server

import (
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/api"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

type HandlerReport func(msg specV1.Message) (*specV1.Message, error)
type HandlerDesire func(msg specV1.Message) (*specV1.Message, error)

const (
	MethodNameReport = "report"
	MethodNameDesire = "desire"
)

type MsgRouter struct {
	SyncAPI api.SyncAPI
	Link    plugin.SyncLink
}

func (m *MsgRouter) InitMsgRouter() {
	m.Link.AddMsgRouter(MethodNameReport, HandlerReport(m.SyncAPI.Report))
	m.Link.AddMsgRouter(MethodNameDesire, HandlerDesire(m.SyncAPI.Desire))
}
