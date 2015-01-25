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
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func connectCmd(cmd *cobra.Command, args []string) {
	var (
		pubKey, addr, pass string
		face               int
		err                error
	)

	switch len(args) {
	case 4:
		pubKey, addr, pass = args[0], args[1], args[2]
		face, err = strconv.Atoi(args[3])
		if err != nil {
			die(err.Error())
		}

	case 3:
		pubKey, addr, pass = args[0], args[1], args[2]

	case 2:
		pubKey, addr = args[0], args[1]

		fmt.Fprintf(os.Stderr, "Password: ")
		pass, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			die(err.Error())
		}

	default:
		cmd.Usage()
		os.Exit(1)
	}

	// Check the address, resolve if nesisary.
	s := strings.SplitN(addr, ":", 2)
	if len(s) != 2 {
		die("Invalid address %s", addr)
	}
	if net.ParseIP(s[0]) == nil {
		ips, err := net.LookupHost(s[0])
		if err != nil {
			die("Could not resolve %q, %s\n", s[0], err)
		}
		if len(ips) < 1 {
			die("%s not found", s[0])
		}
		addr = fmt.Sprintf("%s:%s", ips[0], s[1])
	}

	// Connect to cjdns.
	c := Connect()
	fmt.Println("connecting to", addr)
	err = c.UDPInterface_beginConnection(pubKey, addr, face, pass)
	if err != nil {
		die(err.Error())
	}
}
