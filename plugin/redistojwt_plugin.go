package plugin

import (
	//"encoding/json"
	log "github.com/Sirupsen/logrus"
	//jwt "github.com/dgrijalva/jwt-go"
	"github.com/fzzy/radix/extra/pool"
	"github.com/gocraft/web"
	//"net/http"
	//"time"
  "errors"
)

const (
	PLUGIN_REDIS_TO_JWT string = "redis-jwt"
)

type RedisToJWTPlugin struct {
	connections *pool.Pool
}

func (p *RedisToJWTPlugin) Bootstrap() (Plugin, error) {
	var err error
	if p.connections == nil {
		p.connections, err = pool.NewPool("tcp", "127.0.0.1:6379", 10)
	}
	return p, err
}

func (p *RedisToJWTPlugin) Inbound(req *web.Request) (Plugin, error) {
	var err error
  token := req.Header.Get("Authorization")
  if token == "" {
    err = errors.New("Missing Authorization header")
  } else {
    redis, _ := connections.Get()
    user, er := redis.Cmd("get", token).Str()
    if er != nil {
      err = errors.New("Invalid token")
    } else {
      jwToken, er := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
          // Don't forget to validate the alg is what you expect:
          if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
              return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
          }
          return myLookupKey(token.Header["kid"])
      })

      if er == nil && token.Valid {
          deliverGoodness("!")
      } else {
          deliverUtterRejection(":(")
      }
    }
  }
  return p, err
}

type Session struct {
	PremiseID string `json="data.premise_id"`
	Expires   float64
	JWT       string
}

Looks up Authorization header token in redis.
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
				var session Session
				json.Unmarshal([]byte(user), &session)
				log.Warn(session.PremiseID)
				session.Expires = float64(time.Now().Unix() + 5) // Expires in 5 secs
				token := jwt.New(jwt.SigningMethodRS256)
				token.Claims = map[string]interface{}{"premise_id": session.PremiseID, "exp": session.Expires}
				tokenString, e := token.SignedString(c.Config.Key)
				if e != nil {
					log.Error("Could not sign JWT", session)
				}
				session.JWT = tokenString

				c.Session = &session
				log.Debug("RedisToJWTAuthMiddleware added ", c.Session)
				next(rw, req)
			}
		}
	} else {
		next(rw, req)
	}
}
