package server

import (
	"github.com/SHMEDIALIMITED/apigo/model"
	"github.com/SHMEDIALIMITED/apigo/server/backends"
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

	log.Info(Config.Name, " listening on ", viper.GetString("host"), " in ", viper.Get("env"), " mode")

	rootRouter := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).ConfigMiddleware).
		Middleware((*Context).RootMiddleware)

	rootRouter.Subrouter(ResourceContext{}, "/").
		Middleware((*ResourceContext).ResourceConfigMiddleware).
		Middleware((*ResourceContext).RedisToJWTAuthMiddleware).
		Get("/:resource", (*ResourceContext).BalancedProxy)

	prepareLoadBalancer(Config.Resources)

	err := http.ListenAndServe(viper.GetString("host"), rootRouter)
	if err != nil {
		log.Fatal("Failed to start server ", err)
	}
}

// Creates linked list (golang Ring) from weighted micros array per resource
func prepareLoadBalancer(resources []*model.Resource) {
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
