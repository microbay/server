package plugin

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	PLUGIN_TRANSFORMER string = "transformer"
)

type TransformerPlugin struct{}

func (p *TransformerPlugin) Bootstrap(config map[string]interface{}) (Interface, error) {
	log.Warn("TransformerPlugin::Bootstrap")
	var err error
	return p, err
}

func (p *TransformerPlugin) Inbound(req *web.Request) (int, error) {
	log.Warn("TransformerPlugin::Inbound")
	var err error
	return http.StatusOK, err
}

func (p *TransformerPlugin) Outbound(res *http.Response) (int, error) {
	//log.Warn("TransformerPlugin::Outbound", res.Header.Get("Content-Encoding"))

	var err error

	log.Debug("TransformerPlugin::Outbound:", res.Body)
	//savecl := res.ContentLength
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(res.Body); err != nil {
		return http.StatusInternalServerError, err
	}
	if err = res.Body.Close(); err != nil {
		return http.StatusInternalServerError, err
	}

	// log.Warn("TransformerPlugin::Outbound PARSED")
	var b map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &b)
	if err != nil {
		log.Error("TransformerPlugin::Outbound: unmarshal body failed, ", err)
		return http.StatusInternalServerError, err
	}

	b["name"] = "h"

	var j []byte
	j, err = json.Marshal(b)

	if err != nil {
		log.Error("TransformerPlugin::Outbound: unmarshal body failed, ", err)
		return http.StatusInternalServerError, err
	}

	log.Warn("TransformerPlugin::Outbound: ", len(buf.Bytes()), http.DetectContentType(buf.Bytes()))
	log.Warn("TransformerPlugin::Outbound: ", len(j), http.DetectContentType(j))
	//buf.Reset()
	//buf.Write(j)

	res.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	//res.Body.Close()

	return http.StatusOK, err
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// // debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
