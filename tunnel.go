package main

// These functions are not enabled, as tunneling is not well supported by cjdns.

import (
	"fmt"
	"github.com/inhies/go-cjdns/key"
	"github.com/spf13/cobra"
	"net"
	"os"
)

/*
var (
	TunnelCmd = &cobra.Command{
		Use: "tunnel",
		Run: tunnelCmd,
	}
	TunnelAllowCmd = &cobra.Command{
		Use: "allow PUBLIC_KEY TUNNEL_ADDRESS",
		Run: tunnelAllowCmd,
	}
	TunnelConnectCmd = &cobra.Command{
		Use: "connect PUBLIC_KEY",
		Run: tunnelConnectCmd,
	}
	TunnelDisconnectCmd = &cobra.Command{
		Use: "disconnect PUBLIC_KEY",
		Run: tunnelDisconnectCmd,
	}
)


func init() {
	TunnelCmd.AddCommand(
		TunnelAllowCmd,
		TunnelConnectCmd,
		TunnelDisconnectCmd)
}
*/

func tunnelCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		c := Connect()

		indexes, err := c.IpTunnel_listConnections()
		if err != nil {
			fmt.Println("Error getting tunnel information,", err)
			os.Exit(1)
		}
		for _, i := range indexes {
			info, err := c.IpTunnel_showConnection(i)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var f string
			if info.Outgoing {
				f = "%s(%s) - %s - Outgoing\n"
			} else {
				f = "%s(%s) - %s\n"
			}
			fmt.Printf(f, info.Key, info.Ip6Address, info.Ip6Address)
		}
	}
}

func tunnelAllowCmd(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.Usage()
		os.Exit(1)
	}
	key, err := key.DecodePublic(args[0])
	if err != nil {
		fmt.Println("Error reading client public key,", err)
		os.Exit(1)
	}
	addr := net.ParseIP(args[1])

	c := Connect()

	err = c.IpTunnel_allowConnection(key, addr)
	if err != nil {
		fmt.Println("Error allowing tunnel,", err)
		os.Exit(1)
	}
}

func tunnelConnectCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	key, err := key.DecodePublic(args[0])
	if err != nil {
		fmt.Println("Error reading server public key,", err)
		os.Exit(1)
	}

	c := Connect()

	err = c.IpTunnel_connectTo(key)
	if err != nil {
		fmt.Printf("Errror connecting to %s, %s", key.IP(), err)
		os.Exit(1)
	}
}

func tunnelDisconnectCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}

	key, err := key.DecodePublic(args[0])
	if err != nil {
		fmt.Println("Error reading server public key,", err)
	}

	c := Connect()

	indexes, err := c.IpTunnel_listConnections()
	if err != nil {
		fmt.Println("Error getting tunnel information,", err)
		os.Exit(1)
	}
	for _, i := range indexes {
		info, err := c.IpTunnel_showConnection(i)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if info.Key.Equal(key) {
			if err = c.IpTunnel_removeConnection(i); err != nil {
				fmt.Println("Failed to remove tunnel,", err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Println("tunnel not found")
}
