package entities

import (
	"testing"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/stretchr/testify/assert"
)

func TestFromApplicationModel(t *testing.T) {
	app := &Application{
		Name:      "testApp",
		Namespace: "namespace",
		Content:   `{"name":"testApp","namespace":"namespace","createTime":"0001-01-01T00:00:00Z","updateTime":"0001-01-01T00:00:00Z","cronTime":"0001-01-01T00:00:00Z","replica":0,"ota":{}}`,
	}

	mApp := &specV1.Application{
		Namespace: "namespace",
		Name:      "testApp",
	}
	modelApp, err := FromApplicationModel(mApp)

	assert.NoError(t, err)
	assert.Equal(t, app, modelApp)

}

func TestToApplicationModel(t *testing.T) {
	app := &Application{
		Id:        123,
		Name:      "testApp",
		Namespace: "namespace",
		Content:   `{"name":"testApp","namespace":"namespace"}`,
	}

	modelApp, err := ToApplicationModel(app)
	assert.NoError(t, err)
	assert.Equal(t, app.Name, modelApp.Name)
	assert.Equal(t, app.Namespace, modelApp.Namespace)

	app = &Application{
		Id:        123,
		Name:      "testApp",
		Namespace: "namespace",
		Content:   `{"name":"testApp","namespace":"namespace"''}`,
	}

	_, err = ToApplicationModel(app)
	assert.NotNil(t, err)
}
