package config

import "net/http/httputil"

// instance of a backend.
type Backend struct {
	Url          string                 `json:"url"`
	Port         int                    `json:"port"`
	IsAlive      bool                   `json:"is_alive"`
	ReverseProxy *httputil.ReverseProxy `json:"-"`
}

func NilBackend() Backend {
	return Backend{}
}

// Struct representation of json config file.
// config file SHOULD NOT be manually-edited, as it
// fully managed by the program.
type Config struct {
	IsAlive        bool      `json:"is_alive"`
	PortLowerLimit int       `json:"port_lower_limit"`
	PortUpperLimit int       `json:"port_upper_limit"`
	Backends       []Backend `json:"backends"`
}

func (c Config) Equals(another Config) bool {
	if c.PortLowerLimit != another.PortLowerLimit ||
		c.PortUpperLimit != another.PortUpperLimit ||
		len(c.Backends) != len(another.Backends) {
		return false
	}

	for index, b := range c.Backends {
		if another.Backends[index].Url != b.Url && another.Backends[index].Port != b.Port {
			return false
		}
	}

	return true
}