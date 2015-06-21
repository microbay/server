package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
	"net/http"
	"reflect"
	//"strings"
)

type Config struct {
	RedisPool  *redis.Pool
	Properties map[string]interface{}
}

type Interface interface {
	Bootstrap(config *Config) (Interface, error)
	Inbound(rw web.ResponseWriter, req *web.Request) (int, error)
	Outbound(res *http.Response) (int, error)
}

var pluginRegistry map[string]reflect.Type

func init() {
	pluginRegistry = make(map[string]reflect.Type)
	pluginRegistry[PLUGIN_AUTH_JWT] = reflect.TypeOf(AuthJWTPlugin{})
	pluginRegistry[PLUGIN_RATELIMITER] = reflect.TypeOf(RateLimiterPlugin{})
	pluginRegistry[PLUGIN_NOOP] = reflect.TypeOf(NoopPlugin{})
	log.Debug(log.WithFields(log.Fields{
		"type":    "pluginregistry",
		"content": pluginRegistry,
	}))
}

func New(id string) (Interface, error) {
	if _, ok := pluginRegistry[id]; ok != true {
		log.Fatal("Could not retrieve ", id, " plugin from plugin registry.\n", "Plugins available: ", pluginRegistry)
	}
	return reflect.New(pluginRegistry[id]).Interface().(Interface), nil
}
