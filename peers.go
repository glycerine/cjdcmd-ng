package main

import (
	"fmt"
	"github.com/3M3RY/go-cjdns/cjdns"
	"github.com/spf13/cobra"
	"os"
)

var PeersCmd = &cobra.Command{
	Use:   "peers HOST",
	Short: "prints a host's peers",
	Long:  `Parses the CJDNS routing table and prints out nodes that are a single hop away from a given host.`,
	Run:   peers,
}

func peers(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}
	addr, err := cjdns.Resolve(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve "+args[0]+".")
		os.Exit(1)
	}

	table, err := Admin.NodeStore_dumpTable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get routing table:", err)
		os.Exit(1)
	}

	peers := table.Peers(addr)

	if len(peers) == 0 {
		fmt.Println("no peers found in local routing table\n")
		os.Exit(1)
	}
	if ResolveNodeinfo {
		for _, p := range peers {
			host := NodeinfoReverse(p.IP)
			fmt.Println("\t ", p.IP, host)
		}
	} else {
		for _, p := range peers {
			fmt.Println("\t ", p.IP)
		}
	}
}
