package redis

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func TestRedis(t *testing.T) {
	//redis test
	//	t.Skip(t.Name())
	conf := `
cacheRedisConfig:
 addr: 127.0.0.1:6379
 typeRedis: single
 maxIdle: 5
 maxActive: 5
 idleTimeout: 10
 password: 123456
`
	filename := "cloud.yml"
	err := ioutil.WriteFile(filename, []byte(conf), 0644)
	defer os.Remove(filename)
	common.SetConfFile(filename)

	p, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, p)

	cache := p.(plugin.DataCache)

	err = cache.Set("a", "abc")
	assert.NoError(t, err)

	data, err := cache.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, "abc", data)

	check, err := cache.Exist("a")
	assert.NoError(t, err)
	assert.Equal(t, true, check)

	err = cache.Delete("a")
	assert.NoError(t, err)
	check, err = cache.Exist("a")
	assert.NoError(t, err)
	assert.Equal(t, false, check)
}
