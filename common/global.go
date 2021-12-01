package common

import (
	"time"

	"github.com/gin-contrib/cache/persistence"
)

const (
	KeyConfFile      = "ConfFile"
	ValueConfFile    = "etc/baetyl/cloud.yml"
	KeyTraceKey      = "TraceKey"
	ValueTraceKey    = "requestId"
	KeyTraceHeader   = "TraceHeader"
	ValueTraceHeader = "x-bce-request-id" // TODO: change to x-baetyl-request-id when support configuration

	RegistryAuth     = "system-registry-auth"
)

var cache = persistence.NewInMemoryStore(time.Minute * 10)

func SetConfFile(v string) {
	cache.Set(KeyConfFile, v, -1)
}

func GetConfFile() string {
	res := ValueConfFile
	cache.Get(KeyConfFile, &res)
	return res
}

func SetTraceKey(v string) {
	cache.Set(KeyTraceKey, v, -1)
}

func GetTraceKey() string {
	res := ValueTraceKey
	cache.Get(KeyTraceKey, &res)
	return res
}

func SetTraceHeader(v string) {
	cache.Set(KeyTraceHeader, v, -1)
}

func GetTraceHeader() string {
	res := ValueTraceHeader
	cache.Get(KeyTraceHeader, &res)
	return res
}
