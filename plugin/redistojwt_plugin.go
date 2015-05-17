package plugin

import (
	//"encoding/json"
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fzzy/radix/extra/pool"
	"github.com/gocraft/web"
	"net/http"
	//"time"
	"fmt"
)

const (
	PLUGIN_REDIS_TO_JWT string = "redis-jwt"
)

type RedisToJWTPlugin struct {
	connections *pool.Pool
}

func (p *RedisToJWTPlugin) Inbound(req *web.Request) (Plugin, PluginError) {
	log.Debug("RedisToJWTPlugin::Inbound")
	var err PluginError
	token := req.Header.Get("Authorization")
	if token == "" {
		err = NewError(http.StatusUnauthorized, "Missing Authorization header")
	} else {
		redis, _ := p.connections.Get()
		_, er := redis.Cmd("get", token).Str()
		if er != nil {
			//err = errors.New("Invalid token")
		} else {
			jwToken, er := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
				}

				return "", nil
			})

			if er == nil && jwToken.Valid {
				log.Warn(">>>")
			} else {
				log.Warn("<<<")
			}
		}
	}
	return p, err
}

func (p *RedisToJWTPlugin) Bootstrap(config map[string]interface{}) (Plugin, error) {
	log.Debug("RedisToJWTPlugin::Bootstrap ", config)
	var err error
	if _, ok := config["host"]; ok != true {
		log.Fatal("RedisToJWTPlugin::Bootstrap failed to lookup 'host' key in config ", config)
	}
	if _, ok := config["key"]; ok != true {
		log.Fatal("RedisToJWTPlugin::Bootstrap failed to lookup 'key' key in config ", config)
	}
	if _, ok := config["connections"]; ok != true {
		log.Fatal("RedisToJWTPlugin::Bootstrap failed to lookup 'connections' key in config ", config)
	}
	if p.connections == nil {
		p.connections, err = pool.NewPool("tcp", config["host"].(string), int(config["connections"].(float64)))
		if err != nil {
			log.Fatal("RedisToJWTPlugin::Bootstrap failed to connect to Redis on ", config["host"].(string))
		}
	}
	return p, err
}
