package entities

import (
	"encoding/json"
	"time"

	"github.com/baetyl/baetyl-go/v2/log"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Module struct {
	Id          uint64    `db:"id"`
	Name        string    `db:"name"`
	Version     string    `db:"version"`
	Image       string    `db:"image"`
	Programs    string    `db:"programs"`
	Type        string    `db:"type"`
	Flag        int       `db:"flag"`
	IsLatest    bool      `db:"is_latest"`
	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}

func ToModuleModel(module *Module) (*models.Module, error) {
	m := &models.Module{
		Name:              module.Name,
		Version:           module.Version,
		Image:             module.Image,
		Programs:          make(map[string]string),
		Type:              module.Type,
		Flag:              module.Flag,
		IsLatest:          module.IsLatest,
		Description:       module.Description,
		CreationTimestamp: module.CreateTime,
		UpdateTimestamp:   module.UpdateTime,
	}

	if module.Programs != "" {
		err := json.Unmarshal([]byte(module.Programs), &m.Programs)
		if err != nil {
			log.L().Error("module db to module error",
				log.Any("name", module.Name),
				log.Any("version", module.Version))
			return nil, err
		}
	}
	return m, nil
}

func FromModuleModel(module *models.Module) (*Module, error) {
	s, err := json.Marshal(module.Programs)
	if err != nil {
		log.L().Error("module translate to db model error",
			log.Any("name", module.Name),
			log.Any("version", module.Version))
		return nil, err
	}

	app := &Module{
		Name:        module.Name,
		Version:     module.Version,
		Image:       module.Image,
		Programs:    string(s),
		Type:        module.Type,
		Flag:        module.Flag,
		IsLatest:    module.IsLatest,
		Description: module.Description,
		CreateTime:  time.Time{},
		UpdateTime:  time.Time{},
	}
	return app, nil
}
