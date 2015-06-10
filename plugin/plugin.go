package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"net/http"
	"reflect"
	//"strings"
)

type Interface interface {
	Bootstrap(config map[string]interface{}) (Interface, error)
	Inbound(req *web.Request) (int, error)
	Outbound(res *http.Response) (int, error)
}

var pluginRegistry map[string]reflect.Type

func init() {
	pluginRegistry = make(map[string]reflect.Type)
	pluginRegistry[PLUGIN_REDIS_TO_JWT] = reflect.TypeOf(RedisToJWTPlugin{})
	pluginRegistry[PLUGIN_REDIS_RATELIMITER] = reflect.TypeOf(RedisRateLimiterPlugin{})
	pluginRegistry[PLUGIN_NOOP] = reflect.TypeOf(NoopPlugin{})
	pluginRegistry[PLUGIN_TRANSFORMER] = reflect.TypeOf(TransformerPlugin{})
}

func New(id string) (Interface, error) {

	if _, ok := pluginRegistry[id]; ok != true {
		log.Fatal("Could not retrieve ", id, " plugin from plugin registry.\n", "Plugins available: ", pluginRegistry)
	}

	// one way is to have a value of the type you want already
	//a := 1
	// reflect.New works kind of like the built-in function new
	// We'll get a reflected pointer to a new int value

	intPtr := reflect.New(pluginRegistry[id]).Interface().(Interface)
	// Just to prove it

	// Prints 0
	log.Info(intPtr)

	return intPtr, nil

	// if strings.Index(id, "redis") != -1 {
	//   return pluginPool[id].Bootstrap(Config.plugins["redis-jwt"])
	// } else {
	//   return pluginPool[id].Bootstrap(Config.plugins["redis-jwt"])
	// }

}
