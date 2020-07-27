package common

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/log"
	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

// Context context
type Context struct {
	*gin.Context
}

type User struct {
	ID string
}

// NewContext create a new context with gin context
func NewContext(inner *gin.Context) *Context {
	return &Context{inner}
}

// NewContextEmpty create a new empty context
func NewContextEmpty() *Context {
	return &Context{&gin.Context{}}
}

// SetNamespace sets namespace into context
func (c *Context) SetNamespace(ns string) {
	c.Set("namespace", ns)
}

// GetNamespace gets namespace from context if exists
func (c *Context) GetNamespace() string {
	return c.GetString("namespace")
}

// SetUser sets user into context
func (c *Context) SetUser(user User) {
	c.Set("user", user)
}

// GetUser gets user from context if exists
func (c *Context) GetUser() User {
	user, ok := c.Get("user")
	if !ok {
		return User{}
	}
	return user.(User)
}

// SetName sets name into context
func (c *Context) SetName(n string) {
	c.Set("name", n)
}

// GetName gets name from context if exists
func (c *Context) GetName() string {
	return c.GetString("name")
}

// GetNameFromParam gets name from param if exists
func (c *Context) GetNameFromParam() string {
	return c.Param("name")
}

// SetTrace set the trace key and value
func (c *Context) SetTrace() {
	k := GetTraceHeader()
	v := c.Request.Header.Get(k)
	if v == "" {
		v = uuid.NewV4().String()
	}
	c.Writer.Header().Set(k, v)
}

// GetTrace gets the trace key and value
func (c *Context) GetTrace() (k string, v string) {
	return GetTraceKey(), c.Request.Header.Get(GetTraceHeader())
}

// LoadBody loads json data from body into object and set defaults
func (c *Context) LoadBody(obj interface{}) error {
	err := c.BindJSON(obj)
	if err != nil {
		return err
	}
	err = validate.Struct(obj)
	if err != nil {
		if es, ok := err.(validator.ValidationErrors); ok {
			for _, v := range es {
				return Error(Code(v.Tag()), Field(v.Tag(), v.Field()), Field("error", err.Error()))
			}
		}
		return err
	}
	return utils.SetDefaults(obj)
}

type sucResponse struct {
	Success bool `json:"success"`
}

// PackageResponse PackageResponse
func PackageResponse(res interface{}) (int, interface{}) {
	if res == nil {
		res = &sucResponse{
			Success: true,
		}
	}
	return http.StatusOK, res
}

// PopulateFailedResponse PopulateFailedResponse
func PopulateFailedResponse(cc *Context, err error, abort bool) {
	var code string
	var status int
	switch e := err.(type) {
	case errors.Coder:
		code = e.Code()
		status = getHTTPStatus(Code(e.Code()))
	default:
		code = ErrUnknown
		status = http.StatusInternalServerError
	}

	log.L().Error("process failed.", log.Any(cc.GetTrace()), log.Code(err))

	k, v := cc.GetTrace()
	body := gin.H{
		"code":    code,
		"message": err.Error(),
		k:         v,
	}
	if abort {
		cc.AbortWithStatusJSON(status, body)
	} else {
		cc.JSON(status, body)
	}
}

// HandlerFunc HandlerFunc
type HandlerFunc func(c *Context) (interface{}, error)

// Wrapper Wrapper
// TODO: to use gin.HandlerFunc ?
func Wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := NewContext(c)
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = Error(ErrUnknown, Field("error", r))
				}
				log.L().Info("handle a panic", log.Any(cc.GetTrace()), log.Code(err), log.Error(err), log.Any("panic", string(debug.Stack())))
				PopulateFailedResponse(cc, err, false)
			}
		}()
		res, err := handler(cc)
		if err != nil {
			log.L().Error("failed to handler request", log.Any(cc.GetTrace()), log.Code(err), log.Error(err))
			PopulateFailedResponse(cc, err, false)
			return
		}
		log.L().Debug("process success", log.Any(cc.GetTrace()), log.Any("response", _toJsonString(res)))
		// unlike JSON, does not replace special html characters with their unicode entities. eg: JSON(&)->'\u0026' PureJSON(&)->'&'
		cc.PureJSON(PackageResponse(res))
	}
}

func WrapperRaw(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := NewContext(c)
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = Error(ErrUnknown, Field("error", r))
				}
				log.L().Info("handle a panic", log.Any(cc.GetTrace()), log.Code(err), log.Error(err))
				PopulateFailedResponse(cc, err, false)
			}
		}()
		res, err := handler(cc)
		if err != nil {
			log.L().Error("failed to handler request", log.Any(cc.GetTrace()), log.Code(err), log.Error(err))
			PopulateFailedResponse(cc, err, false)
			return
		}
		if res == nil {
			return
		}
		if data, ok := res.([]byte); ok {
			cc.Data(http.StatusOK, "application/octet-stream", data)
		} else {
			log.L().Error("failed to convert data to []byte", log.Any(cc.GetTrace()))
			PopulateFailedResponse(cc, Error(ErrUnknown), false)
		}
	}
}

func _toJsonString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

// MisHandlerFunc MisHandlerFunc
type MisHandlerFunc func(c *Context) error

func WrapperMis(handler MisHandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		cc := NewContext(c)
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = Error(ErrUnknown, Field("error", r))
				}
				log.L().Info("handle a panic", log.Any(cc.GetTrace()), log.Code(err), log.Error(err), log.Any("panic", string(debug.Stack())))
				PopulateFailedMisResponse(cc, err, false)
			}
		}()
		err := handler(cc)
		if err != nil {
			log.L().Error("failed to handler request", log.Any(cc.GetTrace()), log.Code(err), log.Error(err))
			PopulateFailedMisResponse(cc, err, false)
			return
		}
		log.L().Debug("process success", log.Any(cc.GetTrace()))
		if cc.Request.Method != "GET" {
			// unlike JSON, does not replace special html characters with their unicode entities. eg: JSON(&)->'\u0026' PureJSON(&)->'&'
			cc.PureJSON(PackageMisResponse(nil))
		}
	}
}

// PopulateFailedMisResponse PopulateFailedMisResponse
func PopulateFailedMisResponse(cc *Context, err error, abort bool) {
	var status int = http.StatusOK
	log.L().Error("process failed.", log.Any(cc.GetTrace()), log.Code(err))

	body := gin.H{
		"status": 1,
		"msg":    err.Error(),
	}
	if abort {
		cc.AbortWithStatusJSON(status, body)
	} else {
		cc.JSON(status, body)
	}
}

func PackageMisResponse(res interface{}) (int, interface{}) {
	if res == nil {
		res = "[]"
	}
	return http.StatusOK, gin.H{
		"status": 0,
		"msg":    "ok",
		"data":   res,
	}
}
