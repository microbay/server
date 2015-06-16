package server

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gocraft/web"
	"github.com/microbay/server/core"
	"github.com/microbay/server/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

type CompoundResponse struct {
	Response *http.Response
	Error    error
	Key      string
}

// Reverse proxies and load-balances backend micro services
func (c *Context) BalancedProxy(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	numRequests := len(c.Resource.Backends)

	// Proxy if single request
	if numRequests == 1 {
		for batchKey := range c.Resource.Backends {
			u := c.Resource.Backends[batchKey].Choose().String()
			for key := range c.Params {
				u = strings.Replace(u, ":"+key, c.Params[key], 1)
			}
			serverUrl, err := url.Parse(u)
			if err != nil {
				log.Error("URL failed to parse")
				c.RenderError(rw, errors.New("Internal Server Error"), "", http.StatusInternalServerError)
			}
			reverseProxy := proxy.New(serverUrl, &c.Resource.Middleware)
			res, err := reverseProxy.ServeHTTP(req.Request)
			if err != nil {
				c.RenderError(rw, err, "", http.StatusInternalServerError)
			} else {
				reverseProxy.CopyAndClose(rw, res)
			}
		}
		return
	}

	// Otherwise compound backend calls
	var wg sync.WaitGroup
	respones := make(chan *CompoundResponse, numRequests)
	for batchKey := range c.Resource.Backends {
		wg.Add(1)
		go func(batchKey string) {
			defer wg.Done()
			u := c.Resource.Backends[batchKey].Choose().String()
			for key := range c.Params {
				u = strings.Replace(u, ":"+key, c.Params[key], 1)
			}
			res, err := http.Get(u)
			respones <- &CompoundResponse{res, err, batchKey}
		}(batchKey)
	}

	wg.Wait()
	close(respones)

	// Create Compound response
	output := make(map[string]interface{})
	for composite := range respones {
		defer composite.Response.Body.Close()
		body, err := ioutil.ReadAll(composite.Response.Body)
		var data interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.RenderError(rw, errors.New("Internal Server Error"), "", http.StatusInternalServerError)
			return
		}
		output[composite.Key] = data
	}
	c.Render(rw, output, http.StatusOK)
}
