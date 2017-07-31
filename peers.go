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
	"os"
	"time"

	"github.com/glycerine/go-cjdns/key"
	"github.com/inhies/go-bytesize"
	"github.com/spf13/cobra"
)

func showLocalPeers() {
	c := Connect()
	stats, err := c.InterfaceController_peerStats()
	if err != nil {
		fmt.Println("Error getting local peers,", err)
	}
	for _, node := range stats {
		// remove the 24 byte "v20.0000.0000.0000.0013." prefix
		// so that the Addr will decode to a public key.
		node.PublicKey, err = key.DecodePublic(node.Addr[24:])
		if err != nil {
			fmt.Println("Error decoding peer Addr,", err)
		}
	}

	var tIn, tOut int64

	for _, node := range stats {
		tIn += node.BytesIn
		tOut += node.BytesOut
	}

	//var host addr string
	for _, node := range stats {
		ip := node.PublicKey.IP().String()
		host, _ := resolveIP(ip)

		if len(host) == 0 {
			host = ip
		}
		if Verbose {
			var inP, outP int64
			if tIn == 0 {
				inP = 0
			} else {
				inP = (node.BytesIn * 100) / tIn
			}
			if tOut == 0 {
				outP = 0
			} else {
				outP = (node.BytesOut * 100) / tOut
			}

			lastTm := time.Unix(node.Last/1000, (node.Last%1000)*1e6)
			fmt.Fprintf(os.Stdout, "%s %s\n"+
				"\tIncoming: %-5t\n"+
				"\tState: %s \n"+
				"\tBytes In:  %s (%d%%)\n"+
				"\tBytes Out: %s (%d%%)\n"+
				"\tTraffic Ratio: %s\n"+
				"\tLost Packets: %d\n"+
				"\tLast seen: %s (%s ago)\n\n",
				node.PublicKey, host,
				node.IsIncoming, node.State,
				bytesize.New(float64(node.BytesIn)), inP, bytesize.New(float64(node.BytesOut)), outP,
				ratio(node.BytesIn, node.BytesOut), node.LostPackets,
				// node.Last is a millisecond precision unix epoch timestamp
				lastTm, time.Since(lastTm),
			)
		} else {
			fmt.Fprintln(os.Stdout, ip)
		}
	}
	if Verbose {
		fmt.Fprintln(os.Stdout, "Total Traffic:", bytesize.New(float64(tIn)), ratio(tIn, tOut), bytesize.New(float64(tOut)))
	}
}

func peersCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		showLocalPeers()
		return
	}

	if len(args) > 1 {
		cmd.Usage()
		os.Exit(1)
	}

	_, ip, err := resolve(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not resolve %q.\n", args[0])
		os.Exit(1)
	}

	c := Connect()

	node, err := c.NodeStore_nodeForAddr(ip)
	if err != nil {
		die(err.Error())
	}

	peers, _, err := c.RouterModule_getPeers(node.RouteLabel, 0, "")
	if err != nil {
		die(err.Error())
	}
	for _, s := range peers {
		k, err := key.DecodePublic(s[24:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "received malformed key ", s[24:])
			continue
		}
		if Verbose {
			fmt.Fprintln(os.Stdout, s[:3], s[3:24], s[24:], k.IP())
		} else {
			fmt.Fprintln(os.Stdout, k.IP())
		}
	}
}

const maxSpread = 32

func ratio(in, out int64) string {
	if in == 0 || out == 0 {
		return "∞"
	}

	var factor int64
	if in > out {
		factor = in / maxSpread
	} else if out > in {
		factor = out / maxSpread
	} else {
		return "1/1"
	}

	out /= factor
	in /= factor

	for out%2 == 0 && in%2 == 0 {
		out /= 2
		in /= 2
	}
	return fmt.Sprintf("↓%d/%d↑", in, out)
}
