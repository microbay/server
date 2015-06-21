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
	"strings"
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
	PLUGIN_RATELIMITER                            string = "ratelimiter"
	PLUGIN_RATELIMITER_EXCEEDED                   string = "Rate limit exceeded"
	PLUGIN_RATELIMITER_HEADER_AUTH                string = "Authorization"
	PLUGIN_RATELIMITER_HEADER_CONSUMER            string = "X-Consumer-Id"
	PLUGIN_RATELIMITER_HEADER_ISSUER              string = "X-Issuer-Id"
	PLUGIN_RATELIMITER_ERROR_HEADER_MISSING       string = "Rate Limiter Plugin needs X-Issuer-Id and X-Consumer-Id headers"
	PLUGIN_RATELIMITER_REDIS_CMD_MULTI            string = "MULTI"
	PLUGIN_RATELIMITER_REDIS_CMD_ZREMRANGEBYSCORE string = "ZREMRANGEBYSCORE"
	PLUGIN_RATELIMITER_REDIS_CMD_ZRANGE           string = "ZRANGE"
	PLUGIN_RATELIMITER_REDIS_CMD_ZADD             string = "ZADD"
	PLUGIN_RATELIMITER_REDIS_CMD_EXPIRE           string = "EXPIRE"
	PLUGIN_RATELIMITER_REDIS_CMD_EXEC             string = "EXEC"
)

type RateLimiterPlugin struct {
	*core.Renderer
	redisPool      *redis.Pool
	interval       int64
	reqPerInterval int
	key            string
}

func (p *RateLimiterPlugin) Bootstrap(config *Config) (Interface, error) {
	var err error
	p.redisPool = config.RedisPool
	if _, ok := config.Properties["interval"]; ok != true {
		return p, errors.New("RateLimiterPlugin needs 'interval' int")
	}
	if _, ok := config.Properties["path"]; ok != true {
		log.Fatal("RateLimiterPlugin needs failed to lookup path for ", config)
	}
	if _, ok := config.Properties["max_req_per_interval"]; ok != true {
		return p, errors.New("RateLimiterPlugin needs 'max_req_per_interval' int")
	}
	p.interval = int64(config.Properties["interval"].(float64)) * 1000
	p.reqPerInterval = int(config.Properties["max_req_per_interval"].(float64))
	p.key = config.Properties["path"].(string)
	p.key = strings.Replace(p.key, ":", "=", 10)
	return p, err
}

func (p *RateLimiterPlugin) Inbound(rw web.ResponseWriter, req *web.Request) (int, error) {
	reqPerInterval := p.reqPerInterval
	interval := p.interval * 1000
	now := time.Now().UnixNano() / 1000
	clearBefore := now - interval
	consumer := req.Header.Get(PLUGIN_RATELIMITER_HEADER_CONSUMER)
	issuer := req.Header.Get(PLUGIN_RATELIMITER_HEADER_ISSUER)
	id := issuer + "." + consumer
	if id == "." {
		err := errors.New(PLUGIN_RATELIMITER_ERROR_HEADER_MISSING)

		log.WithFields(log.Fields{
			"type": PLUGIN_RATELIMITER + "." + PLUGIN_RATELIMITER_ERROR_HEADER_MISSING,
		}).Info(err.Error())

		p.RenderError(rw, err, "", 500)
		return http.StatusInternalServerError, err
	}
	key := id + ":" + p.key
	//minDifference := 100
	//time.Since(startTime).Nanoseconds()

	c := p.redisPool.Get()
	c.Send(PLUGIN_RATELIMITER_REDIS_CMD_MULTI)
	c.Send(PLUGIN_RATELIMITER_REDIS_CMD_ZREMRANGEBYSCORE, key, 0, clearBefore)
	c.Send(PLUGIN_RATELIMITER_REDIS_CMD_ZRANGE, key, 0, -1)
	c.Send(PLUGIN_RATELIMITER_REDIS_CMD_ZADD, key, now, now)
	c.Send(PLUGIN_RATELIMITER_REDIS_CMD_EXPIRE, key, math.Ceil(float64(interval/1000000)))
	r, err := c.Do(PLUGIN_RATELIMITER_REDIS_CMD_EXEC)

	if err != nil {
		p.RenderError(rw, err, "", 500)
		return http.StatusInternalServerError, err
	}
	userSet := r.([]interface{})[1].([]interface{})

	l := len(userSet)

	reached := l >= reqPerInterval
	//log.Error(reflect.TypeOf(userSet[len(userSet)-1]))
	//timeSinceLastRequest := minDifference && now-userSet[len(userSet)-1]

	// Todo process timeleft
	//if reached || timeSinceLastRequest < minDifference {
	if reached {
		log.WithFields(log.Fields{
			"type":     PLUGIN_RATELIMITER,
			"issuer":   req.Header.Get(PLUGIN_RATELIMITER_HEADER_ISSUER),
			"consumer": req.Header.Get(PLUGIN_RATELIMITER_HEADER_CONSUMER),
			"path":     p.key,
		}).Info(PLUGIN_RATELIMITER_EXCEEDED)
		p.RenderError(rw, errors.New(PLUGIN_RATELIMITER_EXCEEDED), "", 429)
		//result := math.Min(userSet[0]-now+interval, minDifference-timeSinceLastRequest)
		//rw.Header().Set("Retry-After", string(result/1000))
		return 429, errors.New(PLUGIN_RATELIMITER_EXCEEDED)
	} else {
		return http.StatusOK, nil
	}
}

func (p *RateLimiterPlugin) Outbound(res *http.Response) (int, error) {
	return http.StatusOK, nil
}
