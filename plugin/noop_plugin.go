package plugin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
)

const (
	PLUGIN_NOOP string = "noop"
)

type NoopPlugin struct{}

func (p *NoopPlugin) Bootstrap(config map[string]interface{}) (Plugin, error) {
	var err error
	return p, err
}

func (p *NoopPlugin) Inbound(req *web.Request) (Plugin, PluginError) {
	log.Debug("NoopPlugin::Inbound")
	var err PluginError
	return p, err
}

func (p *NoopPlugin) Outbound(req *web.Request) (Plugin, error) {
	log.Debug("NoopPlugin::Outbound")
	var err error
	return p, err
}
