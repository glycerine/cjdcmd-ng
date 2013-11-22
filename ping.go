package main

import (
	"fmt"
	"github.com/3M3RY/go-cjdns/cjdns"
	"github.com/spf13/cobra"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"
)

const minInterval = time.Millisecond * 200

var PingCmd = &cobra.Command{
	Use:   "ping HOST",
	Short: "pings a host",
	Long:  `Sends a CJDNS level ping to a given host.`,
	Run:   ping,
}

func ping(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	//TODO(emery): count and interval should be set
	var count int
	var interval time.Duration
	//if interval.Nanoseconds() < minInterval {
	//	fmt.Println("reducing interval")
	interval = minInterval
	//}

	var err error
	var addr, host string

	if cjdns.IsAddress(args[0]) {
		addr = args[0]
		host = Resolve(addr)
	} else {
		host = args[0]
		if addr, err = cjdns.Resolve(host); err != nil {
			fmt.Fprintln(os.Stderr, "Could not resolve "+host+": "+err.Error())
			os.Exit(1)
		}
	}

	var version string
	var ms, minT, avgT, maxT, transmitted, received float32
	minT = math.MaxFloat32

	var start time.Time
	printSummary := func() {
		duration := time.Since(start)
		var loss float32
		switch {
		case received == 0:
			loss = 100
		case received == transmitted:
			loss = 0
		default:
			loss = (received / transmitted) * 100.0
		}

		fmt.Fprint(os.Stdout, "\n--- "+host+" ---\n")
		fmt.Fprintf(os.Stdout, "%.0f pings transmitted, %.0f received, %2.0f%% ping loss, time %s\n", transmitted, received, loss, duration)
		if received != 0 {
			avgT /= received
			fmt.Fprintf(os.Stdout, "rtt min/avg/max = %2.f/%.2f/%.2f ms\n", minT, avgT, maxT)
			fmt.Fprintln(os.Stdout, "CJDNS version: "+version)

		}
		os.Exit(0)
	}

	mu := new(sync.Mutex)
	ping := func() {
		version, ms, err = Admin.RouterModule_pingNode(addr, 0)
		mu.Lock()
		transmitted++

		if err != nil {
			fmt.Fprintf(os.Stdout, "error: %s\n", err)
		} else {
			received++
			fmt.Fprintf(os.Stdout, "time=%.3f ms\n", ms)
			switch {
			case ms < minT:
				minT = ms
			case ms > maxT:
				maxT = ms
			}
			avgT += ms
		}
		mu.Unlock()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	fmt.Fprintf(os.Stdout, "PING %s (%s)\n", host, addr)
	start = time.Now()
	go ping()
	for i := count - 1; i != 0; i-- {
		select {
		case <-c:
			printSummary()

		case <-time.After(interval):
			go ping()
		}
	}
	printSummary()
}
