package main

import (
	"fmt"
	"github.com/3M3RY/go-cjdns/cjdns"
	"github.com/spf13/cobra"
	"os"
)

var (
	admin   *cjdns.Admin
	rootCmd = &cobra.Command{Use: os.Args[0]}
)

func main() {
	rootCmd.AddCommand(PingCmd)

	var err error
	if admin, err = cjdns.NewAdmin(nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd.Execute()
}
