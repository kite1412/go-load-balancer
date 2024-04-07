package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	logger = log.New(log.Writer(), "<CONFIG>", log.Lshortfile)
)

// Json config file contains all neccessary attributes
// needed by load balancer.
type LoadBalancerConfig interface {

	// Initialize a config, typically creating a new config file.
	initConfig() error

	// Set config with a new one and return it in struct form.
	SetConfig(Config) (*Config, error)

	// Return the config in a struct form.
	GetConfig() (*Config, error)

	GetPortLowerLimit() int

	GetPortUpperLimit() int

	AddBackend(Backend)
}

type Default struct {
	Filepath       string
	PortLowerLimit int
	PortUpperLimit int
	m              sync.RWMutex
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

// Send signal to all listed backends,
// updating backends' state accordingly.
// 
// TODO: let the load balancer instance handle this instead.
func (c *Default) pingBackends() {
	for {
		con, err := c.GetConfig()

		if err != nil {
			fmt.Println(err.Error())
		}

		if con.IsAlive {
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
			<-time.NewTicker(time.Second * 30).C
		} else {
			logger.Println("Waiting for server...")
			<-time.NewTicker(time.Second * 3).C
		}
	}
}

func (c *Default) initConfig() error {
	con := Config{
		PortLowerLimit: c.PortLowerLimit,
		PortUpperLimit: c.PortUpperLimit,
		Backends:       make([]Backend, 0),
	}
	if _, err := os.Stat(c.Filepath); err == nil {
		cfg, e := c.GetConfig()
		if e != nil {
			logger.Println(e)
			return e
		}
		if cfg.IsAlive {
			cfg.IsAlive = false
			c.SetConfig(*cfg)
		}
		return nil
	}
	_, err := c.SetConfig(con)
	return err
}

func (c *Default) SetConfig(newConf Config) (*Config, error) {
	c.m.Lock()
	defer c.m.Unlock()
	if err := Write(c.Filepath, newConf); err != nil {
		return nil, err
	}
	c.PortLowerLimit = newConf.PortLowerLimit
	c.PortUpperLimit = newConf.PortUpperLimit
	return &newConf, nil
}

func (c *Default) GetConfig() (*Config, error) {
	return Read(c.Filepath)
}

func (c *Default) GetPortLowerLimit() int {
	return c.PortLowerLimit
}

func (c *Default) GetPortUpperLimit() int {
	return c.PortUpperLimit
}

func (c *Default) AddBackend(newB Backend) {
	con, err := c.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Backends = append(con.Backends, newB)
	c.SetConfig(*con)
}

func (c *Default) AddBackends(backends ...Backend) {
	con, err := c.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Backends = append(con.Backends, backends...)
	c.SetConfig(*con)
}
