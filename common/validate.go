package common

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	resourceName     = "resourceName"
	serviceName      = "serviceName"
	fingerprintValue = "fingerprintValue"
	memory           = "memory"
	duration         = "duration"
	setcpus          = "setcpus"
	devicemodel      = "devicemodel"
	nonBaetyl        = "nonBaetyl"
	namespace        = "namespace"
	validLabels      = "validLabels"
	validConfigKeys  = "validConfigKeys"
	maxLength        = "maxLength"
	resourceLength   = 63
)

var regexps = map[string]string{
	namespace:       "^[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$",
	memory:          "^[1-9][0-9]*(k|m|g|t|p|)$",
	duration:        "^[1-9][0-9]*(s|m|h)$",
	setcpus:         "^(([1-9]\\d*|0)-([1-9]\\d*|0)|([1-9]\\d*|0)(,([1-9]\\d*|0))*)$",
	validConfigKeys: "^[-._a-zA-Z0-9]+$",
	devicemodel:     "^[a-zA-Z0-9\\-_]{1,32}$",
}

var (
	validate         *validator.Validate
	labelRegex       *regexp.Regexp
	resourceRegex    *regexp.Regexp
	serviceRegex     *regexp.Regexp
	fingerprintRegex *regexp.Regexp
)

func init() {
	var ok bool
	validate, ok = binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.L().Error("failed to get binding validator")
		return
	}

	labelRegex, _ = regexp.Compile("^([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]?$")
	resourceRegex, _ = regexp.Compile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$")
	serviceRegex, _ = regexp.Compile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$")
	fingerprintRegex, _ = regexp.Compile("^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?(\\.[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?)*$")

	validate.RegisterValidation(nonBaetyl, nonBaetylFunc())
	validate.RegisterValidation(validLabels, validLabelsFunc())
	validate.RegisterValidation(resourceName, validRexAndLengthFunc(resourceLength, resourceRegex))
	validate.RegisterValidation(serviceName, validRexAndLengthFunc(resourceLength, serviceRegex))
	validate.RegisterValidation(fingerprintValue, validRexAndLengthFunc(resourceLength, fingerprintRegex))
	validate.RegisterValidation(maxLength, validMaxLengthFunc())
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
			if strings.Contains(k, "/") {
				ss := strings.Split(k, "/")
				if len(ss) != 2 {
					return false
				}
				if len(ss[0]) > 253 || len(ss[0]) < 1 || !labelRegex.MatchString(ss[0]) || len(ss[1]) > 63 || !labelRegex.MatchString(ss[1]) {
					return false
				}
			} else {
				if len(k) > 63 || !labelRegex.MatchString(k) {
					return false
				}
			}
			if len(v) > 63 || !labelRegex.MatchString(v) {
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

func validMaxLengthFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface().([]string)
		length, err := strconv.Atoi(fl.Param())
		if err != nil {
			return false
		}
		if len(field) > length {
			return false
		}
		return true
	}
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
