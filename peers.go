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

	"github.com/ehmry/go-cjdns/key"
	"github.com/spf13/cobra"
)

func showLocalPeers() {
	c := Connect()
	stats, err := c.InterfaceController_peerStats()
	if err != nil {
		fmt.Println("Error getting local peers,", err)
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

			fmt.Fprintf(os.Stdout, "%s %s\n"+
				"\tIncoming: %-5t\n"+
				"\tState: %s \n"+
				"\tBytes In:  %10d (%d%%)\n"+
				"\tBytes Out: %10d (%d%%)\n"+
				"\tTraffic Ratio: %s\n"+
				"\tLost Packets: %d\n\n",
				// Last seen: %s\n",

				node.PublicKey, host,
				node.IsIncoming, node.State,
				tIn, inP, tOut, outP,
				ratio(node.BytesIn, node.BytesOut), node.LostPackets,
			// time.Duration(node.Last),
			)
		} else {
			fmt.Fprintln(os.Stdout, ip)
		}
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

	for out%2 == 0 {
		out /= 2
		in /= 2
	}
	return fmt.Sprintf("↓%d/%d↑", in, out)
}
