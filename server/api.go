package server

import (
	"encoding/json"
	"github.com/gocraft/web"
	"net/http"
)

type Context struct {
}

func (a *Context) Root(rw web.ResponseWriter, req *web.Request) {
	a.RenderError(rw, "no api root", http.StatusForbidden)
}

////////////
// Helpers

// Renders from struct to JSON
func (a *Context) Render(rw web.ResponseWriter, model interface{}, status int) {
	jsonString, err := json.MarshalIndent(&model, "", "    ")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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

func (a *Context) RenderError(rw web.ResponseWriter, message string, status int) {

	js, err := json.MarshalIndent(&JSONError{message, "---"}, "", "    ")

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	header := rw.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	header.Set("Access-Control-Allow-Headers", "Authorization")
	rw.WriteHeader(status)
	rw.Write(js)
}
