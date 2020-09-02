package service

import (
	"testing"

	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/config"
)

func TestTemplateServiceImpl_UnmarshalTemplate(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	cfg := &config.CloudConfig{}
	cfg.Template.Path = "../scripts/native/templates"
	funcs := map[string]interface{}{
		"GetProperty": func(in string) string {
			return "out-" + in
		},
	}
	sTemplate, err := NewTemplateService(cfg, funcs)
	assert.NoError(t, err)

	params := map[string]interface{}{
		"Namespace":   "ns",
		"NodeName":    "node-name",
		"CoreAppName": "core-app",
	}
	tests := []struct {
		name    string
		out     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "baetyl-core-app.yml",
			out:  &v1.Application{},
			want: &v1.Application{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sTemplate.UnmarshalTemplate(tt.name, params, tt.out)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, tt.out)
		})
	}
}
