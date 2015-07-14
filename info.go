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

func infoCmd(cmd *cobra.Command, args []string) {
	c := Connect()

	if len(args) == 0 {
		self, err := c.NodeStore_nodeForAddr("")
		if err != nil {
			die(err.Error())
		}

		fmt.Fprintf(os.Stdout, self.Key +
			"\n\tProtocol version: %d"+
			"\n\tLink Count: %d\n",
			self.ProtocolVersion, self.LinkCount,
		)
		return
	}

	for _, host := range args {
		_, ip, err := resolve(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		node, err := c.NodeStore_nodeForAddr(ip)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fmt.Fprintf(os.Stdout,
			"%s\n"+
				"\tKey: %s\n"+
				"\tProtocol version: %d\n"+
				"\tBest Parent: %s\n"+
				"\tLink Count: %d\n"+
				"\tReach: %d\n\n",
			host, node.Key, node.ProtocolVersion, node.BestParent.IP, node.LinkCount, node.Reach,
		)
	}
}
