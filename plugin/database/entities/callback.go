package entities

import (
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/models"
	"github.com/baetyl/baetyl-go/log"
	"time"
)

type Callback struct {
	Name        string    `db:"name"`
	Namespace   string    `db:"namespace"`
	Method      string    `db:"method"`
	Url         string    `db:"url"`
	Params      string    `db:"params"`
	Header      string    `db:"header"`
	Body        string    `db:"body"`
	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func ToCallbackModel(c *Callback) *models.Callback {
	var params map[string]string
	if err := json.Unmarshal([]byte(c.Params), &params); err != nil {
		log.L().Error("callback db params unmarshal error", log.Any("params", c.Params))
	}
	var header map[string]string
	if err := json.Unmarshal([]byte(c.Header), &header); err != nil {
		log.L().Error("callback db header unmarshal error", log.Any("header", c.Header))
	}
	var body map[string]string
	if err := json.Unmarshal([]byte(c.Body), &body); err != nil {
		log.L().Error("callback db body unmarshal error", log.Any("body", c.Body))
	}
	res := &models.Callback{
		Name:        c.Name,
		Namespace:   c.Namespace,
		Method:      c.Method,
		Url:         c.Url,
		Description: c.Description,
		CreateTime:  c.CreateTime,
		UpdateTime:  c.UpdateTime,
		Params:      params,
		Header:      header,
		Body:        body,
	}
	return res
}

func FromCallbackModel(c *models.Callback) *Callback {
	params, err := json.Marshal(c.Params)
	if err != nil {
		log.L().Error("callback params marshal error", log.Any("params", c.Params))
		params = []byte("{}")
	}
	header, err := json.Marshal(c.Header)
	if err != nil {
		log.L().Error("callback header marshal error", log.Any("header", c.Header))
		header = []byte("{}")
	}
	body, err := json.Marshal(c.Body)
	if err != nil {
		log.L().Error("callback body marshal error", log.Any("body", c.Body))
		body = []byte("{}")
	}
	res := &Callback{
		Name:        c.Name,
		Namespace:   c.Namespace,
		Method:      c.Method,
		Url:         c.Url,
		Description: c.Description,
		CreateTime:  c.CreateTime,
		UpdateTime:  c.UpdateTime,
		Params:      string(params),
		Header:      string(header),
		Body:        string(body),
	}
	return res
}
