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

	"github.com/inhies/go-cjdns/key"
	"github.com/spf13/cobra"
)

var (
	ConvertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert key forms",
		Run: func(cmd *cobra.Command, args []string) {
			ConvertPrivateCmd.Help()
			ConvertPublicCmd.Help()
			os.Exit(1)
		},
	}

	ConvertPrivateCmd = &cobra.Command{
		Use:   "private",
		Short: "Convert a private key into a public key",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				k, err := key.DecodePrivate(arg)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to decode private key,", err)
					os.Exit(1)
				}
				fmt.Fprintln(os.Stdout, k.Pubkey())
			}
		},
	}
	ConvertPublicCmd = &cobra.Command{
		Use:   "public",
		Short: "Convert a pubic key into an IP address",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				k, err := key.DecodePublic(arg)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Failed to decode public key,", err)
					os.Exit(1)
				}
				fmt.Fprintln(os.Stdout, k.IP())
			}
		},
	}
)

func init() {
	ConvertCmd.AddCommand(
		ConvertPrivateCmd,
		ConvertPublicCmd,
	)
}
