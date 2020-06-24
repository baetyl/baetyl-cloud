package common

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
	"strings"
)

const (
	resourceName     = "resourceName"
	fingerprintValue = "fingerprintValue"
	memory           = "memory"
	duration         = "duration"
	setcpus          = "setcpus"
	nonBaetyl        = "nonBaetyl"
	namespace        = "namespace"
	validLabels      = "validLabels"
)

var regexps = map[string]string{
	namespace:        "^[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$",
	resourceName:     "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$",
	fingerprintValue: "^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?(\\.[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?)*$",
	memory:           "^[1-9][0-9]*(k|m|g|t|p|)$",
	duration:         "^[1-9][0-9]*(s|m|h)$",
	setcpus:          "^(([1-9]\\d*|0)-([1-9]\\d*|0)|([1-9]\\d*|0)(,([1-9]\\d*|0))*)$",
}

var validate *validator.Validate
var labelRegex *regexp.Regexp

func init() {
	labelRegex, _ = regexp.Compile("^([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]?$")
	validate = validator.New()
	validate.RegisterValidation(nonBaetyl, nonBaetylFunc())
	validate.RegisterValidation(validLabels, validLabelsFunc())
	for k, v := range regexps {
		validate.RegisterValidation(k, genValidFunc(v))
	}
}

func genValidFunc(str string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(str, fl.Field().String())
		return match
	}
}

func nonBaetylFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return !strings.Contains(strings.ToLower(fl.Field().String()), "baetyl")
	}
}

func validLabelsFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		labels := fl.Field().Interface().(map[string]string)
		for k, v := range labels {
			if len(k) > 63 || len(v) > 63 || !labelRegex.MatchString(k) || !labelRegex.MatchString(v) {
				return false
			}
		}
		return true
	}
}
