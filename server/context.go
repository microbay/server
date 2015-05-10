package server

import (
	"encoding/json"
	"github.com/SHMEDIALIMITED/apigo/model"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"net/http"
)

// Root context
type Context struct {
	Config   model.API
	Resource *model.Resource
	Session  *Session
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

// Responds from error to JSON
type JSONError struct {
	Message  string `json:"message"`
	MoreInfo string `json:"moreInfo"`
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
