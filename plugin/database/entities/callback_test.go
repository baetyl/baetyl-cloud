package entities

import (
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func genCallback() *models.Callback {
	return &models.Callback{
		Name:        "zx",
		Namespace:   "default",
		Method:      "Post",
		Url:         "www.baidu.com",
		Params:      map[string]string{"1": "1"},
		Header:      map[string]string{"2": "2"},
		Body:        map[string]string{"3": "3"},
		Description: "desc",
		CreateTime:  time.Unix(1000, 10),
		UpdateTime:  time.Unix(1000, 10),
	}
}

func genCallbackDB() *Callback {
	return &Callback{
		Name:        "zx",
		Namespace:   "default",
		Method:      "Post",
		Url:         "www.baidu.com",
		Params:      "{\"1\":\"1\"}",
		Header:      "{\"2\":\"2\"}",
		Body:        "{\"3\":\"3\"}",
		Description: "desc",
		CreateTime:  time.Unix(1000, 10),
		UpdateTime:  time.Unix(1000, 10),
	}
}

func TestConvertCallback(t *testing.T) {
	callback := genCallback()
	callbackDB := genCallbackDB()
	resCallback := ToCallbackModel(callbackDB)
	assert.EqualValues(t, callback, resCallback)
	resCallbackDB := FromCallbackModel(callback)
	assert.EqualValues(t, callbackDB, resCallbackDB)
}
