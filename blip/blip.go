package main

import (
	"fmt"
	"blip/monitor"
)

func run(){
	fmt.Println(" Hello world")
	var input, host string
	host = "google.com"
	mon, err:= monitor.NewPingMonitor(host, 0)
	if err != nil{
		panic(fmt.Sprintf("Failed to create Ping Monitor for %s", host))
	}
	mon.SetInterval(15)
	go mon.Run()
	fmt.Scanln(&input)
	mon.Terminate()
	fmt.Scanln(&input)
}

func main(){
	run()
}