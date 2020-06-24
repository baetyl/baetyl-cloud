package api

import (
	"bytes"
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/common"
	"github.com/baetyl/baetyl-cloud/mock/service"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initCallBackAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }
	v1 := router.Group("v1")
	{
		callback := v1.Group("/callback")
		callback.POST("", mockIM, common.Wrapper(api.CreateCallback))
		callback.PUT("/:callbackName", mockIM, common.Wrapper(api.UpdateCallback))
		callback.DELETE("/:callbackName", mockIM, common.Wrapper(api.DeleteCallback))
		callback.GET("/:callbackName", mockIM, common.Wrapper(api.GetCallback))

	}
	return api, router, mockCtl
}

func genCallback() *models.Callback {
	return &models.Callback{
		Name:        "test",
		Namespace:   "default",
		Method:      http.MethodPost,
		Url:         "http://www.baidu.com",
		Params:      map[string]string{"a": "a"},
		Header:      map[string]string{"b": "b"},
		Body:        map[string]string{"c": "c"},
		Description: "hhhhhh",
	}
}

func TestAPI_CreateCallback(t *testing.T) {
	api, router, ctl := initCallBackAPI(t)
	rs := plugin.NewMockCallbackService(ctl)
	api.callbackService = rs

	callback := genCallback()

	rs.EXPECT().Create(gomock.Any()).Return(callback, nil).AnyTimes()

	body, _ := json.Marshal(callback)
	req, _ := http.NewRequest(http.MethodPost, "/v1/callback", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetCallback(t *testing.T) {
	api, router, ctl := initCallBackAPI(t)
	rs := plugin.NewMockCallbackService(ctl)
	api.callbackService = rs

	callback := genCallback()

	rs.EXPECT().Get(callback.Name, callback.Namespace).Return(callback, nil).AnyTimes()
	req, _ := http.NewRequest(http.MethodGet, "/v1/callback/"+callback.Name, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_GetCallback_Err(t *testing.T) {
	api, router, ctl := initCallBackAPI(t)
	rs := plugin.NewMockCallbackService(ctl)
	api.callbackService = rs

	rs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	req, _ := http.NewRequest(http.MethodGet, "/v1/callback/error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_UpdateCallback(t *testing.T) {
	api, router, ctl := initCallBackAPI(t)
	rs := plugin.NewMockCallbackService(ctl)
	api.callbackService = rs

	callback := genCallback()
	rs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(callback, nil).AnyTimes()
	rs.EXPECT().Update(gomock.Any()).Return(callback, nil).AnyTimes()

	callback.Params = nil
	callback.Body = nil
	callback.Header = nil
	body, _ := json.Marshal(callback)
	req, _ := http.NewRequest(http.MethodPut, "/v1/callback/"+callback.Name, bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPI_UpdateCallback_Err(t *testing.T) {
	api, router, ctl := initCallBackAPI(t)
	rs := plugin.NewMockCallbackService(ctl)
	api.callbackService = rs

	callback := genCallback()
	rs.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	rs.EXPECT().Update(gomock.Any()).Return(callback, nil).AnyTimes()

	callback.Params = nil
	callback.Body = nil
	callback.Header = nil
	body, _ := json.Marshal(callback)
	req, _ := http.NewRequest(http.MethodPut, "/v1/callback/"+callback.Name, bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPI_DeleteCallback(t *testing.T) {
	api, router, ctl := initCallBackAPI(t)
	rs := plugin.NewMockCallbackService(ctl)
	api.callbackService = rs

	callback := genCallback()
	rs.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	body, _ := json.Marshal(callback)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/callback/"+callback.Name, bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
