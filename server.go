package server

import (
	"github.com/microbay/server/backends"
	"github.com/microbay/server/plugin"
	//"github.com/fvbock/endless" ----> Hot reloads
	log "github.com/Sirupsen/logrus"
	"github.com/fzzy/radix/extra/pool"
	"github.com/gocraft/web"
	"github.com/spf13/viper"
	"net/http"
)

var Config API
var redisPool *pool.Pool

// Creates Root and resources routes and starts listening
func Start() {
	Config = LoadConfig()
	connectRedis()
	bootstrapLoadBalancer(Config.Resources)
	bootstrapPlugins(Config.Resources)
	rootRouter := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).RedisMiddleware).
		Middleware((*Context).ConfigMiddleware).
		Middleware((*Context).RootMiddleware).
		Middleware((*Context).ResourceConfigMiddleware).
		Middleware((*Context).PluginMiddleware).
		Middleware((*Context).BalancedProxy)
	log.Info(Config.Name, " listening on ", viper.GetString("host"), " in ", viper.Get("env"), " mode")
	err := http.ListenAndServe(viper.GetString("host"), rootRouter)
	if err != nil {
		log.Fatal("Failed to start server ", err)
	}
}

// Creates redis connection pool from config
func connectRedis() {
	config := viper.GetStringMap("redis")
	var err error
	if _, ok := config["host"]; ok != true {
		log.Fatal("Redis::connectRedis - failed to lookup Redis 'host' key in config ", config)
	}
	if _, ok := config["idle_connections"]; ok != true {
		log.Fatal("Redis::connectRedis - failed to lookup to lookup Redis 'idle_connections' key in config ", config)
	}
	redisPool, err = pool.NewPool("tcp", config["host"].(string), int(config["idle_connections"].(float64)))
	if err != nil {
		log.Fatal("Server::connectRedis - failed to connect to Redis on ", config["host"].(string))
	}

}

func bootstrapPlugins(resources []*Resource) {
	for i := 0; i < len(resources); i++ {
		activePlugins := resources[i].Plugins
		plugins := make([]plugin.Interface, 0)
		for j := 0; j < len(activePlugins); j++ {

			n := activePlugins[j]
			if _, err := plugin.New(n); err != nil {
				log.Fatal(activePlugins[j], " plugin failed to bootstrap: ", err)
			} else {
				plugins = append(plugins, &plugin.NoopPlugin{})
			}
		}
		resources[i].Middleware = plugins
	}
}

// Creates linked list (golang Ring) from weighted micros array per resource
func bootstrapLoadBalancer(resources []*Resource) {
	for i := 0; i < len(resources); i++ {
		micros := resources[i].Micros
		flattenedMicros := make([]string, 0)
		for j := 0; j < len(micros); j++ {
			for n := 0; n < micros[j].Weight; n++ {
				flattenedMicros = append(flattenedMicros, micros[j].URL)
			}
		}
		resources[i].Backends = backends.Build("round-robin", flattenedMicros)
	}
}
