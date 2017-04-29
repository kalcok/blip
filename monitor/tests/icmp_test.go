package monitor

import (
	"testing"
	"blip/monitor"
	"fmt"
	"time"
)

func initPingMonitor(t *testing.T) *monitor.PingMonitor{
	ping_monitor, err := monitor.NewPingMonitor("localhost", monitor.MIN_ICMP_INTERVAL)
	if err != nil{
		t.Error("Failed to initiate pingMonitor")
	}
	return ping_monitor
}

func TestPingMonitorInterfaceCompliance(t *testing.T){
	defer func(){
		err := recover()
		if err != nil{
			t.Error(fmt.Printf("PingMonitor does not complain with Monitor interface. %s", err))
		}
	}()
	var _ monitor.Monitor = (*monitor.PingMonitor)(nil)

}

func TestPingMonitorInitIntervalTooLow(t *testing.T){
	interval_too_low := 9
	ping_monitor, err := monitor.NewPingMonitor("localhost", interval_too_low)
	if err == nil{
		if ping_monitor.Interval() == interval_too_low{
			t.Error("Initiated PingMonitor with interval value below allowed threshold.")
		}
	}
}

func TestPingMonitorSetIntervalTooLow(t *testing.T){
	ping_monitor := initPingMonitor(t)
	err := ping_monitor.SetInterval(9)
	if err ==nil{
		t.Error("Allowed to set PingMonitorInterval below allowed threshold.")
	}
}

func TestPingMonitorInitIntervalDefault(t *testing.T){
	trigger_default := 0
	ping_monitor, err := monitor.NewPingMonitor("localhost", trigger_default)
	if err == nil{
		if ping_monitor.Interval() != monitor.MIN_ICMP_INTERVAL{
			t.Error("Failed to init PingMonitor with default interval value")
		}
	}else{
		t.Error(fmt.Sprintf("Failed to create PingMonitor. %s", err))
	}
}

func TestPingMonitorSetIntervalDefault(t *testing.T){
	ping_monitor := initPingMonitor(t)
	ping_monitor.SetInterval(0)
	if ping_monitor.Interval() != monitor.MIN_ICMP_INTERVAL{
		t.Error("Failed to set PingMonitors default interval value")
	}
}

func TestPingMonitorInitNoHost(t *testing.T){
	no_host := ""
	ping_monitor, err := monitor.NewPingMonitor(no_host, monitor.MIN_ICMP_INTERVAL)
	if err == nil{
		if ping_monitor.Host() == no_host{
			t.Error("Initiated PingMonitor with empty Host.")
		}
	}
}

func TestPingMonitorSetNoHost(t *testing.T){
	ping_monitor := initPingMonitor(t)
	ping_monitor.SetHost("")
	if ping_monitor.Host() == ""{
		t.Error("Allowed to set empty Host")
	}
}

func TestPingMonitorGracefulStop(t *testing.T){
	ping_monitor := initPingMonitor(t)

	go ping_monitor.Run()
	ping_monitor.Terminate()
	time.Sleep(1 * time.Second)
	if ping_monitor.Running(){
		t.Error("Failed to stop PingMonitor routine")
	}

}