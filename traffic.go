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

	"github.com/glycerine/go-cjdns/key"
	"github.com/inhies/go-bytesize"
	"github.com/spf13/cobra"
)

func trafficCmd(cmd *cobra.Command, args []string) {
	c := Connect()
	stats, err := c.InterfaceController_peerStats()
	if err != nil {
		die("Error getting peers stats,", err)
	}

	var tIn, tOut float64

	for _, node := range stats {
		tIn += float64(node.BytesIn)
		tOut += float64(node.BytesOut)
	}

	fmt.Fprintf(os.Stdout,
		"                                         In:               Out:\n"+
			"                                         %8s          %8s\n\n",
		bytesize.New(float64(tIn)),
		bytesize.New(float64(tOut)))

	for _, node := range stats {
		// remove the 24 byte "v20.0000.0000.0000.0013." prefix
		// so that the Addr will decode to a public key.
		node.PublicKey, _ = key.DecodePublic(node.Addr[24:])
		ip := node.PublicKey.IP().String()
		host, _ := resolveIP(ip)
		if host == "" {
			host = ip
		}
		in := float64(node.BytesIn)
		out := float64(node.BytesOut)

		fmt.Fprintf(os.Stdout, "%-39s  %8s(%5.2f%%)  %8s(%5.2f%%)\n",
			host,
			bytesize.New(in), (100 * (in / tIn)),
			bytesize.New(out), (100 * (out / tOut)))
	}

}
