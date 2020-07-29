package entities

import (
	"encoding/json"
	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewShadowFromShadowModel(t *testing.T) {
	shadow := models.NewShadow("default", "node01")
	shadow.Desire[common.DesiredApplications] = []v1.AppInfo{
		{
			"app01",
			"v1",
		},
	}
	shadow.Report[common.DesiredSysApplications] = []v1.AppInfo{
		{
			"sysapp01",
			"v1",
		},
	}

	shd, err := NewShadowFromShadowModel(shadow)
	assert.NoError(t, err)
	desire, _ := json.Marshal(&shadow.Desire)
	report, _ := json.Marshal(&shadow.Report)
	assert.Equal(t, string(report), shd.Report)
	assert.Equal(t, string(desire), shd.Desire)

}

func TestShadow_ToShadowModel(t *testing.T) {
	shadow := &Shadow{
		Namespace: "default",
		Name:      "node01",
		Desire:    `{"apps":[],"sysapps":[{"name":"baetyl-core-node-test-2-qcfpywxuh","version":"211430"},{"name":"baetyl-function-node-test-2-hb8ibamcv","version":"211434"}]}`,
		Report:    `{"apps":[],"sysapps":[{"name":"baetyl-core-node-test-2-qcfpywxuh","version":"211430"},{"name":"baetyl-function-node-test-2-hb8ibamcv","version":"211434"}]}`,
	}

	shd, err := shadow.ToShadowModel()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(shd.Desire.AppInfos(false)))
	assert.Equal(t, 0, len(shd.Report.AppInfos(false)))
	assert.Equal(t, 2, len(shd.Desire.AppInfos(true)))
	assert.Equal(t, 2, len(shd.Report.AppInfos(true)))
}
