package httplink

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/http"
	"github.com/baetyl/baetyl-go/v2/json"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/baetyl/baetyl-cloud/v2/server"
)

const (
	configYaml = `
httplink:
  commonName: "cn"
`
)

var (
	reportMsg = specV1.Message{
		Kind:    specV1.MessageReport,
		Content: specV1.LazyValue{Value: map[string]string{"1": "2"}},
	}

	desiretMsg = specV1.Message{
		Kind:    specV1.MessageDesire,
		Content: specV1.LazyValue{Value: map[string]string{"3": "4"}},
	}

	deviceReportMsg = specV1.Message{
		Kind:    specV1.MessageDeviceReport,
		Content: specV1.LazyValue{Value: map[string]string{"5": "6"}},
	}

	deviceLifecycReportMsg = specV1.Message{
		Kind:    specV1.MessageDeviceLifecycleReport,
		Content: specV1.LazyValue{Value: map[string]string{"7": "8"}},
	}

	deviceDesiretMsg = specV1.Message{
		Kind:    specV1.MessageDeviceDesire,
		Content: specV1.LazyValue{Value: map[string]string{"9": "0"}},
	}
)

func genHTTPLinkConf(t *testing.T) string {
	tempDir := t.TempDir()
	err := os.WriteFile(path.Join(tempDir, "config.yml"), []byte(configYaml), 777)
	assert.NoError(t, err)
	return tempDir
}

type handler struct {
	t                           *testing.T
	expReportMsg                specV1.Message
	expDesireMsg                specV1.Message
	expDeviceReportMsg          specV1.Message
	expDeviceLifecycleReportMsg specV1.Message
	expDeviceDesireMsg          specV1.Message
}

func (h *handler) report(m specV1.Message) (*specV1.Message, error) {
	res := map[string]string{}
	err := m.Content.Unmarshal(&res)
	assert.NoError(h.t, err)
	assert.EqualValues(h.t, h.expReportMsg.Content.Value, res)
	return &reportMsg, nil
}

func (h *handler) desire(m specV1.Message) (*specV1.Message, error) {
	res := map[string]string{}
	err := m.Content.Unmarshal(&res)
	assert.NoError(h.t, err)
	assert.EqualValues(h.t, h.expDesireMsg.Content.Value, res)
	return &desiretMsg, nil
}

func (h *handler) deviceReport(m specV1.Message) (*specV1.Message, error) {
	res := map[string]string{}
	err := m.Content.Unmarshal(&res)
	assert.NoError(h.t, err)
	assert.EqualValues(h.t, h.expDeviceReportMsg.Content.Value, res)
	return &deviceReportMsg, nil
}

func (h *handler) deviceLifecycleReport(m specV1.Message) (*specV1.Message, error) {
	res := map[string]string{}
	err := m.Content.Unmarshal(&res)
	assert.NoError(h.t, err)
	assert.EqualValues(h.t, h.expDeviceLifecycleReportMsg.Content.Value, res)
	return &deviceLifecycReportMsg, nil
}

func (h *handler) deviceDesire(m specV1.Message) (*specV1.Message, error) {
	res := map[string]string{}
	err := m.Content.Unmarshal(&res)
	assert.NoError(h.t, err)
	assert.EqualValues(h.t, h.expDeviceDesireMsg.Content.Value, res)
	return &deviceDesiretMsg, nil
}

func TestNewHTTPLink(t *testing.T) {
	cfg := &CloudConfig{}
	common.SetConfFile(path.Join(genHTTPLinkConf(t), "config.yml"))
	err := common.LoadConfig(cfg)
	assert.NoError(t, err)

	err = os.Setenv("HTTP_LINK_PORT", "9939")
	assert.NoError(t, err)

	pl, err := NewHTTPLink()
	assert.NoError(t, err)
	assert.NotNil(t, pl)

	link, ok := pl.(plugin.SyncLink)
	assert.True(t, ok)

	hd := &handler{
		t:                           t,
		expReportMsg:                reportMsg,
		expDesireMsg:                desiretMsg,
		expDeviceReportMsg:          deviceReportMsg,
		expDeviceLifecycleReportMsg: deviceLifecycReportMsg,
		expDeviceDesireMsg:          deviceDesiretMsg,
	}

	link.AddMsgRouter(string(specV1.MessageReport), server.HandlerMessage(hd.report))
	link.AddMsgRouter(string(specV1.MessageDesire), server.HandlerMessage(hd.desire))
	link.AddMsgRouter(string(specV1.MessageDeviceReport), server.HandlerMessage(hd.deviceReport))
	link.AddMsgRouter(string(specV1.MessageDeviceLifecycleReport), server.HandlerMessage(hd.deviceLifecycleReport))
	link.AddMsgRouter(string(specV1.MessageDeviceDesire), server.HandlerMessage(hd.deviceDesire))
	link.AddMsgRouter(string(specV1.NewMessageDeviceDesire), server.HandlerMessage(hd.deviceDesire))

	go link.Start()
	defer func() {
		err = link.Close()
		assert.NoError(t, err)
	}()

	cli := http.NewClient(&http.ClientOptions{
		Address: "http://0.0.0.0:9939",
	})

	// wait server start
	for {
		_, err := cli.GetURL("http://0.0.0.0:9939/health")
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// report
	dt, err := json.Marshal(reportMsg.Content.Value)
	assert.NoError(t, err)
	resp, err := cli.PostJSON("v1/sync/report", dt, map[string]string{"cn": "default.test"})
	assert.NoError(t, err)

	reportResp := map[string]string{}
	err = json.Unmarshal(resp, &reportResp)
	assert.NoError(t, err)
	assert.EqualValues(t, reportMsg.Content.Value, reportResp)

	// desire
	dt, err = json.Marshal(desiretMsg.Content.Value)
	assert.NoError(t, err)
	resp, err = cli.PostJSON("v1/sync/desire", dt, map[string]string{"cn": "default.test"})
	assert.NoError(t, err)

	desireResp := map[string]string{}
	err = json.Unmarshal(resp, &desireResp)
	assert.NoError(t, err)
	assert.EqualValues(t, desiretMsg.Content.Value, desireResp)

	// MessageDeviceReport
	dt, err = json.Marshal(deviceReportMsg.Content.Value)
	assert.NoError(t, err)
	resp, err = cli.PostJSON("v1/sync/report", dt, map[string]string{
		"cn":   "default.test",
		"kind": string(specV1.MessageDeviceReport),
	})
	assert.NoError(t, err)

	drResp := specV1.Message{}
	err = json.Unmarshal(resp, &drResp)
	assert.NoError(t, err)
	data := map[string]string{}
	err = drResp.Content.Unmarshal(&data)
	assert.NoError(t, err)
	assert.Equal(t, deviceReportMsg.Kind, drResp.Kind)
	assert.EqualValues(t, deviceReportMsg.Content.Value, data)

	// MessageDeviceLifecycleReport
	dt, err = json.Marshal(deviceLifecycReportMsg.Content.Value)
	assert.NoError(t, err)
	resp, err = cli.PostJSON("v1/sync/report", dt, map[string]string{
		"cn":   "default.test",
		"kind": string(specV1.MessageDeviceLifecycleReport),
	})
	assert.NoError(t, err)

	drResp = specV1.Message{}
	err = json.Unmarshal(resp, &drResp)
	assert.NoError(t, err)
	data = map[string]string{}
	err = drResp.Content.Unmarshal(&data)
	assert.NoError(t, err)
	assert.Equal(t, deviceLifecycReportMsg.Kind, drResp.Kind)
	assert.EqualValues(t, deviceLifecycReportMsg.Content.Value, data)

	dt, err = json.Marshal(deviceDesiretMsg.Content.Value)
	assert.NoError(t, err)

	resp, err = cli.PostJSON("v1/sync/desire", dt, map[string]string{
		"cn":   "default.test",
		"kind": string(specV1.MessageDeviceDesire),
	})
	assert.NoError(t, err)

	resp, err = cli.PostJSON("v1/sync/desire", dt, map[string]string{
		"cn":   "default.test",
		"kind": string(specV1.NewMessageDeviceDesire),
	})
	assert.NoError(t, err)

	drResp = specV1.Message{}
	err = json.Unmarshal(resp, &drResp)
	assert.NoError(t, err)
	data = map[string]string{}
	err = drResp.Content.Unmarshal(&data)
	assert.NoError(t, err)
	assert.Equal(t, deviceDesiretMsg.Kind, drResp.Kind)
	assert.EqualValues(t, deviceDesiretMsg.Content.Value, data)
}
