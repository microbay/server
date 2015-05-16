package server

import (
	"github.com/SHMEDIALIMITED/apigo/model"
	"github.com/SHMEDIALIMITED/apigo/plugin"
	//"github.com/fvbock/endless" ----> Hot reloads
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"github.com/spf13/viper"
	"net/http"
)

var Config model.API

// Creates Root and resources routes and starts listening
func Start() {
	Config = model.Load()

	model.PrepareLoadBalancer(Config.Resources)

	plugin.Bootstrap(Config.Resources)
	// if err := ConnectRedis(); err != nil {
	// 	log.Fatal("Failed to connect to Redis ", err)
	// }

	rootRouter := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
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
