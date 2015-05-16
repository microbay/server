package plugin

import (
	"github.com/SHMEDIALIMITED/apigo/model"
	log "github.com/Sirupsen/logrus"
)

const (
	PLUGIN_REDIS_JWT string = "redis-jwt"
	PLUGIN_NOOP      string = "noop"
)

var pluginPool map[string]model.Plugin

func init() {
	pluginPool = make(map[string]model.Plugin)
	pluginPool[PLUGIN_REDIS_JWT] = &RedisToJWTPlugin{}
	pluginPool[PLUGIN_NOOP] = &NoopPlugin{}
}

func getPlugin(id string) model.Plugin {
	if _, ok := pluginPool[id]; ok != true {
		log.Fatal("Could not retrieve ", id, " plugin from plugin pool.\n", "Plugins available: ", pluginPool)
	}
	return pluginPool[id]
}

func Bootstrap(resources []*model.Resource) {
	for i := 0; i < len(resources); i++ {
		activePlugins := resources[i].Plugins
		plugins := make([]model.Plugin, 0)
		for j := 0; j < len(activePlugins); j++ {
			plugins = append(plugins, getPlugin(activePlugins[j]))
		}
		resources[i].Middleware = plugins
		log.Warn(">>>>", resources[i].Plugins)
		log.Error(">>>>", resources[i].Plugins)
	}
}
