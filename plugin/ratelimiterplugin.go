package plugin

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/fzzy/radix/extra/pool"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
	"math"
	"net/http"
	"time"
)

const (
	PLUGIN_RATELIMITER string = "ratelimiter"
)

// var now = microtime.now();
// var key = namespace + id;
// var clearBefore = now - interval;

// var batch = redis.multi();
// batch.zremrangebyscore(key, 0, clearBefore);
// batch.zrange(key, 0, -1);
// batch.zadd(key, now, now);
// batch.expire(key, Math.ceil(interval/1000000)); // convert to seconds, as used by redis ttl.
// batch.exec(function (err, resultArr) {

//         if (err) return cb(err);

//         var userSet = zrangeToUserSet(resultArr[1]);

//         var tooManyInInterval = userSet.length >= maxInInterval;
//         var timeSinceLastRequest = minDifference && (now - userSet[userSet.length - 1]);

//         var result;
//         if (tooManyInInterval || timeSinceLastRequest < minDifference) {
//           result = Math.min(userSet[0] - now + interval, minDifference ? minDifference - timeSinceLastRequest : Infinity);
//           result = Math.floor(result/1000); // convert to miliseconds for user readability.
//         } else {
//           result = 0;
//         }

//         return cb(null, result)

type RateLimiterPlugin struct {
	redisPool *redis.Pool
}

func (p *RateLimiterPlugin) Bootstrap(config *Config) (Interface, error) {
	p.redisPool = config.RedisPool
	var err error
	return p, err
}

func (p *RateLimiterPlugin) Inbound(req *web.Request) (int, error) {
	log.Info("RateLimiterPlugin::Inbound")

	interval := 10000
	now := time.Now().UnixNano()
	clearBefore := now - int64(interval*1000)
	key := "Patrick"
	//time.Since(startTime).Nanoseconds()

	c := p.redisPool.Get()
	c.Send("MULTI")
	c.Send("zremrangebyscore", key, 0, clearBefore)
	c.Send("zrange", key, 0, -1)
	c.Send("zadd", key, now, now)
	c.Send("expire", key, math.Ceil(float64(interval/1000000)))
	r, err := c.Do("EXEC")
	log.Error(r)

	return http.StatusOK, err
}

func (p *RateLimiterPlugin) Outbound(res *http.Response) (int, error) {
	log.Warn("NoopPlugin::Outbound")
	var err error
	return http.StatusOK, err
}
