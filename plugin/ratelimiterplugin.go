package plugin

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/fzzy/radix/extra/pool"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
	"net/http"
)

const (
	PLUGIN_RATELIMITER string = "ratelimiter"
)

type RateLimiterPlugin struct {
	redisPool *redis.Pool
}

func (p *RateLimiterPlugin) Bootstrap(config *Config) (Interface, error) {
	p.redisPool = config.RedisPool
	var err error
	return p, err
}

func (p *RateLimiterPlugin) Inbound(req *web.Request) (int, error) {
	log.Warn("RateLimiterPlugin::Inbound")
	c := p.redisPool.Get()
	c.Send("MULTI")
	c.Send("INCR", "foo")
	c.Send("INCR", "bar")
	r, err := c.Do("EXEC")
	log.Error(r) // prints [1, 1]

	return http.StatusOK, err
}

func (p *RateLimiterPlugin) Outbound(res *http.Response) (int, error) {
	log.Warn("NoopPlugin::Outbound")
	var err error
	return http.StatusOK, err
}
