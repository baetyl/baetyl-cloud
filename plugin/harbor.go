package plugin

import (
	"io"
)

//go:generate mockgen -destination=../mock/plugin/harbor.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Harbor

type Harbor interface {
	GetImageDigest(projects, repo, tags string) (string, error)
	ListProjects() ([]string, error)
	io.Closer
}
