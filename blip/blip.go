package main

import (
	"fmt"
	"blip/monitor"
	"os"
	"os/signal"
	"time"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func stop(monitors []monitor.Monitor, sig_chan chan os.Signal){
	fmt.Println("Gracefully stopping Blip.")
	for _, mon := range monitors{
		mon.Terminate()
	}

	for true{
		select {
		case <- sig_chan:
			fmt.Println("Force stopping Blip")
			return
		case <-time.After(100 * time.Millisecond):
			for _, mon := range monitors{
				if mon.Running(){break}
			}
			return
		}
	}
}

func run(blip_conf *monitor.ConfTemplate){
	var running_monitors []monitor.Monitor
	for hostname, host_config := range blip_conf.Monitored {
		if host_config.Ping.Active{
			ping_conf := host_config.Ping
			mon, err := monitor.NewPingMonitor(hostname, ping_conf.Interval)
			if err != nil {
				panic(fmt.Sprintf("Failed to create Ping Monitor for %s", hostname))
			}
			go mon.Run()
			running_monitors = append(running_monitors, mon)
		}
	}
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, os.Interrupt)
	<- sig_chan
	stop(running_monitors, sig_chan)

}

func parse_conf(path string, dump bool) (cnf *monitor.ConfTemplate){
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil{
		panic(err)
	}

	cnf = new(monitor.ConfTemplate)

	err = yaml.Unmarshal(yamlFile, &cnf)
	if err != nil{
		panic(err)
	}
	if dump {
		for host, config := range cnf.Monitored {
			fmt.Printf("Will monitor %s\n", host)
			if config.Ping.Active {
				fmt.Println("  Using Ping")
				fmt.Printf("    Interval: %d s\n", config.Ping.Interval)
				fmt.Printf("    Warning after %d failed attempts\n", config.Ping.DeadAfter)

			} else {
				fmt.Println("  Skipping Ping")
			}
		}

	}
	return
}
func main(){
	conf := parse_conf("./config.yaml", true)
	run(conf)
}