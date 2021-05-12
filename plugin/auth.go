package plugin

import (
	"io"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

//go:generate mockgen -destination=../mock/plugin/auth.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Auth

// Auth interfaces of auth
type Auth interface {
	Authenticate(c *common.Context) error
	AuthAndVerify(c *common.Context, pr *PermissionRequest) error
	Verify(c *common.Context, pr *PermissionRequest) error

	SignToken(meta []byte) ([]byte, error)
	VerifyToken(meta, sign []byte) bool

	io.Closer
}

const (
	PermissionRead = "READ"
	PermissionFull = "FULL_CONTROL"

	PermissionResourceConfig = "config"
	PermissionResourceSecret = "secret"
	PermissionResourceApp    = "app"
	PermissionResourceNode   = "node"
	PermissionResourceBatch  = "batch"
)

type PermissionRequest struct {
	Region         string         `json:"region"`
	Resource       string         `json:"resource"`
	Permission     []string       `json:"permission"`
	RequestContext RequestContext `json:"request_context"`
}

type RequestContext struct {
	IpAddress  string                 `json:"ip_address"`
	Referer    string                 `json:"referer"`
	Conditions map[string]interface{} `json:"conditions"`
}
