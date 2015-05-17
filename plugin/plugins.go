package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
)

type Plugin interface {
	Bootstrap(config map[string]interface{}) (Plugin, error)
	Inbound(req *web.Request) (Plugin, PluginError)
}

var pluginPool map[string]Plugin

func init() {
	pluginPool = make(map[string]Plugin)
	pluginPool[PLUGIN_REDIS_TO_JWT] = &RedisToJWTPlugin{}
	pluginPool[PLUGIN_NOOP] = &NoopPlugin{}
}

func Get(id string) Plugin {
	if _, ok := pluginPool[id]; ok != true {
		log.Fatal("Could not retrieve ", id, " plugin from plugin pool.\n", "Plugins available: ", pluginPool)
	}
	return pluginPool[id]
}
