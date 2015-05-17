package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Root context
type Context struct {
	Config   API
	Resource *Resource
}

// Assigns global config to context --> muset be a better way to pass that onto context
func (c *Context) ConfigMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Config = Config
	next(rw, req)
}

// 403 on API root
func (c *Context) RootMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if req.URL.Path == "/" {
		c.RenderError(rw, c.Config.Name+" root access forbidden", http.StatusForbidden)
	} else {
		next(rw, req)
	}
}

func (c *Context) ResourceConfigMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	var err error
	c.Resource, err = c.Config.FindResourceByRequest(req.Request)
	if err != nil {
		c.RenderError(rw, "Access Forbidden", http.StatusForbidden)
	} else {
		next(rw, req)
	}
}

func (c *Context) PluginMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	for i := range c.Resource.Plugins {
		if _, err := c.Resource.Middleware[i].Inbound(req); err != nil {
			message, status := err.Error()
			c.RenderError(rw, message, status)
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

	serverUrl, err := url.Parse(backend.String())
	if err != nil {
		log.Fatal("URL failed to parse")
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)
	//if c.Resource.Auth == REDIS_JWT {
	//combinedHeaders := headerCombiner(reverseProxy, c.Session.JWT)
	//}

	//log.Debug(">>>", c.Session.JWT)

	// if c.Resource.Auth == REDIS_JWT {
	//  c.RenderError(rw, "Invalid token", http.StatusUnauthorized)
	// } else {
	//  next(rw, req)
	// }

	req.URL.Path = ""
	reverseProxy.ServeHTTP(rw, req.Request)
}

// Append additional query params to the original URL query.
func headerCombiner(handler http.Handler, token string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", token)
		r.Header.Set("X-Premise", "-2287340764")
		handler.ServeHTTP(w, r)
	})
}

// Responds from error to JSON
type JSONError struct {
	Message  string `json:"message"`
	MoreInfo string `json:"moreInfo"`
}

// Renders from struct to JSON
func (a *Context) Render(rw web.ResponseWriter, model interface{}, status int) {
	jsonString, err := json.MarshalIndent(&model, "", "    ")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		log.Error("Render failed to marshall model", err)
		return
	}
	header := rw.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	header.Set("Access-Control-Allow-Headers", "Authorization")
	rw.WriteHeader(status)
	rw.Write(jsonString)
}

// Format Error to JSON with moreInfo link
func (c *Context) RenderError(rw web.ResponseWriter, message string, status int) {
	js, err := json.MarshalIndent(&JSONError{message, c.Config.Portal}, "", "    ")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		log.Error("RenderError failed to marshall JSONError", err)
		return
	}
	header := rw.Header()
	header.Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)
}
