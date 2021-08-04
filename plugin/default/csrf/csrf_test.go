package csrf

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/plugin"
	"github.com/stretchr/testify/assert"
)

const (
	confData = `
defaultcsrf:
  whitelist:
  - test.com`
)

func genConfig(workspace string) error {
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(workspace, "cloud.yml"), []byte(confData), 0755); err != nil {
		return err
	}
	return nil
}

func TestVerify(t *testing.T) {
	ctx := common.NewContextEmpty()
	ctx.Request = &http.Request{
		Header: http.Header{},
	}
	cookieName := "csrftoken"
	headerName := "csrftoken"

	err := genConfig("etc/baetyl")
	assert.NoError(t, err)
	defer os.RemoveAll(path.Dir("etc/baetyl"))

	plg, err := plugin.GetPlugin("defaultcsrf")
	assert.NoError(t, err)
	csrf, ok := plg.(*defaultCsrf)
	assert.True(t, ok)

	ctx.Request.Header.Set(headerName, "csrftoken-1")
	ctx.Request.AddCookie(&http.Cookie{Name: cookieName, Value: "csrftoken-1"})
	ctx.Request.Header.Set("Referer", "http://test.com/a/b")

	err = csrf.Verify(ctx)
	assert.NoError(t, err)

	csrf.cfg.CsrfConfig.Whitelist = []string{"test.cn", "test.com"}
	ctx.Request.Header.Set("Referer", "http://test.cn/a/b")
	err = csrf.Verify(ctx)
	assert.NoError(t, err)
	ctx.Request.Header.Set("Referer", "http://test.com/a/b")
	csrf.cfg.CsrfConfig.Whitelist = []string{"test.com"}

	csrf.cfg.CsrfConfig.Whitelist = []string{}
	err = csrf.Verify(ctx)
	assert.NoError(t, err)
	csrf.cfg.CsrfConfig.Whitelist = []string{"test.com"}

	ctx.Request.Header.Set(headerName, "csrftoken-2")
	err = csrf.Verify(ctx)
	assert.Error(t, err)
	ctx.Request.Header.Set(headerName, "csrftoken-1")

	ctx.Request.Header.Set("Referer", "")
	err = csrf.Verify(ctx)
	assert.Error(t, err)
	ctx.Request.Header.Set("Referer", "http://test.com/a/b")
}
