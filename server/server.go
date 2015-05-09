package server

import (
	"github.com/SHMEDIALIMITED/apigo/model"
	//"github.com/fvbock/endless"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"github.com/spf13/viper"
	"net/http"
)

var Config model.API

func Start() {

	Config = model.Load()

	log.Info(Config.Name, " listening on ", viper.GetString("host"), " in ", viper.Get("env"), " mode")

	rootRouter := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).ConfigMiddleware).
		Middleware((*Context).RootMiddleware)

	rootRouter.Subrouter(ResourceContext{}, "/").
		Middleware((*ResourceContext).ResourceConfigMiddleware).
		Middleware((*ResourceContext).RedisToJWTAuthMiddleware).
		Get("/:resource", (*ResourceContext).Proxy)

	err := http.ListenAndServe(viper.GetString("host"), rootRouter)
	if err != nil {
		log.Fatal("Failed to start server ", err)
	}
}
