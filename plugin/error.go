package plugin

import (
	log "github.com/Sirupsen/logrus"
)

type PluginError interface {
	Error() (string, int)
}

type pluginError struct {
	Status  int
	Message string
}

func NewError(status int, message string) PluginError {
	if status < 200 {
		log.Error("plugin.NewError: Status code ", status, " can not be smaller than 200")
	} else if status > 599 {
		log.Error("plugin.NewError: Status code ", status, " can not be greater than 599")
	}
	return &pluginError{status, message}
}

func (e *pluginError) Error() (string, int) {
	return e.Message, e.Status
}
