package entities

import (
	"encoding/json"
	"github.com/baetyl/baetyl-go/log"
	specV1 "github.com/baetyl/baetyl-go/spec/v1"
	"time"
)

type Application struct {
	Id         uint64    `db:"id"`
	Namespace  string    `db:"namespace"`
	Name       string    `db:"name"`
	Version    string    `db:"version"`
	IsDeleted  int       `db:"is_deleted"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
	Content    string    `db:"content"`
}

func ToApplicationModel(application *Application) (*specV1.Application, error) {

	app := &specV1.Application{
		Name:      application.Name,
		Namespace: application.Namespace,
		Version:   application.Version,
	}
	if application.Content != "" {
		err := json.Unmarshal([]byte(application.Content), app)
		if err != nil {
			log.L().Error("application db to application error",
				log.Any("namespace", application.Namespace),
				log.Any("name", application.Name),
				log.Any("version", application.Version))
			return nil, err
		}
	}

	return app, nil
}

func FromApplicationModel(application *specV1.Application) (*Application, error) {
	content, err := json.Marshal(application)
	if err != nil {
		log.L().Error("application translate to db model error",
			log.Any("namespace", application.Namespace),
			log.Any("name", application.Name),
			log.Any("version", application.Version))
		return nil, err
	}

	app := &Application{
		Name:      application.Name,
		Namespace: application.Namespace,
		Version:   application.Version,
		Content:   string(content),
	}
	return app, nil
}
