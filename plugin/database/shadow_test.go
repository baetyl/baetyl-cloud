package database

import (
	"fmt"
	"testing"

	"github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	shaodowTables = []string{
		`
CREATE TABLE baetyl_node_shadow(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(128) NOT NULL DEFAULT '',
    namespace   VARCHAR(64) NOT NULL DEFAULT '',
    create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    report      BLOB,
    desire      BLOB 
);
`,
	}
)

func (d *DB) MockCreateShadowTable() {
	for _, sql := range shaodowTables {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create table exception: %s", err.Error()))
		}
	}
}

func TestShadow(t *testing.T) {
	db, err := MockNewDB()
	isSysApp := false
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateShadowTable()

	namespace := "test"
	shadow := &models.Shadow{
		Namespace: namespace,
		Name:      "node01",
		Desire: v1.Desire{
			"apps": []v1.AppInfo{
				{
					Name:    "app01",
					Version: "1",
				},
			},
		},
		Report: v1.Report{
			"apps": []v1.AppInfo{
				{
					Name:    "app01",
					Version: "1",
				},
			},
		},
	}
	result, err := db.Create(shadow)
	assert.NoError(t, err)
	assert.Equal(t, shadow.Name, result.Name)
	assert.Equal(t, shadow.Desire.AppInfos(isSysApp), result.Desire.AppInfos(isSysApp))
	assert.Equal(t, shadow.Report.AppInfos(isSysApp), result.Report.AppInfos(isSysApp))

	n, err := db.Get(namespace, shadow.Name)
	assert.NoError(t, err)
	assert.Equal(t, shadow.Name, n.Name)
	assert.Equal(t, shadow.Desire.AppInfos(isSysApp), n.Desire.AppInfos(isSysApp))
	assert.Equal(t, shadow.Report.AppInfos(isSysApp), n.Report.AppInfos(isSysApp))

	report := v1.Report{
		"apps": []v1.AppInfo{
			{
				Name:    "app02",
				Version: "2",
			},
		},
	}

	shadow.Report = report
	result, err = db.UpdateReport(shadow)
	assert.NoError(t, err)
	assert.Equal(t, report.AppInfos(isSysApp), result.Report.AppInfos(isSysApp))
	assert.Equal(t, shadow.Desire.AppInfos(isSysApp), result.Desire.AppInfos(isSysApp))

	desire := v1.Desire{
		"apps": []v1.AppInfo{
			{
				Name:    "app02",
				Version: "2",
			},
		},
	}

	shadow.Desire = desire
	result, err = db.UpdateDesire(shadow)
	assert.NoError(t, err)
	assert.Equal(t, desire.AppInfos(isSysApp), result.Desire.AppInfos(isSysApp))
	assert.Equal(t, report.AppInfos(isSysApp), result.Report.AppInfos(isSysApp))

	nodeList := &models.NodeList{
		Items: []v1.Node{
			{
				Namespace: namespace,
				Name:      "node01",
			},
		},
	}

	list, err := db.List(namespace, nodeList)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list.Items))
	assert.Equal(t, "node01", list.Items[0].Name)

	err = db.Delete(namespace, shadow.Name)
	assert.NoError(t, err)

}
