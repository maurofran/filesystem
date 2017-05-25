package filesystem

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

// FindPlugin will find a plugin for given method.
func (p *Pluggable) FindPlugin(method string) (Plugin, error) {
	plugin, ok := p.plugins[method]
	if !ok {
		return nil, pluginNotFoundError(method)
	}
	return plugin, nil
}

// InvokePlugin will invoke the plugin on provided manager.
func (p *Pluggable) InvokePlugin(filesystem Interface, method string, args ...interface{}) (interface{}, error) {
	plugin, err := p.FindPlugin(method)
	if err != nil {
		return nil, err
	}
	plugin.SetFileSystem(filesystem)
	return plugin.Handle(args...)
}
