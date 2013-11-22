package main

import (
	"fmt"
	"github.com/3M3RY/go-cjdns/cjdns"
	"github.com/3M3RY/go-nodeinfo"
	"github.com/spf13/cobra"
	"net"
	"os"
	"sync"
)

var (
	ConfFile        string
	NmapOutput      bool
	ResolveNodeinfo bool
	ResolveSystem   bool
	Verbose         bool
	Admin           *cjdns.Admin
)

var rootCmd = &cobra.Command{Use: os.Args[0]}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&ResolveNodeinfo, "nodeinfo", "N", false, "resolve hostnames with nodeinfo")
	rootCmd.PersistentFlags().BoolVarP(&ResolveSystem, "system-resolver", "S", true, "resolve hostnames with system resolver")
	TraceCmd.Flags().BoolVarP(&NmapOutput, "nmap", "x", false, "format output as nmap XML")
	//AddPassCmd.Flags().StringVarP(&ConfFile, "conf", "c", "", "path to cjdroute.conf")

	//rootCmd.AddCommand(AddPassCmd)
	rootCmd.AddCommand(GetNodeinfoCmd)
	rootCmd.AddCommand(SetNodeinfoCmd)
	rootCmd.AddCommand(NickCmd)
	rootCmd.AddCommand(PeersCmd)
	rootCmd.AddCommand(PingCmd)
	rootCmd.AddCommand(TraceCmd)

	var err error
	if Admin, err = cjdns.NewAdmin(nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd.Execute()
}

var (
	// BUG(emery): oft useless instantiation
	rCache = make(map[string]string)
	rMutex = new(sync.RWMutex)
)

func Resolve(addr string) string {
	rMutex.RLock()
	s, ok := rCache[addr]
	rMutex.RUnlock()
	if ok {
		return s
	}
	rMutex.Lock()
	defer rMutex.Unlock()
	s, ok = rCache[addr]
	if ok {
		return s
	}

	if ResolveNodeinfo {
		if hostnames, err := nodeinfo.LookupAddr(addr); err == nil {
			for _, h := range hostnames {
				s = s + h + " "
			}
		} else if Verbose {
			fmt.Fprintf(os.Stderr, "failed to resolve %s, %s\n", addr, err)
		}
	}

	if ResolveSystem {
		if hostnames, err := net.LookupAddr(addr); err == nil {
			for _, h := range hostnames {
				s = s + h + " "
			}
		} else if Verbose {
			fmt.Fprintf(os.Stderr, "failed to resolve %s, %s\n", addr, err)
		}
	}
	if len(s) > 0 {
		rCache[addr] = s
	}
	return s
}
