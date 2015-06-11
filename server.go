package server

import (
	"github.com/garyburd/redigo/redis"
	"github.com/microbay/server/backends"
	"github.com/microbay/server/plugin"
	//"github.com/fvbock/endless" ----> Hot reloads
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"github.com/spf13/viper"
	"net/http"
)

var Config API
var redisPool *redis.Pool

// Creates Root and resources routes and starts listening
func Start() {
	Config = LoadConfig()
	redisPool = connectRedis()
	defer redisPool.Close()
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
func connectRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", viper.GetString("redis_host"))
			if err != nil {
				log.Fatal(err.Error())
			}
			return c, err
		},
	}
}

func bootstrapPlugins(resources []*Resource) {
	for i := 0; i < len(resources); i++ {
		activePlugins := resources[i].Plugins
		plugins := make([]plugin.Interface, 0)
		for j := 0; j < len(activePlugins); j++ {
			n := activePlugins[j]
			if p, err := plugin.New(n); err != nil {
				log.Fatal(activePlugins[j], " plugin failed to bootstrap: ", err)
			} else {
				if rp, err := p.Bootstrap(&plugin.Config{redisPool}); err != nil {
					log.Fatal("Failed to bootstrap plugin ", err)
				} else {
					plugins = append(plugins, rp)
				}
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
