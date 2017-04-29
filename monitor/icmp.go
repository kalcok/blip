package monitor

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"time"
	"errors"
)

const (
	PING_OK   = iota
	PING_FAIL = iota
)

var MIN_ICMP_INTERVAL = 10

// XXX: Hack to pick ipv4 address from list of addresses for resolved hostname
func pick_ipv4(addresses []string) string{
	for i := 0; i < len(addresses); i++{
		ip := net.ParseIP(addresses[i])
		if ip.To4() != nil{
			return addresses[i]
		}
	}
	return ""
}

type PingMonitor struct {
	baseMonitor
}

func NewPingMonitor(host string, interval int) (monitor *PingMonitor, err error) {
	monitor = nil
	err = nil
	monitor = new(PingMonitor)
	ok := monitor.SetHost(host)
	if ok != nil {
		monitor = nil
		err = ok
		return
	}
	ok = monitor.SetInterval(interval)
	if ok != nil {
		monitor = nil
		err = ok
		return
	}
	monitor.terminus = make(chan bool)
	monitor.running = false
	return
}

func (monitor *PingMonitor) Host() string {
	return monitor.host
}

func (monitor *PingMonitor) SetHost(host string) (err error) {
	if len(host) < 1 {
		err = errors.New("Hostname can't be empty.")
		return
	}
	monitor.host = host
	return
}

func (monitor *PingMonitor) Interval() int {
	return int(monitor.interval / time.Second)
}

func (monitor *PingMonitor) SetInterval(interval int) (err error) {
	if interval == 0 {
		interval = MIN_ICMP_INTERVAL
	}
	if interval < MIN_ICMP_INTERVAL {
		err = errors.New(fmt.Sprintf("Can't execute monitoring actions in interval lower than %d seconds.", MIN_ICMP_INTERVAL))
		return
	}

	monitor.interval = time.Duration(interval) * time.Second
	return
}

func (monitor *PingMonitor) Run() {
	run := true
	monitor.running = true

	for run {
		fmt.Printf("Monitoring %s with ping...", monitor.host)
		err := monitor.ping()
		if err != nil {
			fmt.Printf("\nPing to %s failed. %s", monitor.host, err)
		} else {
			fmt.Print("OK\n")
		}
		select {
		case _, ok := <-monitor.terminus:
			if !ok {
				fmt.Println("Recieved termination order.")
				run = false
				break
			}
		case <-time.After(monitor.interval):
		}
	}
	fmt.Println("Graciously exiting go routine")
	monitor.running = false
}

func (monitor *PingMonitor) Terminate() {
	close(monitor.terminus)
}

func (monitor *PingMonitor) Running() bool{
	return monitor.running
}

func (monitor *PingMonitor) ping() (err error) {
	host_ip, err := net.LookupHost(monitor.host)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to resolve hostname. %s", err))
		return
	}
	// TODO: Loose this hack, start supporting IPv6
	host_ip_4 := pick_ipv4(host_ip)
	if host_ip_4 == ""{
		err = errors.New(fmt.Sprintf("Failed to resolve hostname '%s' to IPv4 address. IPv6 not supported", monitor.host))
		return
	}
	target := &net.UDPAddr{IP: net.ParseIP(host_ip_4)}
	// TODO: Check for platform compatibility (darwin+linux)
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	defer conn.Close()

	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to open connection . %s", err))
		return
	}

	message := icmp.Message{Type: ipv4.ICMPTypeEcho,
		Code:                 0,
		Body: &icmp.Echo{ID: 666,
			Seq:         1,
			Data:        []byte("You are blip on my radar")}}

	serial_message, err := message.Marshal(nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to serialize ICMP body. %s", err))
	}
	_, err = conn.WriteTo(serial_message, target)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to send ICMP request to host '%s'. %s", monitor.host, err))
		return
	}

	rcv_buffer := make([]byte, 1500)

	reply_len, _, err := conn.ReadFrom(rcv_buffer)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed read response from %s. %s", monitor.host, err))
		return
	}

	_, err = icmp.ParseMessage(1, rcv_buffer[:reply_len])
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to parse repsponse message. %s", err))
		return
	}

	return
}