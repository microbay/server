package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"net/http"
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

func (p *NoopPlugin) Outbound(rw web.ResponseWriter, req *web.Request) (int, error) {
	log.Warn("NoopPlugin::Outbound")
	var err error
	return http.StatusOK, err
}
