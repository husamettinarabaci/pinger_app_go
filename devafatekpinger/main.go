package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/AfatekDevelopers/ping_lib_go/devafatekping"
)

func main() {

	var pingCount int
	var timeOutSecond int
	var ipAddress string
	var verbose bool
	flag.IntVar(&pingCount, "count", 5, "set ping count")
	flag.IntVar(&timeOutSecond, "timeout", 5, "set ping timeout second")
	flag.StringVar(&ipAddress, "address", "192.168.1.1", "set ping target address")
	flag.BoolVar(&verbose, "v", true, "show ping log values")
	flag.Parse()

	pinger, err := devafatekping.NewPinger(ipAddress)
	pinger.SetPrivileged(true)
	pinger.Count = pingCount
	pinger.Timeout = time.Duration(time.Second * time.Duration(timeOutSecond))
	logErr(err)

	pinger.OnRecv = func(pkt *devafatekping.Packet) {
		if verbose {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
	}
	var retVal bool = false
	pinger.OnFinish = func(stats *devafatekping.Statistics) {
		if verbose {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

			if stats.PacketLoss == 100 {
				retVal = false
			} else {
				retVal = true
			}
		}
	}
	err = pinger.Run()
	logErr(err)
	if verbose {
		if retVal {
			logStr("PING SUCCESS")
		} else {
			logStr("PING FAIL")
		}
	}
	return
}

func logErr(err error) {
	if err != nil {
		logStr(err.Error())
		panic(err)
	}
}

func logStr(value string) {
	fmt.Println(value)
}
