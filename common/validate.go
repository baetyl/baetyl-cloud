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
	validBatchOp     = "validBatchOp"

	resourceLength = 63
	batchOpNum     = 20
)

var regexps = map[string]string{
	namespace: "^[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$",
	memory:    "^[1-9][0-9]*(k|m|g|t|p|)$",
	duration:  "^[1-9][0-9]*(s|m|h)$",
	setcpus:   "^(([1-9]\\d*|0)-([1-9]\\d*|0)|([1-9]\\d*|0)(,([1-9]\\d*|0))*)$",
}

var validate *validator.Validate
var labelRegex *regexp.Regexp
var resourceRegex *regexp.Regexp
var fingerprintRegex *regexp.Regexp

func init() {
	labelRegex, _ = regexp.Compile("^([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]?$")
	resourceRegex, _ = regexp.Compile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$")
	fingerprintRegex, _ = regexp.Compile("^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?(\\.[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?)*$")
	validate = validator.New()
	validate.RegisterValidation(nonBaetyl, nonBaetylFunc())
	validate.RegisterValidation(validLabels, validLabelsFunc())
	validate.RegisterValidation(resourceName, validRexAndLengthFunc(resourceLength, resourceRegex))
	validate.RegisterValidation(fingerprintValue, validRexAndLengthFunc(resourceLength, fingerprintRegex))
	validate.RegisterValidation(validBatchOp, validBatchOpNumFunc(batchOpNum))
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

func validRexAndLengthFunc(length int, reg *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface().(string)
		if len(field) > length || !reg.MatchString(field) {
			return false
		}
		return true
	}
}

func validBatchOpNumFunc(length int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface().([]string)
		if len(field) > length {
			return false
		}
		return true
	}
}
