package httplink

import (
	"strings"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/server"
)

func (l *httpLink) wrapper(tp specV1.MessageKind) common.HandlerFunc {
	switch tp {
	case specV1.MessageReport:
		return func(c *common.Context) (interface{}, error) {
			ns, n := c.GetNamespace(), c.GetName()
			if ns == "" || n == "" {
				return nil, common.Error(common.ErrRequestParamInvalid)
			}
			body, err := c.GetRawData()
			if err != nil {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
			}

			msg := specV1.Message{
				Kind:     specV1.MessageReport,
				Content:  specV1.LazyValue{},
				Metadata: map[string]string{},
			}
			err = msg.Content.UnmarshalJSON(body)
			if err != nil {
				return nil, err
			}
			for k := range c.Request.Header {
				msg.Metadata[strings.ToLower(k)] = c.GetHeader(k)
			}
			msg.Metadata["name"] = n
			msg.Metadata["namespace"] = ns
			msg.Metadata["address"] = c.ClientIP()
			resp, err := l.msgRouter[string(specV1.MessageReport)].(server.HandlerMessage)(msg)
			if err != nil {
				return nil, err
			}
			return resp.Content.Value, nil
		}
	case specV1.MessageDesire:
		return func(c *common.Context) (interface{}, error) {
			ns := c.GetNamespace()
			if ns == "" {
				return nil, common.Error(common.ErrRequestParamInvalid)
			}

			body, err := c.GetRawData()
			if err != nil {
				return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
			}

			msg := specV1.Message{
				Kind:     specV1.MessageDesire,
				Content:  specV1.LazyValue{},
				Metadata: map[string]string{},
			}
			err = msg.Content.UnmarshalJSON(body)
			if err != nil {
				return nil, err
			}

			for k := range c.Request.Header {
				msg.Metadata[strings.ToLower(k)] = c.GetHeader(k)
			}
			msg.Metadata["namespace"] = ns
			resp, err := l.msgRouter[string(specV1.MessageDesire)].(server.HandlerMessage)(msg)
			if err != nil {
				return nil, err
			}
			return resp.Content.Value, nil
		}
	}
	return func(c *common.Context) (interface{}, error) {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "messageType"))
	}
}
