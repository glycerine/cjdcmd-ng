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

	//"github.com/ehmry/go-cjdns/admin"
	"github.com/spf13/cobra"
)

func routeCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}
	c := Connect()

	var ip string
	var err error
	for _, host := range args {
		_, ip, err = resolve(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		node, err := c.NodeStore_nodeForAddr(ip)

		for err == nil && node.RouteLabel != "0000.0000.0000.0001" {
			if Verbose {
				fmt.Fprintf(os.Stdout, "%-39s - %s - v%d - Reach: %10d\n", ip, node.RouteLabel, node.ProtocolVersion, node.Reach)
			} else {
				fmt.Fprintln(os.Stdout, ip)
			}
			node, err = c.NodeStore_nodeForAddr(ip)
			ip = node.BestParent.IP
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Fprintln(os.Stdout)
	}
}
