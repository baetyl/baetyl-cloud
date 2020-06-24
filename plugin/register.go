package plugin

import (
	"github.com/baetyl/baetyl-go/log"
	"io"
	"strings"
	"sync"

	"github.com/baetyl/baetyl-cloud/common"
)

// Plugin interfaces
type Plugin interface {
	io.Closer
}

// Factory create engine by given config
type Factory func() (Plugin, error)

// PluginFactory contains all supported plugin factory
var pluginFactory = make(map[string]Factory)
var plugins = map[string]Plugin{}
var mu sync.Mutex

// RegisterFactory adds a supported plugin
func RegisterFactory(name string, f Factory) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := pluginFactory[name]; ok {
		log.L().Info("plugin already exists, skip", log.Any("plugin", name))
		return
	}
	pluginFactory[name] = f
	log.L().Info("plugin is registered", log.Any("plugin", name))
}

// GetPlugin GetPlugin
func GetPlugin(name string) (Plugin, error) {
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	if p, ok := plugins[name]; ok {
		return p, nil
	}
	f, ok := pluginFactory[name]
	if !ok {
		return nil, common.Error(common.ErrPluginNotFound, common.Field("name", name))
	}
	p, err := f()
	if err != nil {
		log.L().Error("plugin create failed", log.Error(err))
		return nil, err
	}
	plugins[name] = p
	return p, nil
}

// ClosePlugins ClosePlugins
func ClosePlugins() {
	mu.Lock()
	defer mu.Unlock()
	for _, v := range plugins {
		v.Close()
	}
}
