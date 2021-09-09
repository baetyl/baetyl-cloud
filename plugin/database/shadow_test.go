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
    desire_version VARCHAR(36) NOT NULL DEFAULT '',
    report      BLOB,
    desire      BLOB,
    report_meta BLOB,
    desire_meta BLOB
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
	assert.NoError(t, err)
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
		ReportMeta: map[string]interface{}{},
		DesireMeta: map[string]interface{}{},
	}
	result, err := db.Create(nil, shadow)
	assert.NoError(t, err)
	assert.Equal(t, shadow.Name, result.Name)
	assert.Equal(t, shadow.Desire.AppInfos(isSysApp), result.Desire.AppInfos(isSysApp))
	assert.Equal(t, shadow.Report.AppInfos(isSysApp), result.Report.AppInfos(isSysApp))

	n, err := db.Get(nil, namespace, shadow.Name)
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
	result, err = db.GetShadowTx(nil, shadow.Namespace, shadow.Name)
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
	err = db.UpdateDesire(nil, shadow)
	assert.NoError(t, err)
	result, err = db.GetShadowTx(nil, shadow.Namespace, shadow.Name)
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

func TestShadowTx(t *testing.T) {
	db, err := MockNewDB()
	assert.NoError(t, err)
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
		ReportMeta: map[string]interface{}{},
		DesireMeta: map[string]interface{}{},
	}

	tx, err := db.BeginTx()
	assert.NoError(t, err)

	result, err := db.Create(tx, shadow)
	assert.NoError(t, err)
	assert.Equal(t, shadow.Name, result.Name)
	err = tx.Commit()
	assert.NoError(t, err)

	desire := v1.Desire{
		"apps": []v1.AppInfo{
			{
				Name:    "app02",
				Version: "2",
			},
		},
	}

	shadow.Desire = desire
	tx, err = db.BeginTx()
	assert.NoError(t, err)

	err = db.UpdateDesire(tx, shadow)
	assert.NoError(t, err)
	err = tx.Commit()
	assert.NoError(t, err)
}

func TestBatchShadows(t *testing.T)  {
	db, err := MockNewDB()
	assert.NoError(t, err)
	db.MockCreateShadowTable()

	namespace := "test"
	shadow1 := &models.Shadow{
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
		Report: v1.Report{},
		ReportMeta: map[string]interface{}{},
		DesireMeta: map[string]interface{}{},
	}

	shadow2 := &models.Shadow{
		Namespace: namespace,
		Name:      "node02",
		Desire: v1.Desire{
			"apps": []v1.AppInfo{
				{
					Name:    "app01",
					Version: "1",
				},
			},
		},
		Report: v1.Report{},
		ReportMeta: map[string]interface{}{},
		DesireMeta: map[string]interface{}{},
	}

	tx, err := db.BeginTx()
	assert.NoError(t, err)

	_, err = db.Create(tx, shadow1)
	assert.NoError(t, err)

	_, err = db.Create(tx, shadow2)
	assert.NoError(t, err)

	names := []string{shadow1.Name, shadow2.Name}
	_, err = db.ListShadowByNames(tx, namespace, nil)
	assert.NoError(t, err)
	shadows, err := db.ListShadowByNames(tx, namespace, names)
	assert.NoError(t, err)
	assert.Equal(t, len(shadows), 2)

	err = tx.Commit()
	assert.NoError(t, err)

	updateShadows := []*models.Shadow{shadow1, shadow2}
	err = db.UpdateDesires(nil, nil)
	assert.NoError(t, err)
	err = db.UpdateDesires(nil, updateShadows)
	assert.NotNil(t, err)
}
