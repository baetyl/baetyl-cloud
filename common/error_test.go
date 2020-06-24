package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	type args struct {
		c  Code
		fs []*F
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "ErrRequestAccessDenied-1",
			args: args{
				c: ErrRequestAccessDenied,
			},
			wantErr: "The request access is denied.",
		},
		{
			name: "ErrRequestAccessDenied-2",
			args: args{
				c:  ErrRequestAccessDenied,
				fs: []*F{Field("dummy", "dummy")},
			},
			wantErr: "The request access is denied.",
		},
		{
			name: "ErrRequestParamInvalid-1",
			args: args{
				c: ErrRequestParamInvalid,
			},
			wantErr: "The request parameter is invalid.",
		},
		{
			name: "ErrRequestParamInvalid-2",
			args: args{
				c:  ErrRequestParamInvalid,
				fs: []*F{Field("error", "missing name")},
			},
			wantErr: "The request parameter is invalid. (missing name)",
		},
		{
			name: "ErrResourceNotFound-1",
			args: args{
				c: ErrResourceNotFound,
			},
			wantErr: "The resource is not found.",
		},
		{
			name: "ErrResourceNotFound-2",
			args: args{
				c:  ErrResourceNotFound,
				fs: []*F{Field("name", "xxx")},
			},
			wantErr: "The resource (xxx) is not found.",
		},
		{
			name: "ErrVolumeType-1",
			args: args{
				c: ErrVolumeType,
			},
			wantErr: "The volume type should be.",
		},
		{
			name: "ErrVolumeType-2",
			args: args{
				c:  ErrVolumeType,
				fs: []*F{Field("type", "yyy"), Field("name", "baetyl")},
			},
			wantErr: "The volume (baetyl) type should be (yyy).",
		},
		{
			name: "ErrUnknown-1",
			args: args{
				c: ErrUnknown,
			},
			wantErr: "There is a unknown error. If the attempt to retry does not work, please contact us.",
		},
		{
			name: "ErrUnknown-2",
			args: args{
				c:  ErrUnknown,
				fs: []*F{Field("error", "zzz")},
			},
			wantErr: "There is a unknown error (zzz). If the attempt to retry does not work, please contact us.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error(tt.args.c, tt.args.fs...)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}
