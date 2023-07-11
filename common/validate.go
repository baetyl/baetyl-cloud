package common

import (
	"regexp"
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

func ValidateResourceName(s string) error {
	resourceRegex, _ := regexp.Compile("^[a-z0-9][-a-z0-9.]{0,61}[a-z0-9]$")
	if !resourceRegex.MatchString(s) {
		return Error(ErrRequestParamInvalid, Field("error", "resource name invalid"))
	}
	return nil
}

func ValidateKeyValue(k string) error {
	resourceRegex, _ := regexp.Compile("^[-._a-zA-Z0-9]+$")
	if !resourceRegex.MatchString(k) {
		return Error(ErrRequestParamInvalid, Field("error", "config data key value invalid"))
	}
	return nil
}
