package types

// This file defines the Go structs that map to the Clash config YAML structure.

// Proxy represents a single proxy entry.
type Proxy map[string]interface{}

// ProxyGroup represents a proxy-group entry.
type ProxyGroup struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Proxies  []string `yaml:"proxies"`
	URL      string   `yaml:"url,omitempty"`
	Interval int      `yaml:"interval,omitempty"`
}

// Config represents the overall Clash configuration structure.
type Config struct {
	Port               int          `yaml:"port,omitempty"`
	SocksPort          int          `yaml:"socks-port,omitempty"`
	AllowLan           bool         `yaml:"allow-lan,omitempty"`
	Mode               string       `yaml:"mode,omitempty"`
	LogLevel           string       `yaml:"log-level,omitempty"`
	ExternalController string       `yaml:"external-controller,omitempty"`
	Proxies            []Proxy      `yaml:"proxies"`
	ProxyGroups        []ProxyGroup `yaml:"proxy-groups"`
	Rules              []string     `yaml:"rules"`
}
