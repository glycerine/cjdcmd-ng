package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/3M3RY/go-cjdns/cjdns"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var TraceCmd = &cobra.Command{
	Use:   "trace HOST [HOST...]",
	Short: "prints routes to hosts",
	Long:  `Parses the local routing table and prints the routes to the given hosts.`,
	Run:   trace,
}

func trace(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Usage()
		os.Exit(1)
	}
	var run *NmapRun
	startTime := time.Now()
	if NmapOutput {
		args := fmt.Sprint(os.Args[:])
		run = &NmapRun{
			Scanner:          "cjdmap",
			Args:             args,
			Start:            startTime.Unix(),
			Startstr:         startTime.String(),
			Version:          "0.1",
			XMLOutputVersion: "1.04",
			Hosts:            make([]*Host, 0, len(args)),
		}
	}

	targets := make([]*target, 0, len(args))
	for _, arg := range args {
		target, err := newTarget(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping %s: %s\n", arg, err)
			continue
		}
		targets = append(targets, target)
	}

	table, err := Admin.NodeStore_dumpTable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get routing table:", err)
	}
	table.SortByPath()

	for _, target := range targets {
		if NmapOutput {
			fmt.Fprintln(os.Stderr, target)
		} else {
			fmt.Fprintln(os.Stdout, target)
		}
		traces, err := target.trace(table)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Faild to trace %s: %s\n", target, err)
			continue
		}
		if NmapOutput {
			run.Hosts = append(run.Hosts, traces[0])
		}
	}
	if NmapOutput {
		stopTime := time.Now()
		run.Finished = &Finished{
			Time:    stopTime.Unix(),
			TimeStr: stopTime.String(),
			//Elapsed: (stopTime.Sub(startTime) * time.Millisecond).String(),
			Exit: "success",
		}

		fmt.Fprint(os.Stdout, xml.Header)
		fmt.Fprintln(os.Stdout, `<?xml-stylesheet href="file:///usr/bin/../share/nmap/nmap.xsl" type="text/xsl"?>`)
		xEnc := xml.NewEncoder(os.Stdout)
		err = xEnc.Encode(run)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

type target struct {
	addr string
	name string
	rtt  time.Duration
	xml  *Host
}

func (t *target) String() string {
	if len(t.name) != 0 {
		return t.name + " (" + t.addr + ")"
	}
	return t.addr
}

func newTarget(host string) (t *target, err error) {
	t = new(target)
	if cjdns.IsAddress(host) {
		t.addr = host
		if ResolveNodeinfo {
			t.name = NodeinfoReverse(host)
		}
	} else {
		t.name = host
		t.addr, err = cjdns.Resolve(host)
	}
	return
}

var notInTableError = errors.New("not found in routing table")

func (t *target) trace(table cjdns.Routes) (hostTraces []*Host, err error) {
	hostTraces = make([]*Host, 0, 2)
	for _, r := range table {
		if r.IP == t.addr {
			hops := table.Hops(r)
			if hostTrace, err := t.traceHops(hops); err != nil {
				fmt.Fprintf(os.Stderr, "failed to trace %s, %s\n", t, err)
			} else {
				hostTraces = append(hostTraces, hostTrace)
			}
			fmt.Fprintln(os.Stderr)
		}
	}
	if len(hostTraces) == 0 {
		hostTraces = nil
		err = notInTableError
	}
	return
}

func (t *target) traceHops(hops cjdns.Routes) (*Host, error) {
	hops.SortByPath()
	startTime := time.Now().Unix()
	trace := &Trace{Proto: "CJDNS"}
	for y, p := range hops {
		if y == 0 {
			continue
		}

		// Ping by path so we don't get RTT for a different route.
		_, rtt, err := Admin.SwitchPinger_ping(p.Path, "", 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
		if rtt == 0 {
			rtt = 1
		}

		hop := &Hop{
			TTL:    y,
			RTT:    rtt,
			IPAddr: p.IP,
		}
		if ResolveNodeinfo {
			hop.Host = NodeinfoReverse(p.IP)
		}
		if NmapOutput {
			fmt.Fprintf(os.Stderr, "  %02d.% 4dms %s %s %s\n", y, rtt, p.Path, p.IP, hop.Host)
		} else {
			fmt.Fprintf(os.Stdout, "  %02d.% 4dms %s %s %s\n", y, rtt, p.Path, p.IP, hop.Host)
		}
		trace.Hops = append(trace.Hops, hop)
	}

	endTime := time.Now().Unix()
	h := &Host{
		StartTime: startTime,
		EndTime:   endTime,
		Status: &Status{
			State:     HostStateUp,
			Reason:    "pingNode",
			ReasonTTL: 56,
		},
		Address: newAddress(t.addr),
		Trace:   trace,
		//Times: &Times{ // Don't know what to do with this element yet.
		//	SRTT:   1,
		//	RTTVar: 1,
		//	To:     1,
		//},
	}

	if t.name != "" {
		h.Hostnames = []*Hostname{&Hostname{Name: t.name, Type: HostnameTypeUser}}
	}
	return h, nil
}
