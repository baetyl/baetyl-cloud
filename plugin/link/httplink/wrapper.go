package httplink

import (
	"strings"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/server"
)

// wrapper 现在 http link 支持以下类型消息的处理
// 上报类：上报 report(默认)、设备消息上报 deviceReport、设备生命周期上报 thing.lifecycle.post
// 同步类：期望 desire(默认)、设备消息同步 deviceDesire
func (l *httpLink) wrapper(tp specV1.MessageKind) common.HandlerFunc {
	switch tp {
	case specV1.MessageReport:
		return l.MsgReport
	case specV1.MessageDesire:
		return l.MsgDesire
	}
	return func(c *common.Context) (interface{}, error) {
		return nil, common.Error(common.ErrResourceNotFound, common.Field("type", "messageType"))
	}
}

// MsgReport report
func (l *httpLink) MsgReport(c *common.Context) (interface{}, error) {
	switch c.GetHeader("kind") {
	case string(specV1.MessageDeviceReport):
		return l.deviceReport(c)
	case string(specV1.MessageDeviceLifecycleReport):
		return l.deviceLifecycleReport(c)
	default:
		return l.report(c)
	}
}

func (l *httpLink) report(c *common.Context) (interface{}, error) {
	msg, err := genReportMsg(c)
	if err != nil {
		return nil, err
	}
	resp, err := l.msgRouter[string(specV1.MessageReport)].(server.HandlerMessage)(*msg)
	if err != nil {
		return nil, err
	}
	return resp.Content.Value, nil
}

func (l *httpLink) deviceLifecycleReport(c *common.Context) (interface{}, error) {
	msg, err := genReportMsg(c)
	if err != nil {
		return nil, err
	}
	return l.msgRouter[string(specV1.MessageDeviceLifecycleReport)].(server.HandlerMessage)(*msg)
}

func (l *httpLink) deviceReport(c *common.Context) (interface{}, error) {
	msg, err := genReportMsg(c)
	if err != nil {
		return nil, err
	}
	return l.msgRouter[string(specV1.MessageDeviceReport)].(server.HandlerMessage)(*msg)
}

func genReportMsg(c *common.Context) (*specV1.Message, error) {
	msg, err := genDesireMsg(c)
	if err != nil {
		return nil, err
	}
	n := c.GetName()
	if n == "" {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}
	msg.Metadata["name"] = n
	return msg, nil
}

// MsgDesire desire
func (l *httpLink) MsgDesire(c *common.Context) (interface{}, error) {
	switch c.GetHeader("kind") {
	case string(specV1.NewMessageDeviceDesire), string(specV1.MessageDeviceDesire):
		return l.deviceDesire(c)
	default:
		return l.desire(c)
	}
}

func (l *httpLink) desire(c *common.Context) (interface{}, error) {
	msg, err := genDesireMsg(c)
	if err != nil {
		return nil, err
	}
	resp, err := l.msgRouter[string(specV1.MessageDesire)].(server.HandlerMessage)(*msg)
	if err != nil {
		return nil, err
	}
	return resp.Content.Value, nil
}

func (l *httpLink) deviceDesire(c *common.Context) (interface{}, error) {
	msg, err := genDesireMsg(c)
	if err != nil {
		return nil, err
	}
	if msg.Kind == specV1.MessageDeviceDesire {
		l.msgRouter[string(specV1.MessageDeviceDesire)].(server.HandlerMessage)(*msg)
	}
	return l.msgRouter[string(specV1.NewMessageDeviceDesire)].(server.HandlerMessage)(*msg)
}

func genDesireMsg(c *common.Context) (*specV1.Message, error) {
	ns := c.GetNamespace()
	if ns == "" {
		return nil, common.Error(common.ErrRequestParamInvalid)
	}

	body, err := c.GetRawData()
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}

	msg := specV1.Message{
		Content:  specV1.LazyValue{},
		Metadata: map[string]string{},
	}
	for k := range c.Request.Header {
		msg.Metadata[strings.ToLower(k)] = c.GetHeader(k)
	}
	msg.Metadata["namespace"] = ns
	msg.Metadata["clientIP"] = c.ClientIP()
	msg.Kind = specV1.MessageKind(msg.Metadata["kind"])

	if msg.Kind == "" {
		msg.Kind = specV1.MessageReport
	}

	err = msg.Content.UnmarshalJSON(body)
	if err != nil {
		return nil, common.Error(common.ErrRequestParamInvalid, common.Field("error", err.Error()))
	}
	return &msg, nil
}
