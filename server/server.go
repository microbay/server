package server

import (
	"github.com/SHMEDIALIMITED/apigo/model"
	//"github.com/fvbock/endless"
	//"container/ring"
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

	createMicroServiceRings(Config.Resources)

	err := http.ListenAndServe(viper.GetString("host"), rootRouter)
	if err != nil {
		log.Fatal("Failed to start server ", err)
	}

}

// Creates linked list (golang Ring) of micro services for round-robin lb.
func createMicroServiceRings(resources []*model.Resource) {
	for i := 0; i < len(resources); i++ {
		micros := resources[i].Micros
		resources[i].Index = 0
		flattenedMicros := make([]model.Micro, 0)

		for j := 0; j < len(micros); j++ {
			for n := 0; n < micros[j].Weight; n++ {
				flattenedMicros = append(flattenedMicros, micros[j])
			}
		}

		resources[i].BalancedMicros = flattenedMicros

		// RING Balancer
		// r := ring.New(len(flattenedMicros))
		// for q := 0; q < r.Len(); q++ {
		// 	r.Value = flattenedMicros[q].URL
		// 	log.Debug("ring ", q, " : ", r.Value)
		// 	r = r.Next()
		// }
		// resources[i].Ring = *r
		////////////////////////

	}
}
