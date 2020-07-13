package common

import (
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/pkg/errors"
	"testing"
)

func TestLogDirtyData(t *testing.T) {
	err := errors.New("custom")
	LogDirtyData(err, log.Any("name", "baetyl"))
}
