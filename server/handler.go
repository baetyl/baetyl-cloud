package server

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/gin-gonic/gin"

	"github.com/baetyl/baetyl-cloud/v2/common"
)

var (
	HeaderCommonName = "common-name"
)

func NoRouteHandler(c *gin.Context) {
	common.PopulateFailedResponse(common.NewContext(c), common.Error(common.ErrRequestMethodNotFound), true)
}

func NoMethodHandler(c *gin.Context) {
	common.PopulateFailedResponse(common.NewContext(c), common.Error(common.ErrRequestMethodNotFound), true)
}

func RequestIDHandler(c *gin.Context) {
	cc := common.NewContext(c)
	cc.SetTrace()
	cc.Next()
}

func LoggerHandler(c *gin.Context) {
	cc := common.NewContext(c)
	log.L().Info("logger handler start request",
		log.Any(cc.GetTrace()),
		log.Any("method", cc.Request.Method),
		log.Any("url", cc.Request.URL.Path),
		log.Any("host", cc.Request.Host),
		log.Any("header", cc.Request.Header),
		log.Any("clientip", cc.ClientIP()),
	)
	if c.Request.Header.Get("Content-type") == "application/json" && c.Request.Body != nil {
		if buf, err := io.ReadAll(c.Request.Body); err == nil {
			c.Request.Body = io.NopCloser(bytes.NewReader(buf[:]))
			log.L().Info("logger handler request body",
				log.Any(cc.GetTrace()),
				log.Any("body", string(buf)),
			)
		}
	}
	start := time.Now()
	c.Next()
	log.L().Info("logger handler finish request",
		log.Any(cc.GetTrace()),
		log.Any("status", strconv.Itoa(c.Writer.Status())),
		log.Any("latency", time.Since(start)),
		log.Any("size", c.Writer.Size()),
	)
}

func Health(c *gin.Context) {
	c.JSON(common.PackageResponse(nil))
}

func ExtractNodeCommonNameFromCert(c *gin.Context) {
	cc := common.NewContext(c)
	if len(c.Request.TLS.PeerCertificates) == 0 {
		common.PopulateFailedResponse(cc, common.Error(common.ErrRequestAccessDenied), true)
		return
	}
	cert := c.Request.TLS.PeerCertificates[0]
	extractNodeCommonName(cc, cert.Subject.CommonName)
}

func ExtractNodeCommonNameFromHeader(c *gin.Context) {
	cc := common.NewContext(c)
	extractNodeCommonName(cc, c.GetHeader(HeaderCommonName))
}

func extractNodeCommonName(cc *common.Context, commonName string) {
	res := strings.SplitN(commonName, ".", 2)
	if len(res) != 2 || res[0] == "" || res[1] == "" {
		log.L().Error("extract node common name error",
			log.Any(cc.GetTrace()),
			log.Any("commonName", commonName),
			log.Any("HeaderCommonName", HeaderCommonName))
		common.PopulateFailedResponse(cc, common.Error(common.ErrRequestAccessDenied), true)
		return
	}
	cc.SetNamespace(res[0])
	cc.SetName(res[1])
}
