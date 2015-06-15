package server

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
	"github.com/microbay/server/core"
	"github.com/microbay/server/proxy"
	"net/http"
	"net/url"
	"strings"
)

// Root context
type Context struct {
	*core.Renderer
	Config   API
	Resource *Resource
	Redis    redis.Conn
	Params   core.URLParams
}

// Assigns global config to context --> muset be a better way to pass that onto context
func (c *Context) ConfigMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Config = Config
	next(rw, req)
}

// Redis Middleware
func (c *Context) RedisMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Redis = redisPool.Get()
	defer c.Redis.Close()
	next(rw, req)
}

// 403 on API root
func (c *Context) RootMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if req.URL.Path == "/" {
		c.RenderError(rw, errors.New(c.Config.Name+" root access forbidden"), "", http.StatusForbidden)
	} else {
		next(rw, req)
	}
}

func (c *Context) ResourceConfigMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	var err error
	c.Resource, err = c.Config.FindResourceByRequest(req.Request)
	if err != nil {
		if err.Error() == "Method" {
			rw.Header().Set("Allow", strings.Join(c.Resource.Methods, ", "))
			c.RenderError(rw, errors.New("Method Not Allowed"), "", http.StatusMethodNotAllowed)
		} else {
			c.RenderError(rw, errors.New("Access Forbidden"), "", http.StatusForbidden)
		}

	} else {
		c.Params = core.Params(req.URL.Path, c.Resource.Regex, c.Resource.Keys)
		next(rw, req)
	}
}

func (c *Context) PluginMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	for i := range c.Resource.Plugins {
		if _, err := c.Resource.Middleware[i].Inbound(rw, req); err != nil {
			log.Warn(err.Error())
			return
		}
	}
	next(rw, req)
}

// Reverse proxies and load-balances backend micro services
func (c *Context) BalancedProxy(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	backend := c.Resource.Backends.Choose()
	if backend == nil {
		log.Error("no backend for client %s", req.RemoteAddr)
	}

	u := backend.String()
	for key := range c.Params {
		u = strings.Replace(u, ":"+key, c.Params[key], 1)
	}

	serverUrl, err := url.Parse(u)
	if err != nil {
		log.Error("URL failed to parse")
	}

	reverseProxy := proxy.New(serverUrl, &c.Resource.Middleware)
	reverseProxy.ServeHTTP(rw, req.Request)
}
