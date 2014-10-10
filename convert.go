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

var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert key forms",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
		for _, arg := range args {
			if len(arg) == 64 {
				k, err := key.DecodePrivate(arg)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to decode key,", err)
					os.Exit(1)
				}
				fmt.Fprintln(os.Stdout, k.Pubkey())
			} else {
				k, err := key.DecodePublic(arg)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to decode key,", err)
					os.Exit(1)
				}
				fmt.Fprintln(os.Stdout, k.IP())
			}
		}

	},
}
