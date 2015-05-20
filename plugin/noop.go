package plugin

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"
	//	"net/http/httputil"
	"unicode/utf8"
)

const (
	PLUGIN_NOOP string = "noop"
)

type NoopPlugin struct{}

func (p *NoopPlugin) Bootstrap(config map[string]interface{}) (Interface, error) {
	log.Warn("NoopPlugin::Bootstrap")
	var err error
	return p, err
}

func (p *NoopPlugin) Inbound(req *web.Request) (int, error) {
	log.Warn("NoopPlugin::Inbound")
	var err error
	return http.StatusOK, err
}

func (p *NoopPlugin) Outbound(res *http.Response) (int, error) {
	log.Warn("NoopPlugin::Outbound")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// copying the response body did not work
		return http.StatusInternalServerError, err
	}

	log.Warn("NoopPlugin::Outbound PARSED", utf8.Valid(body))
	var b map[string]interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Error("NoopPlugin::Outbound: unmarshal body failed, ", err)
		resp.Write(&b)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, err
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
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
