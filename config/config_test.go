package config

import (
	"testing"
	"time"

	"github.com/baetyl/baetyl-go/v2/utils"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValue(t *testing.T) {
	expect := &CloudConfig{}
	expect.InitServer.Port = ":9003"
	expect.InitServer.WriteTimeout = time.Second * 30
	expect.InitServer.ReadTimeout = time.Second * 30
	expect.InitServer.ShutdownTime = time.Second * 3

	expect.AdminServer.Port = ":9004"
	expect.AdminServer.WriteTimeout = time.Second * 30
	expect.AdminServer.ReadTimeout = time.Second * 30
	expect.AdminServer.ShutdownTime = time.Second * 3

	expect.MisServer.Port = ":9006"
	expect.MisServer.WriteTimeout = time.Second * 30
	expect.MisServer.ReadTimeout = time.Second * 30
	expect.MisServer.ShutdownTime = time.Second * 3
	expect.MisServer.AuthToken = "baetyl-cloud-token"
	expect.MisServer.TokenHeader = "baetyl-cloud-token"
	expect.MisServer.UserHeader = "baetyl-cloud-user"

	expect.LogInfo.Level = "info"
	expect.LogInfo.MaxAge = 15
	expect.LogInfo.MaxSize = 50
	expect.LogInfo.MaxBackups = 15
	expect.LogInfo.Encoding = "json"

	expect.Plugin.PKI = "defaultpki"
	expect.Plugin.Auth = "defaultauth"
	expect.Plugin.License = "defaultlicense"
	expect.Plugin.Resource = "kube"
	expect.Plugin.Shadow = "database"
	expect.Plugin.Index = "database"
	expect.Plugin.Batch = "databaseext"
	expect.Plugin.Record = "databaseext"
	expect.Plugin.Callback = "databaseext"
	expect.Plugin.AppHistory = "database"
	expect.Plugin.Functions = []string{}
	expect.Plugin.Objects = []string{}
	expect.Plugin.Property = "database"
	expect.Plugin.Module = "database"
	expect.Plugin.SyncLinks = []string{"httplink"}
	expect.Plugin.Pubsub = "defaultpubsub"
	expect.Plugin.Locker = "defaultlocker"
	expect.Plugin.Task = "database"
	expect.Lock.ExpireTime = 5

	expect.Template.Path = "/etc/baetyl/templates"

	expect.Cache.ExpirationDuration = time.Minute * 10

	expect.Task.ScheduleTime = 30
	expect.Task.ConcurrentNum = 10
	expect.Task.QueueLength = 100
	expect.Task.LockExpiredTime = 60
	expect.Task.BatchNum = 100
	// case 0
	cfg := &CloudConfig{}
	err := utils.UnmarshalYAML(nil, cfg)
	assert.NoError(t, err)
	assert.EqualValues(t, expect, cfg)
}
