package monitor

import "time"

type baseMonitor struct{
	host string
	interval time.Duration
	running bool
	terminus chan bool
}

type Monitor interface {
	Run()
	Terminate()
	Host() string
	SetHost(string) error
	Interval() int
	SetInterval(int) error
	Running() bool
}

// Template for unmarshaling of yaml config file
type ConfTemplate struct{
	Monitored map[string]confMonitorModules `yaml:"Monitored"` // List of monitored hosts and their parameters
}


// List of modules for specific host
type confMonitorModules struct{
	Ping confPingModule `yaml:"ping"`
}

// Configuration for icmp module
type confPingModule struct {
	Active bool `yaml:"active"`
	Interval int `yaml:"interval"`
	DeadAfter int `yaml:"deadAfter"`
}