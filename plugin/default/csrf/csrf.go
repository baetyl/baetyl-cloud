package csrf

import (
	"net/url"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-go/v2/log"
)

type defaultCsrf struct {
	cfg CloudConfig
	log *log.Logger
}

func init() {
	plugin.RegisterFactory("defaultcsrf", New)
}

func New() (plugin.Plugin, error) {
	var cfg CloudConfig
	if err := common.LoadConfig(&cfg); err != nil {
		return nil, err
	}
	return &defaultCsrf{
		cfg: cfg,
		log: log.With(log.Any("plugin", "defautlcsrf")),
	}, nil
}

func (d *defaultCsrf) Verify(c *common.Context) error {
	referer, err := url.Parse(c.Request.Referer())
	if err != nil {
		return common.Error(common.ErrRequestAccessDenied)
	}
	var inList bool
	for _, host := range d.cfg.CsrfConfig.Whitelist {
		if referer.Host == host {
			d.log.Debug("referer host in whitelist", log.Any("host", host))
			inList = true
			break
		}
	}
	if len(d.cfg.CsrfConfig.Whitelist) > 0 && !inList {
		d.log.Error("referer host not in whitelist",
			log.Any(c.GetTrace()),
			log.Any("referer host", referer.Host))
		return common.Error(common.ErrRequestAccessDenied)
	}

	csrfCookieValue, err := c.Cookie(d.cfg.CsrfConfig.CookieName)
	if err != nil || csrfCookieValue == "" {
		d.log.Error("fetch csrf cookie failed", log.Any(c.GetTrace()))
		return common.Error(common.ErrRequestAccessDenied)
	}
	csrfHeaderValue := c.GetHeader(d.cfg.CsrfConfig.HeaderName)
	if csrfCookieValue != csrfHeaderValue {
		d.log.Error("csrf cookie value not equal to header value",
			log.Any(c.GetTrace()),
			log.Any("header value", csrfHeaderValue),
			log.Any("cookie value", csrfCookieValue))
		return common.Error(common.ErrRequestAccessDenied)
	}
	return nil
}

// Close Close
func (d *defaultCsrf) Close() error {
	return nil
}
