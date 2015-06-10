package plugin

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fzzy/radix/extra/pool"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const (
	PLUGIN_REDIS_TO_JWT             string = "redis-jwt"
	PLUGIN_REDIS_TO_JWT_HEADER      string = "Authorization"
	PLUGIN_REDIS_TO_JWT_MSG_MISSING string = "Missing Authorization header"
	PLUGIN_REDIS_TO_JWT_MSG_INVALID string = "Invalid token"
)

type RedisToJWTPlugin struct {
	key         []byte
	keyFunc     jwt.Keyfunc
	connections *pool.Pool
}

func (p *RedisToJWTPlugin) Inbound(req *web.Request) (int, error) {
	var err error
	token := req.Header.Get(PLUGIN_REDIS_TO_JWT_HEADER)
	if token == "" {
		return http.StatusUnauthorized, errors.New(PLUGIN_REDIS_TO_JWT_MSG_MISSING)
	} else {
		redis, _ := p.connections.Get()
		_, er := redis.Cmd("get", token).Str()
		if er != nil {
			return http.StatusUnauthorized, errors.New(PLUGIN_REDIS_TO_JWT_MSG_INVALID)
		} else {
			jwToken, er := jwt.Parse(token, p.keyFunc)

			if er != nil || !jwToken.Valid {
				return http.StatusUnauthorized, errors.New(PLUGIN_REDIS_TO_JWT_MSG_INVALID)
			}
		}
	}
	return http.StatusOK, err
}

func (p *RedisToJWTPlugin) Outbound(res *http.Response) (int, error) {
	log.Warn("RedisToJWTPlugin::Outbound")
	var err error
	return http.StatusOK, err
}

func (p *RedisToJWTPlugin) Bootstrap(config map[string]interface{}) (Interface, error) {
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
	keyPath, _ := config["key"].(string)
	keyAbsPath, _ := filepath.Abs(keyPath)
	key, err := ioutil.ReadFile(keyAbsPath)
	if err != nil {
		log.Fatal("RedisToJWTPlugin::Bootstrap failed loading public key file", keyPath, err)
	}
	p.key = key
	p.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return p.key, nil
	}
	if p.connections == nil {
		p.connections, err = pool.NewPool("tcp", config["host"].(string), int(config["connections"].(float64)))
		if err != nil {
			log.Fatal("RedisToJWTPlugin::Bootstrap failed to connect to Redis on ", config["host"].(string))
		}
	}
	return p, err
}
