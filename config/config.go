package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"
)

var (
	logger = log.New(log.Writer(), "<CONFIG>", log.Lshortfile)
)

type Backend struct {
	Url          string                 `json:"url"`
	Port         int                    `json:"port"`
	IsAlive      bool                   `json:"is_alive"`
	ReverseProxy *httputil.ReverseProxy `json:"-"`
}

func NilBackend() Backend {
	return Backend{}
}

// struct representation of json config file.
type Config struct {
	PortLowerLimit int       `json:"port_lower_limit"`
	PortUpperLimit int       `json:"port_upper_limit"`
	Backends       []Backend `json:"backends"`
}

func (c Config) Equal(another Config) bool {
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

// json config file contains all neccessary attributes
// needed by load balancer.
type LoadBalancerConfig interface {
	// initialize a config, typically creating a new config file.
	initConfig() error
	// set a new config with a new one and return it in struct form.
	SetConfig(Config) (*Config, error)
	// return the config in a struct form.
	GetConfig() (*Config, error)
	GetPortLowerLimit() int
	GetPortUpperLimit() int
	AddBackend(Backend)
}

type Default struct {
	Filepath       string
	PortLowerLimit int
	PortUpperLimit int
	// accessing this field is never recommended, use GetConfig instead.
	// just make sure to access this field only when there's no on-going modification
	// of the config file.
	config Config
	m      sync.RWMutex
}

func DefaultConfig(filepath string, ll int, ul int) (*Default, error) {
	c := &Default{
		Filepath:       filepath,
		PortLowerLimit: ll,
		PortUpperLimit: ul,
	}
	if err := c.initConfig(); err != nil {
		return nil, err
	}
	go c.pingBackends()
	return c, nil
}

// send signal to all listed backends,
// updating backends' state accordingly.
func (c *Default) pingBackends() {
	for {
		con, err := c.GetConfig()
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(con.Backends) != 0 {
			logger.Println("<<<SENDING PING TO ALL BACKENDS>>>")
			backends := con.Backends
			for i, b := range backends {
				_, err := http.Get(b.Url)
				if err != nil {
					con.Backends[i].IsAlive = false
				} else {
					con.Backends[i].IsAlive = true
				}
			}
			c.SetConfig(*con)
		} else {
			logger.Println("<<<NO BACKEND INSTANCES>>>")
		}
		<- time.NewTicker(time.Minute).C
	}
}

func (c *Default) initConfig() error {
	con := Config{
		PortLowerLimit: c.PortLowerLimit,
		PortUpperLimit: c.PortUpperLimit,
		Backends:       make([]Backend, 0),
	}
	if _, err := os.Stat(c.Filepath); err == nil {
		con, e := c.GetConfig()
		if e != nil {
			fmt.Println(e)
			return e
		}
		c.config = *con
		return nil
	}
	_, err := c.SetConfig(con)
	return err
}

func (c *Default) SetConfig(newConf Config) (*Config, error) {
	c.m.Lock()
	defer c.m.Unlock()
	json, err := json.Marshal(newConf)
	if err != nil {
		fmt.Println("fail to create a json")
		fmt.Println(err)
		return nil, err
	}
	if err := os.WriteFile(c.Filepath, json, os.ModePerm); err != nil {
		fmt.Println("fail to create the config file")
		fmt.Println(err)
		return nil, err
	}
	c.config = newConf
	c.PortLowerLimit = c.config.PortLowerLimit
	c.PortUpperLimit = c.config.PortUpperLimit
	return &newConf, nil
}

func (c *Default) GetConfig() (*Config, error) {
	con, err := os.ReadFile(c.Filepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	prevC := c.config
	cfg := Config{}
	if err := json.Unmarshal(con, &cfg); err != nil {
		fmt.Println(err)
		return nil, err
	}
	copy(cfg.Backends, prevC.Backends)
	if !cfg.Equal(c.config) {
		c.m.Lock()
		c.config = cfg
		c.m.Unlock()
	}
	return &cfg, nil
}

func (c *Default) GetPortLowerLimit() int {
	return c.PortLowerLimit
}

func (c *Default) GetPortUpperLimit() int {
	return c.PortUpperLimit
}

func (c *Default) AddBackend(newB Backend) {
	con := &c.config
	con.Backends = append(con.Backends, newB)
	c.SetConfig(*con)
}

func (c *Default) AddBackends(backends ...Backend) {
	con := &c.config
	con.Backends = append(con.Backends, backends...)
	c.SetConfig(*con)
}
