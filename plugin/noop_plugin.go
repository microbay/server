package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
)

const (
	PLUGIN_NOOP string = "noop"
)

type NoopPlugin struct{}

func (p *NoopPlugin) Inbound(req *web.Request) {
	log.Debug("NoopPlugin::Inbound")
}

func (p *NoopPlugin) Outbound(req *web.Request) {
	log.Debug("NoopPlugin::Outbound")
}
