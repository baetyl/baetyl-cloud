package common

import "github.com/baetyl/baetyl-go/log"

func LogDirtyData(err error, fields ...log.Field) {
	fields = append(fields, log.Error(err))
	log.L().Warn("dirty data", fields...)
}
