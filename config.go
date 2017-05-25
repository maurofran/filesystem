package filesystem

// Config is a configuration object.
type Config struct {
	settings map[string]interface{}
	fallback *Config
}

// EmptyConfig will create a new empty configuration.
func EmptyConfig() *Config {
	return &Config{}
}

// Get a setting.
func (c *Config) Get(key string, def interface{}) interface{} {
	if v, ok := c.settings[key]; ok {
		return v
	}
	if c.fallback != nil {
		return c.fallback.Get(key, def)
	}
	return def
}

// Has will check if an item exists by key.
func (c *Config) Has(key string) bool {
	if _, ok := c.settings[key]; ok {
		return true
	}
	if c.fallback != nil {
		return c.fallback.Has(key)
	}
	return false
}

// GetDefault wil try to retrieve a default setting from a config fallback.
func (c *Config) GetDefault(key string, def interface{}) interface{} {
	if c.fallback == nil {
		return def
	}
	return c.fallback.Get(key, def)
}

// Set a setting.
func (c *Config) Set(key string, val interface{}) {
	c.settings[key] = val
}

// SetFallback will set the fallback.
func (c *Config) SetFallback(fallback *Config) {
	c.fallback = fallback
}

// Configurable is a struct holding a configuration object instance and provide methods to interact with this configuration.
type Configurable struct {
	config *Config
}

// Config is the getter method for configuration object.
func (c *Configurable) Config() *Config {
	return c.config
}

// SetConfig will set the configuration.
func (c *Configurable) SetConfig(config *Config) {
	c.config = config
}

// PrepareConfig will convert a map into a configuration object with right fallback values.
func (c *Configurable) PrepareConfig(config map[string]interface{}) *Config {
	cfg := EmptyConfig()
	for k, v := range config {
		cfg.Set(k, v)
	}
	cfg.SetFallback(c.Config())
	return cfg
}
