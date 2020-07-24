package api

import (
	"bytes"
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/common"
	plugin "github.com/baetyl/baetyl-cloud/mock/service"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initPropertyAPI(t *testing.T) (*API, *gin.Engine, *gomock.Controller) {
	api := &API{}
	router := gin.Default()
	mockCtl := gomock.NewController(t)
	mockIM := func(c *gin.Context) { common.NewContext(c).SetNamespace("default") }

	v1 := router.Group("v1")
	{
		property := v1.Group("/properties")

		property.GET("/:key", mockIM, common.Wrapper(api.GetProperty))
		property.GET("", mockIM, common.Wrapper(api.ListProperty))

		property.POST("", mockIM, common.Wrapper(api.CreateProperty))
		property.DELETE("/:key", mockIM, common.Wrapper(api.DeleteProperty))
		property.PUT("/:key", mockIM, common.Wrapper(api.UpdateProperty))
	}
	return api, router, mockCtl
}

func genProperty() *models.Property {
	return &models.Property{
		Key:   "bae",
		Value: "http://test",
	}
}

func TestCreateProperty(t *testing.T) {
	api, router, ctl := initPropertyAPI(t)
	rs := plugin.NewMockPropertyService(ctl)
	api.propertyService = rs

	property := genProperty()

	rs.EXPECT().CreateProperty(property).Return(nil, nil).Times(1)

	body, _ := json.Marshal(property)
	req, _ := http.NewRequest(http.MethodPost, "/v1/properties", bytes.NewReader(body))
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "baetyl-cloud-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteProperty(t *testing.T) {
	api, router, ctl := initPropertyAPI(t)
	rs := plugin.NewMockPropertyService(ctl)
	api.propertyService = rs

	property := genProperty()

	rs.EXPECT().DeleteProperty(gomock.Any()).Return(nil).Times(1)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/properties/"+property.Key, nil)
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "baetyl-cloud-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProperty(t *testing.T) {
	api, router, ctl := initPropertyAPI(t)
	rs := plugin.NewMockPropertyService(ctl)
	api.propertyService = rs

	property := genProperty()

	rs.EXPECT().GetProperty(property.Key).Return(property, nil).Times(2)

	req, _ := http.NewRequest(http.MethodGet, "/v1/properties/"+property.Key, nil)
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "baetyl-cloud-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListProperty(t *testing.T) {
	api, router, ctl := initPropertyAPI(t)
	rs := plugin.NewMockPropertyService(ctl)
	api.propertyService = rs

	mConf := genProperty()
	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
		Name:     "%",
	}
	rs.EXPECT().ListProperty(page).Return([]models.Property{*mConf}, 1, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/properties?pageNo=1&pageSize=2", nil)
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "baetyl-cloud-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateProperty(t *testing.T) {
	api, router, ctl := initPropertyAPI(t)
	rs := plugin.NewMockPropertyService(ctl)
	api.propertyService = rs

	property := genProperty()

	rs.EXPECT().GetProperty(property.Key).Return(nil, nil).Times(1)
	rs.EXPECT().UpdateProperty(property).Return(nil, nil).Times(1)

	body, _ := json.Marshal(property)
	req, _ := http.NewRequest(http.MethodPut, "/v1/properties/"+property.Key, bytes.NewReader(body))
	req.Header.Set("baetyl-cloud-token", "baetyl-cloud-token")
	req.Header.Set("baetyl-cloud-user", "baetyl-cloud-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
