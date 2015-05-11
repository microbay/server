package server

import (
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gocraft/web"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (c *Context) ResourceConfigMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	var err error
	c.Resource, err = c.Config.FindResourceByRequest(req.Request)
	if err != nil {
		c.RenderError(rw, "Access Forbidden", http.StatusForbidden)
	} else {
		next(rw, req)
	}
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
	combinedHeaders := headerCombiner(reverseProxy, c.Config.Key)
	//}

	// if c.Resource.Auth == REDIS_JWT {
	//  c.RenderError(rw, "Invalid token", http.StatusUnauthorized)
	// } else {
	//  next(rw, req)
	// }

	req.URL.Path = ""
	combinedHeaders.ServeHTTP(rw, req.Request)
}

// Append additional query params to the original URL query.
func headerCombiner(handler http.Handler, key []byte) http.Handler {

	c := map[string]interface{}{"premise_id": -2287340764, "exp": float64(time.Now().Unix() + 100)}

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = c
	tokenString, e := token.SignedString(key)

	if e != nil {
		panic(e.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", tokenString)
		r.Header.Set("X-Premise", "-2287340764")
		handler.ServeHTTP(w, r)
	})
}
