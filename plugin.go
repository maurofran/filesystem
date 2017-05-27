package filesystem

import (
	"fmt"
)

// Plugin is the interface implemented by plugins.
type Plugin interface {
	// The method name exposed by plugin.
	Method() string
	// Set the manager on plugin.
	SetFileSystem(Interface)
	// Handle the method invocation.
	Handle(args ...interface{}) (interface{}, error)
}

// Pluggable is a base struct for pluggable behavior.
type Pluggable struct {
	plugins map[string]Plugin
}

// AddPlugin will add a plugin to pluggable
func (p *Pluggable) AddPlugin(plugin Plugin) {
	p.plugins[plugin.Method()] = plugin
}

// FindPlugin will find a plugin for given method, returning an error if no plugin exists for provided method.
func (p *Pluggable) FindPlugin(method string) (Plugin, error) {
	plugin, ok := p.plugins[method]
	if !ok {
		return nil, pluginNotFoundError{missingPlugin: method}
	}
	return plugin, nil
}

// InvokePlugin will invoke the plugin on provided manager.
func (p *Pluggable) InvokePlugin(filesystem Interface, method string, args ...interface{}) (interface{}, error) {
	plugin, err := p.FindPlugin(method)
	if err != nil {
		// We don't wrap error because it's a call of internal method and any returned error is already the final error.
		return nil, err
	}
	plugin.SetFileSystem(filesystem)
	return plugin.Handle(args...)
}

type pluginNotFoundError struct {
	missingPlugin string
}

func (e pluginNotFoundError) Error() string {
	return fmt.Sprintf("Plugin not found for method: %s", e.missingPlugin)
}

func (e pluginNotFoundError) MissingPlugin() string {
	return e.missingPlugin
}

// IsPluginNotFound will check if provided error is a "plugin not found" error.
func IsPluginNotFound(err error) (bool, string) {
	type pluginNotFound interface {
		MissingPlugin() string
	}

	if pnf, ok := err.(pluginNotFound); ok {
		return ok, pnf.MissingPlugin()
	}
	return false, ""
}
