package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
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
  flag        int(10)          NOT NULL DEFAULT '0',
  is_latest   int(1)           NOT NULL DEFAULT '0',
  description varchar(1024)    NOT NULL DEFAULT '',
  create_time timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`, `
CREATE UNIQUE INDEX name_version
on baetyl_module (name, version);`,
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
		Name:    "baetyl",
		Version: "v2.0.0",
		Image:   "baetyl:v1",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "a",
		IsLatest:    false,
		Description: "for desp",
	}

	res, err := db.CreateModule(module)
	assert.NoError(t, err)
	checkModule(t, module, res)

	res, err = db.GetModuleByVersion(module.Name, module.Version)
	assert.NoError(t, err)
	checkModule(t, module, res)

	res, err = db.CreateModule(module)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "UNIQUE constraint failed: baetyl_module.name, baetyl_module.version")

	module01 := &models.Module{
		Name:    "baetyl",
		Version: "v2.1.0",
		Image:   "baetyl:v11",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "a",
		IsLatest:    false,
		Description: "for desp",
	}

	res, err = db.CreateModule(module01)
	assert.NoError(t, err)
	checkModule(t, module01, res)

	res, err = db.GetModuleByVersion(module01.Name, module01.Version)
	assert.NoError(t, err)
	checkModule(t, module01, res)

	ress, err := db.GetModules(module01.Name)
	assert.NoError(t, err)
	assert.Equal(t, len(ress), 2)
	checkModule(t, &ress[0], module)
	checkModule(t, &ress[1], module01)

	res, err = db.GetModuleByImage(module01.Name, module01.Image)
	assert.NoError(t, err)
	checkModule(t, module01, res)

	res, err = db.GetModuleByImage(module.Name, module.Image)
	assert.NoError(t, err)
	checkModule(t, module, res)

	module2 := &models.Module{
		Name:        "baetyl2",
		Version:     "v2.1.0",
		Image:       "baetyl:v2",
		Type:        "a",
		IsLatest:    false,
		Description: "for desp",
	}

	res, err = db.CreateModule(module2)
	assert.NoError(t, err)
	checkModule(t, module2, res)

	res, err = db.GetModuleByVersion(module2.Name, module2.Version)
	assert.NoError(t, err)
	checkModule(t, module2, res)

	module3 := &models.Module{
		Name:        "baetyl3",
		Version:     "v2.1.0",
		Image:       "baetyl:v2",
		Type:        "a",
		IsLatest:    true,
		Description: "for desp",
	}

	res, err = db.CreateModule(module3)
	assert.NoError(t, err)
	checkModule(t, module3, res)

	res, err = db.GetModuleByVersion(module3.Name, module3.Version)
	assert.NoError(t, err)
	checkModule(t, module3, res)

	res, err = db.GetLatestModule(module3.Name)
	assert.NoError(t, err)
	checkModule(t, module3, res)

	module31 := &models.Module{
		Name:    "baetyl3",
		Version: "v2.2.0",
		Image:   "baetyl:v2.2",
		Type:    "a",
		Programs: map[string]string{
			"linux/amd64": "aaa-url",
		},
		IsLatest:    true,
		Description: "for desp",
	}

	res, err = db.CreateModule(module31)
	assert.NoError(t, err)
	checkModule(t, module31, res)

	res, err = db.GetLatestModule(module3.Name)
	assert.NoError(t, err)
	checkModule(t, module31, res)

	i, err := db.GetLatestModuleImage(module3.Name)
	assert.NoError(t, err)
	assert.Equal(t, i, module31.Image)

	p, err := db.GetLatestModuleProgram(module3.Name, "linux/amd64")
	assert.NoError(t, err)
	assert.Equal(t, p, "aaa-url")

	module31.Description = "a"
	module31.Programs = map[string]string{
		"b": "b-url",
	}
	res, err = db.UpdateModuleByVersion(module31)
	assert.NoError(t, err)
	checkModule(t, module31, res)

	res, err = db.GetLatestModule(module3.Name)
	assert.NoError(t, err)
	checkModule(t, module31, res)

	module3.Description = "aaa"
	module3.Programs = map[string]string{
		"aaa": "aaa-url",
	}
	module3.IsLatest = true
	res, err = db.UpdateModuleByVersion(module3)
	assert.NoError(t, err)
	checkModule(t, module3, res)

	res, err = db.GetLatestModule(module3.Name)
	assert.NoError(t, err)
	checkModule(t, module3, res)

	ress, err = db.GetModules(module3.Name)
	assert.NoError(t, err)
	assert.Equal(t, len(ress), 2)
	checkModule(t, &ress[0], module3)
	module31.IsLatest = false
	checkModule(t, &ress[1], module31)

	err = db.DeleteModuleByVersion(module3.Name, module3.Version)
	assert.NoError(t, err)

	ress, err = db.GetModules(module3.Name)
	assert.NoError(t, err)
	assert.Equal(t, len(ress), 1)
	checkModule(t, &ress[0], module31)

	err = db.DeleteModules(module3.Name)
	assert.NoError(t, err)

	ress, err = db.GetModules(module3.Name)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "The (module) resource (baetyl3) is not found.")
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
		Name:    "baetyl",
		Version: "v2.0.1",
		Image:   "baetyl:v1",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "a",
		IsLatest:    false,
		Description: "for desp",
	}

	module2 := &models.Module{
		Name:    "baetyl2",
		Version: "v2.0.2",
		Image:   "baetyl:v2",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "b",
		IsLatest:    false,
		Description: "for desp",
	}

	module3 := &models.Module{
		Name:    "baetyl3",
		Version: "v2.0.3",
		Image:   "baetyl:v3",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "a",
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module)
	assert.NoError(t, err)

	_, err = db.CreateModule(module2)
	assert.NoError(t, err)

	_, err = db.CreateModule(module3)
	assert.NoError(t, err)

	page := &models.Filter{
		PageNo:   1,
		PageSize: 2,
	}

	resList, err := db.ListModules(page, common.ModuleType("a"))
	assert.NoError(t, err)
	assert.Len(t, resList, 2)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module2.Name, resList[1].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module2, &resList[1])

	page = &models.Filter{}
	resList, err = db.ListModules(page, common.ModuleType("a"))
	assert.NoError(t, err)
	assert.Len(t, resList, 3)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module2.Name, resList[1].Name)
	assert.Equal(t, module3.Name, resList[2].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module2, &resList[1])
	checkModule(t, module3, &resList[2])

	resList, err = db.listModulesByTypeTx(nil, "a", page)
	assert.NoError(t, err)
	assert.Len(t, resList, 1)
	assert.Equal(t, module3.Name, resList[0].Name)
	checkModule(t, module3, &resList[0])

	resList, err = db.listModulesByTypeTx(nil, "b", page)
	assert.NoError(t, err)
	assert.Len(t, resList, 0)

	module4 := &models.Module{
		Name:    "baetyl4",
		Version: "v2.0.4",
		Image:   "baetyl:v4",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "b",
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module4)
	assert.NoError(t, err)

	module5 := &models.Module{
		Name:    "baetyl5",
		Version: "v2.0.5",
		Image:   "baetyl:v5",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        "b",
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module5)
	assert.NoError(t, err)

	resList, err = db.listModulesByTypeTx(nil, "b", page)
	assert.NoError(t, err)
	assert.Len(t, resList, 2)

	module6 := &models.Module{
		Name:    "baetyl6",
		Version: "v2.0.6",
		Image:   "baetyl:v6",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeUserRuntime),
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module6)
	assert.NoError(t, err)

	module7 := &models.Module{
		Name:    "baetyl7",
		Version: "v2.0.7",
		Image:   "baetyl:v7",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeUserRuntime),
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module7)
	assert.NoError(t, err)

	resList, err = db.listModulesByTypeTx(nil, common.TypeUserRuntime, page)
	assert.NoError(t, err)
	assert.Len(t, resList, 2)

	module8 := &models.Module{
		Name:    "baetyl8",
		Version: "v2.0.8",
		Image:   "baetyl:v8",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeSystemOptional),
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module8)
	assert.NoError(t, err)

	resList, err = db.listModulesByTypeTx(nil, common.TypeSystemOptional, page)
	assert.NoError(t, err)
	assert.Len(t, resList, 1)
}

func TestListModule(t *testing.T) {
	db, err := MockNewDB()
	if err != nil {
		fmt.Printf("get mock sqlite3 error = %s", err.Error())
		t.Fail()
		return
	}
	db.MockCreateModuleTable()

	module := &models.Module{
		Name:    "baetyl",
		Version: "v2.0.1",
		Image:   "baetyl:v1",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeSystemOptional),
		IsLatest:    true,
		Description: "for desp",
	}

	module2 := &models.Module{
		Name:    "baetyl2",
		Version: "v2.0.2",
		Image:   "baetyl:v2",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeSystemOptional),
		IsLatest:    false,
		Description: "for desp",
	}

	module3 := &models.Module{
		Name:    "baetyl3",
		Version: "v2.0.3",
		Image:   "baetyl:v3",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeSystemKube),
		IsLatest:    true,
		Description: "for desp",
	}

	module4 := &models.Module{
		Name:    "baetyl4",
		Version: "v2.0.4",
		Image:   "baetyl:v4",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeSystemNative),
		IsLatest:    true,
		Description: "for desp",
	}

	module5 := &models.Module{
		Name:    "baetyl6",
		Version: "v2.0.6",
		Image:   "baetyl:v6",
		Programs: map[string]string{
			"linux-amd64":    "url-linux-amd64",
			"linux-arm64-v8": "url-linux-arm64-v8",
		},
		Type:        string(common.TypeUserRuntime),
		IsLatest:    true,
		Description: "for desp",
	}

	_, err = db.CreateModule(module)
	assert.NoError(t, err)
	_, err = db.CreateModule(module2)
	assert.NoError(t, err)
	_, err = db.CreateModule(module3)
	assert.NoError(t, err)
	_, err = db.CreateModule(module4)
	assert.NoError(t, err)
	_, err = db.CreateModule(module5)
	assert.NoError(t, err)

	page := &models.Filter{}
	resList, err := db.ListModules(page, common.TypeSystemOptional)
	assert.NoError(t, err)
	assert.Len(t, resList, 3)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module3.Name, resList[1].Name)
	assert.Equal(t, module4.Name, resList[2].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module3, &resList[1])
	checkModule(t, module4, &resList[2])

	resList, err = db.ListModules(page, common.TypeSystemKube)
	assert.NoError(t, err)
	assert.Len(t, resList, 2)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module3.Name, resList[1].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module3, &resList[1])

	resList, err = db.ListModules(page, common.TypeSystemNative)
	assert.NoError(t, err)
	assert.Len(t, resList, 2)
	assert.Equal(t, module.Name, resList[0].Name)
	assert.Equal(t, module4.Name, resList[1].Name)
	checkModule(t, module, &resList[0])
	checkModule(t, module4, &resList[1])

	resList, err = db.ListModules(page, common.TypeUserRuntime)
	assert.NoError(t, err)
	assert.Len(t, resList, 1)
	assert.Equal(t, module5.Name, resList[0].Name)
	checkModule(t, module5, &resList[0])
}

func checkModule(t *testing.T, expect, actual *models.Module) {
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.Version, actual.Version)
	assert.Equal(t, expect.Image, actual.Image)
	assert.Equal(t, expect.Programs, actual.Programs)
	assert.EqualValues(t, expect.Type, actual.Type)
	assert.Equal(t, expect.IsLatest, actual.IsLatest)
	assert.Equal(t, expect.Description, actual.Description)
}
