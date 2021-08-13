package plugin

import (
	"io"
	"time"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

//go:generate mockgen -destination=../mock/plugin/jwt.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin JWT

type JWT interface {
	GetJWT(c *common.Context) (string, error)
	GenerateJWT(c *common.Context) (*JWTInfo, error)
	RefreshJWT(c *common.Context) (*JWTInfo, error)
	CheckAndParseJWT(c *common.Context) (map[string]interface{}, error)
	io.Closer
}

type JWTInfo struct {
	Token      string    `json:"token"`
	Expire     time.Time `json:"expire"`
	MaxRefresh time.Time `json:"maxRefresh"`
}
