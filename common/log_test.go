package common

import (
	"testing"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/pkg/errors"
)

func TestLogDirtyData(t *testing.T) {
	err := errors.New("custom")
	LogDirtyData(err, log.Any("name", "baetyl"))
}
