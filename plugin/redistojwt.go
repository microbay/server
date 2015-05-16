package plugin

import (
	//"encoding/json"
	log "github.com/Sirupsen/logrus"
	//jwt "github.com/dgrijalva/jwt-go"
	"github.com/fzzy/radix/extra/pool"
	"github.com/gocraft/web"
	//"net/http"
	//"time"
)

type RedisToJWTPlugin struct{}

func (p *RedisToJWTPlugin) Inbound(req *web.Request) {
	log.Debug("RedisToJWTPlugin")
}

type Session struct {
	PremiseID string `json="data.premise_id"`
	Expires   float64
	JWT       string
}

var connections *pool.Pool

func ConnectRedis() error {
	var err error
	connections, err = pool.NewPool("tcp", "127.0.0.1:6379", 10)
	return err
}

// Looks up Authorization header token in redis.
// func (c *Context) RedisToJWTAuthMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
// 	if c.Resource.Auth == REDIS_JWT {
// 		token := req.Header.Get("Authorization")
// 		if token == "" {
// 			c.RenderError(rw, "Missing Authorization header", http.StatusUnauthorized)
// 		} else {
// 			redis, _ := connections.Get()
// 			user, err := redis.Cmd("get", token).Str()
// 			if err != nil {
// 				c.RenderError(rw, "Invalid token", http.StatusUnauthorized)
// 			} else {
// 				var session Session
// 				json.Unmarshal([]byte(user), &session)
// 				log.Warn(session.PremiseID)
// 				session.Expires = float64(time.Now().Unix() + 5) // Expires in 5 secs
// 				token := jwt.New(jwt.SigningMethodRS256)
// 				token.Claims = map[string]interface{}{"premise_id": session.PremiseID, "exp": session.Expires}
// 				tokenString, e := token.SignedString(c.Config.Key)
// 				if e != nil {
// 					log.Error("Could not sign JWT", session)
// 				}
// 				session.JWT = tokenString

// 				c.Session = &session
// 				log.Debug("RedisToJWTAuthMiddleware added ", c.Session)
// 				next(rw, req)
// 			}
// 		}
// 	} else {
// 		next(rw, req)
// 	}
// }
