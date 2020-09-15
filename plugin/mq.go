package plugin

import (
	"github.com/baetyl/baetyl-go/v2/mq"
)

//go:generate mockgen -destination=../mock/plugin/mq.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin MessageQueue

type MessageQueue interface {
	mq.MessageQueue
}
