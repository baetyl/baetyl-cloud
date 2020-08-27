package common

import (
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
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
	arrayLength    = 3
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
var trans ut.Translator

func init() {
	labelRegex, _ = regexp.Compile("^([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]?$")
	resourceRegex, _ = regexp.Compile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$")
	fingerprintRegex, _ = regexp.Compile("^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?(\\.[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?)*$")
	validate = validator.New()
	validate.RegisterValidation(nonBaetyl, nonBaetylFunc())
	validate.RegisterValidation(validLabels, validLabelsFunc())
	validate.RegisterValidation(resourceName, validRexAndLengthFunc(resourceLength, resourceRegex))
	validate.RegisterValidation(fingerprintValue, validRexAndLengthFunc(resourceLength, fingerprintRegex))
	for k, v := range regexps {
		validate.RegisterValidation(k, genValidFunc(v))
	}
	validate.RegisterValidation(validBatchOp, validArrayLengthFunc(arrayLength))
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("en")
	validate.RegisterTranslation(validBatchOp, trans, registerTranslator(validBatchOp, ErrBatchOpNum.String()), translate)
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

func validArrayLengthFunc(length int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface().([]string)
		if len(field) > length {
			return false
		}
		return true
	}
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}
