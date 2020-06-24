package server

import (
	"bytes"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-go/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderCommonName = "common-name"
)

func noRouteHandler(c *gin.Context) {
	common.PopulateFailedResponse(common.NewContext(c), common.Error(common.ErrRequestMethodNotFound), true)
}

func noMethodHandler(c *gin.Context) {
	common.PopulateFailedResponse(common.NewContext(c), common.Error(common.ErrRequestMethodNotFound), true)
}

func requestIDHandler(c *gin.Context) {
	cc := common.NewContext(c)
	cc.SetTrace()
	cc.Next()
}

func loggerHandler(c *gin.Context) {
	cc := common.NewContext(c)
	log.L().Info("start request",
		log.Any(cc.GetTrace()),
		log.Any("method", cc.Request.Method),
		log.Any("url", cc.Request.URL.Path),
		log.Any("clientip", cc.ClientIP()),
	)
	if c.Request.Body != nil {
		if buf, err := ioutil.ReadAll(c.Request.Body); err == nil {
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf[:]))
			log.L().Info("request body",
				log.Any(cc.GetTrace()),
				log.Any("body", string(buf)),
			)
		}
	}
	start := time.Now()
	c.Next()
	log.L().Info("finish request",
		log.Any(cc.GetTrace()),
		log.Any("status", strconv.Itoa(c.Writer.Status())),
		log.Any("latency", time.Since(start)),
	)
}

func health(c *gin.Context) {
	c.JSON(common.PackageResponse(nil))
}

func extractNodeCommonNameFromCert(c *gin.Context) {
	cc := common.NewContext(c)
	if len(c.Request.TLS.PeerCertificates) == 0 {
		common.PopulateFailedResponse(cc, common.Error(common.ErrRequestAccessDenied), true)
		return
	}
	cert := c.Request.TLS.PeerCertificates[0]
	extractNodeCommonName(cc, cert.Subject.CommonName)
}

func extractNodeCommonNameFromHeader(c *gin.Context) {
	cc := common.NewContext(c)
	extractNodeCommonName(cc, c.GetHeader(HeaderCommonName))
}

func extractNodeCommonName(cc *common.Context, commonName string) {
	res := strings.SplitN(commonName, ".", 2)
	if len(res) != 2 || res[0] == "" || res[1] == "" {
		common.PopulateFailedResponse(cc, common.Error(common.ErrRequestAccessDenied), true)
		return
	}
	cc.SetNamespace(res[0])
	cc.SetName(res[1])
}
