package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/matcher.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Matcher

type Matcher interface {
	IsLabelMatch(labelSelector string, labels map[string]string) (bool, error)
	io.Closer
}
