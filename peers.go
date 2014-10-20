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
	"os"
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

		fmt.Fprintf(os.Stdout, "%-39s %s\n"+
			"\tIncoming: %-5t      State: %s \n"+
			"\tBytes In:  %10d (%d%%)\n"+
			"\tBytes Out: %10d (%d%%)\n"+
			"\tIn/Out: %s  Lost Packets: %d\n\n",
			// Last seen: %s\n",

			ip, host,
			node.IsIncoming, node.State,
			node.BytesIn, (node.BytesIn * 100 / tIn),
			node.BytesOut, (node.BytesOut * 100 / tOut),
			ratio(node.BytesIn, node.BytesOut), node.LostPackets,
		// time.Duration(node.Last),
		)
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
	table, err := c.NodeStore_dumpTable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get routing table:", err)
		os.Exit(1)
	}
	peers := table.Peers(ip)
	if len(peers) == 0 {
		fmt.Fprintln(os.Stderr, "no peers found in local routing table")
		os.Exit(1)
	}
	peers.SortByQuality()

	for _, p := range peers {
		host, _, _ := resolve(p.IP.String())
		//fmt.Printf("\t%-39s %s\n", p.IP, host)
		fmt.Printf("%-39s %s Link: %s %s\n", p.IP, p.Path, p.Link, host)
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
	return fmt.Sprintf("%d/%d", in, out)
}
