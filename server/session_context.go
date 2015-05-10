package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fzzy/radix/extra/pool"
	"github.com/gocraft/web"
	"net/http"
)

const (
	REDIS_JWT string = "redis-jwt"
)

var connections *pool.Pool

func ConnectRedis() error {
	var err error
	connections, err = pool.NewPool("tcp", "127.0.0.1:6379", 10)
	return err
}

// Looks up Authorization header token in redis.
func (c *Context) RedisToJWTAuthMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if c.Resource.Auth == REDIS_JWT {
		token := req.Header.Get("Authorization")
		if token == "" {
			c.RenderError(rw, "Missing Authorization header", http.StatusUnauthorized)
		} else {
			redis, _ := connections.Get()
			user, err := redis.Cmd("get", token).Str()
			if err != nil {
				c.RenderError(rw, "Invalid token", http.StatusUnauthorized)
			} else {
				c.Session = &Session{user}
				log.Debug("RedisToJWTAuthMiddleware added ", c.Session)
				next(rw, req)
			}
		}
	} else {
		next(rw, req)
	}
}
