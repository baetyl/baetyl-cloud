package common

import (
	"strconv"
	"strings"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	var ok bool
	validate, ok = binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.L().Error("failed to get binding validator")
		return
	}
	utils.RegisterValidate(validate)
}

func ValidNonBaetyl(name string) bool {
	return !strings.Contains(name, "baetyl")
}

func ValidIsInvisible(labels map[string]string) bool {
	v, ok := labels[ResourceInvisible]
	if !ok {
		return false
	}
	if res, _ := strconv.ParseBool(v); !res {
		return false
	}
	return true
}
