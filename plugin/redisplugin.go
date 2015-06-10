package plugin

import (
//log "github.com/Sirupsen/logrus"
//"github.com/gocraft/web"
//"net/http"
)

// type RedisPlugin struct {
// 	conns *pool.Pool
// }

// func (p *NoopPlugin) Bootstrap(config map[string]interface{}) (Interface, error) {
// 	if _, ok := config["host"]; ok != true {
// 		log.Fatal("RedisToJWTPlugin::Bootstrap failed to lookup 'host' key in config ", config)
// 	}
// 	if p.conns == nil {
// 		p.conns, err = pool.NewPool("tcp", config["host"].(string), int(config["connections"].(float64)))
// 		if err != nil {
// 			log.Fatal("RedisToJWTPlugin::Bootstrap failed to connect to Redis on ", config["host"].(string))
// 		}
// 	}
// 	var err error
// 	return p, err
// }

// func (p *NoopPlugin) Inbound(req *web.Request) (int, error) {
// 	log.Warn("NoopPlugin::Inbound")
// 	var err error
// 	return http.StatusOK, err
// }

// func (p *NoopPlugin) Outbound(res *http.Response) (int, error) {
// 	log.Warn("NoopPlugin::Outbound")
// 	var err error
// 	return http.StatusOK, err
// }
