package plugin

import (
	"io"
)

type SyncLink interface {
	Start()
	AddMsgRouter(k string, v interface{})
	io.Closer
}
