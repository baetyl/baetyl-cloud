package common

import "github.com/baetyl/baetyl-go/v2/log"

func LogDirtyData(err error, fields ...log.Field) {
	fields = append(fields, log.Error(err))
	log.L().Error("dirty data", fields...)
}
