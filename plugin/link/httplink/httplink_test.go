package httplink

import (
	"io/ioutil"
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
)

func genHTTPLinkConf(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(tempDir, "config.yml"), []byte(configYaml), 777)
	assert.NoError(t, err)
	return tempDir
}

type handler struct {
	t            *testing.T
	expReportMsg specV1.Message
	expDesireMsg specV1.Message
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

	handler := &handler{
		t:            t,
		expReportMsg: reportMsg,
		expDesireMsg: desiretMsg,
	}

	link.AddMsgRouter(string(specV1.MessageReport), server.HandlerMessage(handler.report))
	link.AddMsgRouter(string(specV1.MessageDesire), server.HandlerMessage(handler.desire))

	go link.Start()

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

	err = link.Close()
	assert.NoError(t, err)
}
