/*
 * You may redistribute this program and/or modify it under the terms of
 * the GNU General Public License as published by the Free Software Foundation,
 * either version 3 of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ehmry/go-cjdns/admin"
	"github.com/ehmry/go-cjdns/key"
	"github.com/spf13/cobra"
)

func init() {
	TracerouteCmd.Flags().BoolVarP(&NmapOutput, "nmap", "x", false, "print result in nmap XML to stdout")
}

func tracerouteCmd(cmd *cobra.Command, args []string) {
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

	c := Connect()

	for _, target := range targets {
		if NmapOutput {
			fmt.Fprintln(os.Stderr, target)
		} else {
			fmt.Fprintln(os.Stdout, target)
		}
		trace, err := target.trace(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to trace %s: %s\n", target, err)
			continue
		}
		if NmapOutput {
			run.Hosts = append(run.Hosts, trace)
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
		err := xEnc.Encode(run)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

type target struct {
	addr, name string
	rtt  time.Duration
	xml  *Host
}

func (t *target) String() string {
	if len(t.name) != 0 {
		return fmt.Sprintf("%s (%s)", t.name, t.addr)
	}
	return t.addr
}

func newTarget(host string) (t *target, err error) {
	t = new(target)
	t.name, t.addr, err = resolve(host)
	return
}

var notInTableError = errors.New("not found in routing table")

func (t *target) trace(c *admin.Conn) (hostTrace *Host, err error) {
	var node *admin.StoreNode
	var addr string
	node, err = c.NodeStore_nodeForAddr(t.addr)
	if err != nil {
		return
	}

	var nodes []*admin.StoreNode
	for err == nil && node.RouteLabel != "0000.0000.0000.0001" {
		nodes = append(nodes, node)
		addr = node.BestParent.IP
		node, err = c.NodeStore_nodeForAddr(addr)
	}
	if err != nil {
		return
	}

	startTime := time.Now().Unix()
	trace := &Trace{Proto: "CJDNS"}

	var ttl int
	for i := len(nodes) - 1; i > -1; i-- {
		ttl++
		node = nodes[i]

		// Ping by path so we don't get RTT for a different route.
		rtt, _, err := c.RouterModule_pingNode(node.RouteLabel, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
		if rtt == 0 {
			rtt = 1
		}
		//hostname, _, _ := resolve(p.IP.String())
		pubKey, _ := key.DecodePublic(node.Key)
		ip := pubKey.IP()
		hop := &Hop{
			TTL:    ttl,
			RTT:    rtt,
			IPAddr: ip.String(),
			//Host:   hostname,
		}

		if NmapOutput {
			fmt.Fprintf(os.Stderr, "  %02d.% 4dms %s %s\n", ttl, rtt, node.RouteLabel, ip)
		} else {
			fmt.Fprintf(os.Stdout, "  %02d.% 4dms %s %s\n", ttl, rtt, node.RouteLabel, ip)
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
		Address: &Address{Addr: t.addr, AddrType: "ipv6"},
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
