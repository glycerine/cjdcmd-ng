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
	count    int
	interval time.Duration
)

func init() {
	PingCmd.PersistentFlags().IntVarP(&count, "count", "c", -1, "Stop after sending c packets.")
	PingCmd.PersistentFlags().DurationVarP(&interval, "interval", "i", time.Second, " Wait time between sending each packet.")
}

func pingCmd(cmd *cobra.Command, args []string) {
	c := Connect()

	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	if interval < minInterval {
		fmt.Println("increasing interval to", minInterval)
		interval = minInterval
	}

	host, ip, err := resolve(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not resolve %s: %s\n", args[0], err)
		os.Exit(1)
	}

	node, err := c.NodeStore_nodeForAddr(ip)
	if err != nil {
		die("Failed to find %s, %s", ip, err)
	}
	path := node.RouteLabel

	var (
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
	
		if host != "" {
			fmt.Fprint(os.Stdout, "--- "+host+" ---\n")
		}
		fmt.Fprintf(os.Stdout, "%.0f pings transmitted, %.0f received, %2.0f%% ping loss, time %s\n", transmitted, received, loss, duration)
		if received != 0 {
			avgT /= received
			fmt.Fprintf(os.Stdout, "rtt min/avg/max = %2.f/%.2f/%.2f ms\n", minT, avgT, maxT)
		}
		os.Exit(0)
	}

	mu := new(sync.Mutex)
	ping := func() {
		mu.Lock()
		msInt, _, err = c.RouterModule_pingNode(path, 0)
		transmitted++
		ms = float32(msInt)

		if err != nil {
			fmt.Fprintf(os.Stdout, "error: %s\n", err)
		} else {
			received++
			fmt.Fprintf(os.Stdout, "Reply from %v req=%v time=%03v ms\n",
				path, transmitted, ms)
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

	fmt.Fprintf(os.Stdout, "PING %s (%s)\n", host, ip)
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
