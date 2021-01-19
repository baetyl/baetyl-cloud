package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

var (
	moduletabales = []string{
		`
CREATE TABLE IF NOT EXISTS baetyl_module
(
  id          integer          PRIMARY KEY AUTOINCREMENT,
  name        varchar(255)     NOT NULL DEFAULT '',
  image       varchar(1024)    NOT NULL DEFAULT '',
  programs    varchar(2048)    NOT NULL DEFAULT '',
  version     varchar(36)      NOT NULL DEFAULT '',
  type        varchar(36)      NOT NULL DEFAULT '',
  is_hidden   int(1)           NOT NULL DEFAULT '0',
  description varchar(1024)    NOT NULL DEFAULT '',
  create_time timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`,
	}
)

func (d *DB) MockCreateModuleTable() {
	for _, sql := range moduletabales {
		_, err := d.Exec(nil, sql)
		if err != nil {
			panic(fmt.Sprintf("create module exception: %s", err.Error()))
		}
	}
}

func TestModule(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateModuleTable()

	module := &models.Module{
		Name:              "baetyl",
		Version:           "v2.0.0",
		Image:             "baetyl:v1",
		Programs: map[string]string{
			"linux-amd64": "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:              "a",
		IsHidden:          false,
		Description:       "for desp",
	}

	err = db.CreateModule(module)
	assert.NoError(t, err)

	res, err := db.GetModule(module.Name)
	assert.NoError(t, err)
	checkModule(t, module, res)

	module2 := &models.Module{
		Name:              "baetyl2",
		Version:           "v2.1.0",
		Image:             "baetyl:v2",
		Type:              "a",
		IsHidden:          false,
		Description:       "for desp",
	}

	err = db.CreateModule(module2)
	assert.NoError(t, err)

	res, err = db.GetModule(module2.Name)
	assert.NoError(t, err)
	checkModule(t, module2, res)

	module3 := &models.Module{
		Name:              "baetyl3",
		Version:           "v2.1.0",
		Image:             "baetyl:v2",
		Type:              "a",
		IsHidden:          true,
		Description:       "for desp",
	}

	err = db.CreateModule(module3)
	assert.NoError(t, err)

	res, err = db.GetModule(module3.Name)
	assert.NoError(t, err)
	checkModule(t, module3, res)

	module.IsHidden = true
	module.Description = "a"
	module.Type = "b"
	module.Image = "image-b"
	module.Version = "version"
	module.Programs = map[string]string{
		"b": "b-url",
	}
	err = db.UpdateModule(module)
	assert.NoError(t, err)

	res, err = db.GetModule(module.Name)
	assert.NoError(t, err)
	checkModule(t, module, res)

	module22 := &models.Module{
		Name:              "baetyl2",
		Version:           "v2.2.0",
		Image:             "baetyl:v2.2",
		Type:              "a",
		IsHidden:          false,
		Description:       "for desp",
	}

	err = db.CreateModule(module22)
	assert.NoError(t, err)

	res, err = db.GetModuleByVersion(module22.Name, module22.Version)
	assert.NoError(t, err)
	checkModule(t, module22, res)

	res, err = db.GetModuleByVersion(module2.Name, module2.Version)
	assert.NoError(t, err)
	checkModule(t, module2, res)
}

func TestListApp(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateModuleTable()

	module := &models.Module{
		Name:              "baetyl",
		Version:           "v2.0.0",
		Image:             "baetyl:v1",
		Programs: map[string]string{
			"linux-amd64": "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:              "a",
		IsHidden:          false,
		Description:       "for desp",
	}

	module2 := &models.Module{
		Name:              "baetyl2",
		Version:           "v2.0.2",
		Image:             "baetyl:v2",
		Programs: map[string]string{
			"linux-amd64": "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:              "b",
		IsHidden:          false,
		Description:       "for desp",
	}

	module3 := &models.Module{
		Name:              "baetyl2",
		Version:           "v2.0.2",
		Image:             "baetyl:v2",
		Programs: map[string]string{
			"linux-amd64": "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:              "a",
		IsHidden:          true,
		Description:       "for desp",
	}

	err = db.CreateModule(module)
	assert.NoError(t, err)

	err = db.CreateModule(module2)
	assert.NoError(t, err)

	err = db.CreateModule(module3)
	assert.NoError(t, err)

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
	}

	resList, err := db.ListModule(page)
	assert.NoError(t, err)
	assert.Len(t, resList, 2)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module2.Name, resList[1].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module2, &resList[1])

	page = &models.Filter{}
	resList, err = db.ListModule(page)
	assert.NoError(t, err)
	assert.Len(t, resList, 3)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module2.Name, resList[1].Name)
	assert.Equal(t, module3.Name, resList[2].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module2, &resList[1])
	checkModule(t, module3, &resList[2])

	resList, err = db.ListModuleWithOptions("a", false, page)
	assert.NoError(t, err)
	assert.Len(t, resList, 1)
	assert.Equal(t, module.Name, resList[0].Name)
	checkModule(t, module, &resList[0])

	resList, err = db.ListModuleWithOptions("a", true, page)
	assert.NoError(t, err)
	assert.Len(t, resList, 1)
	assert.Equal(t, module3.Name, resList[0].Name)
	checkModule(t, module3, &resList[0])

	resList, err = db.ListModuleWithOptions("b", false, page)
	assert.NoError(t, err)
	assert.Len(t, resList, 1)
	assert.Equal(t, module2.Name, resList[0].Name)
	checkModule(t, module2, &resList[0])

	resList, err = db.ListModuleWithOptions("b", true, page)
	assert.NoError(t, err)
	assert.Len(t, resList, 0)
}

func checkModule(t *testing.T, expect, actual *models.Module) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Version, actual.Version)
	assert.Equal(t, expect.Image, actual.Image)
	assert.Equal(t, expect.Programs, actual.Programs)
	assert.EqualValues(t, expect.Type, actual.Type)
	assert.Equal(t, expect.IsHidden, actual.IsHidden)
	assert.Equal(t, expect.Description, actual.Description)
}