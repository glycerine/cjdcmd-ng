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

	"github.com/inhies/go-bytesize"
	"github.com/spf13/cobra"
)

func trafficCmd(cmd *cobra.Command, args []string) {
	c := Connect()
	stats, err := c.InterfaceController_peerStats()
	if err != nil {
		die("Error getting peers stats,", err)
	}

	fmt.Fprint(os.Stdout, "Peer:                                     In:       Out:\n")

	for _, node := range stats {
		ip := node.PublicKey.IP().String()
		host, _ := resolveIP(ip)
		if host == "" {
			host = ip
		}
		fmt.Fprintf(os.Stdout, "%-39s  %8s  %8s\n",
			host,
			bytesize.New(float64(node.BytesIn)),
			bytesize.New(float64(node.BytesOut)))
	}

}
