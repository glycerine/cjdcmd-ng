package main

import (
	"fmt"
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
	if len(args) == 0 {
		stats, err := Admin.InterfaceController_peerStats()
		if err != nil {
			fmt.Println("Error getting local peers,", err)
		}

		var addr string
		for _, stat := range stats {
			addr = stat.PublicKey.IP().String()
			host, _, _ := resolve(addr)
			fmt.Println("\t ", stat.PublicKey, addr, host)
		}
		return
	}

	if len(args) > 1 {
		cmd.Usage()
		os.Exit(1)
	}
	_, ip, err := resolve(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve "+args[0]+".")
		os.Exit(1)
	}

	table, err := Admin.NodeStore_dumpTable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get routing table:", err)
		os.Exit(1)
	}

	peers := table.Peers(ip)

	if len(peers) == 0 {
		fmt.Println("no peers found in local routing table\n")
		os.Exit(1)
	}
	for _, p := range peers {
		host, _, _ := resolve(p.IP.String())
		fmt.Println("\t ", p.IP, host)
	}
}
