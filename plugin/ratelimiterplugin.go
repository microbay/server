package plugin

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	//"github.com/fzzy/radix/extra/pool"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
	"github.com/microbay/server/core"
	"math"
	"net/http"
	//"reflect"
	"time"
)

/**
  Rate limiter Plugin based on rolling window.
  Thanks to https://engineering.classdojo.com/blog/2015/02/06/rolling-rate-limiter/

  Exmaple config:
  {
    "id" : "ratelimiter",
    "interval" : 10,
    "max_req_per_interval" : 5
  }
*/

const (
	PLUGIN_RATELIMITER          string = "ratelimiter"
	PLUGIN_RATELIMITER_EXCEEDED string = "Rate limit exceeded"
)

type RateLimiterPlugin struct {
	*core.Renderer
	redisPool      *redis.Pool
	interval       int64
	reqPerInterval int
}

func (p *RateLimiterPlugin) Bootstrap(config *Config) (Interface, error) {
	var err error
	p.redisPool = config.RedisPool
	if _, ok := config.Properties["interval"]; ok != true {
		return p, errors.New("RateLimiterPlugin needs 'interval' int")
	}
	if _, ok := config.Properties["max_req_per_interval"]; ok != true {
		return p, errors.New("RateLimiterPlugin needs 'max_req_per_interval' int")
	}
	p.interval = int64(config.Properties["interval"].(float64)) * 1000
	p.reqPerInterval = int(config.Properties["max_req_per_interval"].(float64))

	return p, err
}

func (p *RateLimiterPlugin) Inbound(rw web.ResponseWriter, req *web.Request) (int, error) {
	log.Info("RateLimiterPlugin::Inbound")

	reqPerInterval := p.reqPerInterval
	interval := p.interval * 1000
	now := time.Now().UnixNano() / 1000
	clearBefore := now - interval
	key := "Patrick"
	//time.Since(startTime).Nanoseconds()

	c := p.redisPool.Get()
	c.Send("MULTI")
	c.Send("zremrangebyscore", key, 0, clearBefore)
	c.Send("zrange", key, 0, -1)
	c.Send("zadd", key, now, now)
	c.Send("expire", key, math.Ceil(float64(interval/1000000)))
	r, err := c.Do("EXEC")
	if err != nil {
		p.RenderError(rw, err, "", 500)
		return http.StatusInternalServerError, err
	}
	userSet := r.([]interface{})[1].([]interface{})

	l := len(userSet)

	reached := l >= reqPerInterval

	// Todo process timeleft
	if reached {
		p.RenderError(rw, errors.New(PLUGIN_RATELIMITER_EXCEEDED), "", 429)
		return 429, errors.New(PLUGIN_RATELIMITER_EXCEEDED)
	} else {
		return http.StatusOK, nil
	}
}

func (p *RateLimiterPlugin) Outbound(res *http.Response) (int, error) {
	return http.StatusOK, nil
}
