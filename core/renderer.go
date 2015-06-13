package core

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"net/http"
)

type Renderer struct{}

// Responds from error to JSON
type JSONError struct {
	Message  string `json:"message"`
	MoreInfo string `json:"moreInfo"`
}

// Renders from struct to JSON
func (a *Renderer) Render(rw web.ResponseWriter, model interface{}, status int) {
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
func (c *Renderer) RenderError(rw web.ResponseWriter, err error, info string, status int) {
	js, err := json.MarshalIndent(&JSONError{err.Error(), info}, "", "    ")
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

// func (r *Renderer) Render500(rw *web.ResponseWriter, err error) {
// 	r.RenderError(rw, err, "http://api.exmaple.com/errors/#/500", http.StatusInternalServerError)
// }
