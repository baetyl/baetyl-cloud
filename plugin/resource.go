package plugin

import "io"

//go:generate mockgen -destination=../mock/plugin/resource.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Resource

type Resource interface {
	Node
	Application
	Configuration
	Secret
	Namespace
	io.Closer
}
