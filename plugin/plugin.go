package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
)

type Interface interface {
	Bootstrap(config map[string]interface{}) (Interface, error)
	Inbound(req *web.Request) (int, error)
	Outbound(rw web.ResponseWriter, req *web.Request) (int, error)
}

var pluginPool map[string]Interface

func init() {
	pluginPool = make(map[string]Interface)
	pluginPool[PLUGIN_REDIS_TO_JWT] = &RedisToJWTPlugin{}
	pluginPool[PLUGIN_NOOP] = &NoopPlugin{}
}

func Get(id string) Interface {
	if _, ok := pluginPool[id]; ok != true {
		log.Fatal("Could not retrieve ", id, " plugin from plugin pool.\n", "Plugins available: ", pluginPool)
	}
	return pluginPool[id]
}
