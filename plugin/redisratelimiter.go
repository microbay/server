package plugin

import (
//log "github.com/Sirupsen/logrus"
//"github.com/fzzy/radix/extra/pool"
///"github.com/gocraft/web"
//"net/http"
)

const (
	PLUGIN_REDIS_RATELIMITER string = "redis-ratelimiter"
)

type RedisRateLimiterPlugin struct{}

// func (p *NoopPlugin) Bootstrap(config map[string]interface{}) (Interface, error) {
// 	log.Warn("NoopPlugin::Bootstrap")
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
