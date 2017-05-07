package main

// Template for unmarshaling of yaml config file
type ConfTemplate struct{
	Monitored map[string]confMonitorModules `yaml:"Monitored"` // List of monitored hosts and their parameters
	Logging confLoggers `yaml:"Logging"`
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

// Configuration for logging modules
type confLoggers struct{
	//FileLogger confFileLogger `yaml:"FileLogger"`
	FileLogger confFileLogger `yaml:"FileLogger"`
}
// Configuration for File logger
type confFileLogger struct{
	Level string `yaml:"level"`
	LogFile string `yaml:"logFile"`
}
