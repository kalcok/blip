package main

import (
	"fmt"
	"blip/monitor"
	"os"
	"os/signal"
	"time"
)

func stop(mon monitor.Monitor, sig_chan chan os.Signal){
	fmt.Println("Gracefully stopping Blip.")
	mon.Terminate()
	for true{
		select {
		case <- sig_chan:
			fmt.Println("Force stopping Blip")
			return
		case <-time.After(100 * time.Millisecond):
			if !mon.Running(){return}
		}
	}
}

func run(){
	fmt.Println(" Hello world")
	var host string
	host = "google.com"
	mon, err:= monitor.NewPingMonitor(host, 0)
	if err != nil{
		panic(fmt.Sprintf("Failed to create Ping Monitor for %s", host))
	}
	mon.SetInterval(15)
	go mon.Run()
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, os.Interrupt)
	<- sig_chan
	stop(mon, sig_chan)

}

func main(){
	run()
}