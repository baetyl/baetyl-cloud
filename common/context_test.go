package common

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWrapper(t *testing.T) {
	test200 := func(c *Context) (interface{}, error) {
		return nil, nil
	}
	test404 := func(c *Context) (interface{}, error) {
		return nil, Error(ErrResourceNotFound, Field("name", "test"))
	}

	test401 := func(c *Context) (interface{}, error) {
		return nil, Error(ErrRequestAccessDenied)
	}
	test400 := func(c *Context) (interface{}, error) {
		return nil, Error(ErrRequestParamInvalid)
	}

	testPanic := func(c *Context) (interface{}, error) {
		panic("panic test")
	}
	router := gin.Default()
	router.GET("/200", Wrapper(test200))
	router.GET("/404", Wrapper(test404))
	router.GET("/400", Wrapper(test400))
	router.GET("/401", Wrapper(test401))
	router.GET("/panic", Wrapper(testPanic))

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/200", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/404", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	// 500
	req, _ = http.NewRequest(http.MethodGet, "/panic", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)

	// 401
	req, _ = http.NewRequest(http.MethodGet, "/401", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req)
	assert.Equal(t, http.StatusUnauthorized, w4.Code)

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/400", nil)
	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req)
	assert.Equal(t, http.StatusBadRequest, w5.Code)
}

func TestWrapperRaw(t *testing.T) {
	test200 := func(c *Context) (interface{}, error) {
		return nil, nil
	}
	test200Byte := func(c *Context) (interface{}, error) {
		return []byte("200"), nil
	}
	test500 := func(c *Context) (interface{}, error) {
		return "other", nil
	}
	test404 := func(c *Context) (interface{}, error) {
		return nil, Error(ErrResourceNotFound, Field("name", "test"))
	}

	test401 := func(c *Context) (interface{}, error) {
		return nil, Error(ErrRequestAccessDenied)
	}
	test400 := func(c *Context) (interface{}, error) {
		return nil, Error(ErrRequestParamInvalid)
	}

	testPanic := func(c *Context) (interface{}, error) {
		panic("panic test")
	}

	router := gin.Default()
	router.GET("/200", WrapperRaw(test200))
	router.GET("/200Byte", WrapperRaw(test200Byte))
	router.GET("/500", WrapperRaw(test500))
	router.GET("/404", WrapperRaw(test404))
	router.GET("/400", WrapperRaw(test400))
	router.GET("/401", WrapperRaw(test401))
	router.GET("/panic", WrapperRaw(testPanic))

	// 200
	req, _ := http.NewRequest(http.MethodGet, "/200", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 201Byte
	req, _ = http.NewRequest(http.MethodGet, "/200Byte", nil)
	w1 = httptest.NewRecorder()
	router.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 500
	req, _ = http.NewRequest(http.MethodGet, "/500", nil)
	w1 = httptest.NewRecorder()
	router.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusInternalServerError, w1.Code)

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/404", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	// 500
	req, _ = http.NewRequest(http.MethodGet, "/panic", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)

	// 401
	req, _ = http.NewRequest(http.MethodGet, "/401", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req)
	assert.Equal(t, http.StatusUnauthorized, w4.Code)

	// 400
	req, _ = http.NewRequest(http.MethodGet, "/400", nil)
	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req)
	assert.Equal(t, http.StatusBadRequest, w5.Code)
}

func TestValid(t *testing.T) {
	test1 := &struct {
		Mem string `validate:"memory"`
	}{Mem: "10g"}
	err := validate.Struct(test1)
	assert.NoError(t, err)
	test2 := &struct {
		Mem string `validate:"memory"`
	}{Mem: "tt"}
	err = validate.Struct(test2)
	assert.NotNil(t, err)
}

func TestContext_LoadBody(t *testing.T) {
	var model struct {
		Name string `json:"name" validate:"nonBaetyl,resourceName"`
	}
	gCtx := &gin.Context{
		Request: &http.Request{
			Body: newStringReaderColser(`{"name":"baetyl-test"}`),
		},
	}
	ctx := NewContext(gCtx)
	err := ctx.LoadBody(&model)
	assert.Error(t, err)
}

type stringReaderCloser struct {
	reader *strings.Reader
	io.Closer
}

func newStringReaderColser(str string) *stringReaderCloser {
	return &stringReaderCloser{reader: strings.NewReader(str)}
}

func (s *stringReaderCloser) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s *stringReaderCloser) Close() error {
	return nil
}
