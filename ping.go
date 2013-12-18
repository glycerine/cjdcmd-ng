package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"
)

const minInterval = time.Millisecond * 200

var (
	PingCmd = &cobra.Command{
		Use:   "ping HOST",
		Short: "pings a host",
		Long:  `Sends a CJDNS level ping to a given host.`,
		Run:   ping,
	}

	count    int
	interval time.Duration
)

func init() {
	PingCmd.PersistentFlags().IntVarP(&count, "count", "c", -1, "Stop after sending c packets.")
	PingCmd.PersistentFlags().DurationVarP(&interval, "interval", "i", time.Second, " Wait time between sending each packet.")
}

func ping(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	if interval < minInterval {
		fmt.Println("increasing interval to", minInterval)
		interval = minInterval
	}

	host, ip, err := resolve(args[0])
	addr := ip.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not resolve %s: %s", args[0], err)
		os.Exit(1)
	}

	var (
		version                                     string
		ms, minT, avgT, maxT, transmitted, received float32
		msInt                                       int
		start                                       time.Time
	)
	minT = math.MaxFloat32

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
		mu.Lock()
		msInt, version, err = Admin.RouterModule_pingNode(addr, 0)
		transmitted++
		ms = float32(msInt)

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fmt.Fprintf(os.Stdout, "PING %s (%s)\n", host, addr)
	start = time.Now()
	go ping()
	for i := count; i != 0; i-- {
		select {
		case <-sig:
			printSummary()

		case <-time.After(interval):
			go ping()
		}
	}
	printSummary()
}
