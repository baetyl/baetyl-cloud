package localcache

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
)

func TestCache(t *testing.T) {
	conf := `
freeCacheConfig:
 maxBytes: 104857600

`
	filename := "cloud.yml"
	err := ioutil.WriteFile(filename, []byte(conf), 0644)
	assert.NoError(t, err)
	defer os.Remove(filename)
	common.SetConfFile(filename)

	p, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, p)

	cache := p.(plugin.DataCache)

	err = cache.SetString("a", "abc")
	assert.NoError(t, err)

	data, err := cache.GetString("a")
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
	str := ""
	for {
		str += "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		if len(str) > 65535 {
			break
		}
	}

	err = cache.SetByte("c", []byte(str))
	assert.NoError(t, err)

	dataByte, err := cache.GetByte("c")
	assert.NoError(t, err)
	assert.Equal(t, []byte(str), dataByte)

	err = cache.SetString("b", str)
	assert.NoError(t, err)

	data, err = cache.GetString("b")
	assert.NoError(t, err)
	assert.Equal(t, str, data)

	check, err = cache.Exist("b")
	assert.NoError(t, err)
	assert.Equal(t, true, check)

	err = cache.Delete("b")
	assert.NoError(t, err)
	check, err = cache.Exist("b")
	assert.NoError(t, err)
	assert.Equal(t, false, check)
}
