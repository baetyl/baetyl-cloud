package plugin

import (
	"io"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

type SyncLink interface {
	Run()
	Wrapper(tp specV1.MessageKind) common.HandlerFunc
	AddMsgRouter(k string, v interface{})
	io.Closer
}
