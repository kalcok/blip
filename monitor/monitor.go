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