package plugin

import (
	"github.com/baetyl/baetyl-go/v2/pubsub"
)

//go:generate mockgen -destination=../mock/plugin/pubsub.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Pubsub

type Pubsub interface {
	pubsub.Pubsub
}
