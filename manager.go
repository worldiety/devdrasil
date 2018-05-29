package main

import (
	"fmt"
	"github.com/worldiety/devdrasil/plugin"
	"github.com/worldiety/devdrasil/plugin/session"
)

var _ plugin.Manager = (*sessionPluginManager)(nil)

type sessionPluginManager struct {
	devdrasil *Devdrasil
	apiDoc    *plugin.ApiDoc
}

func (p *sessionPluginManager) Id() string {
	return session.PLUGIN_ID
}

func (p *sessionPluginManager) Deploy() (string, error) {
	return p.Id(), nil
}

func (p *sessionPluginManager) List() ([]plugin.ManagedInstance, error) {
	return []plugin.ManagedInstance{{p.Id(), plugin.STATUS_RUNNING, p.devdrasil.host, p.devdrasil.port, "/session/","asdf"}}, nil
}

func (p *sessionPluginManager) Start(id string) error {
	return fmt.Errorf("session plugin is builtin and always running")
}

func (p *sessionPluginManager) Stop(id string) error {
	return fmt.Errorf("session plugin is builtin and cannot be stopped")
}

func (p *sessionPluginManager) Doc() *plugin.ApiDoc {
	return p.apiDoc
}
