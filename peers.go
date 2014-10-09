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

func peersCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		c := Connect()
		stats, err := c.InterfaceController_peerStats()
		if err != nil {
			fmt.Println("Error getting local peers,", err)
		}

		//var host addr string
		for _, node := range stats {
			ip := node.PublicKey.IP().String()
			host, _ := resolveIP(ip)

			fmt.Fprintf(os.Stdout, "%-39s %s %s\n"+
				"\tIncoming: %-5t      State: %s \n"+
				"\tBytes In: %-10d Bytes Out: %d\n"+
				"\tLost Packets: %d\n\n", // Last seen: %s\n",

				ip, node.SwitchLabel, host,
				node.IsIncoming, node.State,
				node.BytesIn, node.BytesOut,
				node.LostPackets, // time.Duration(node.Last),
			)
		}
		return
	}

	if len(args) > 1 {
		cmd.Usage()
		os.Exit(1)
	}
}
