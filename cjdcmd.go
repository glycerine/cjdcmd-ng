package main

import (
	"fmt"
	"github.com/3M3RY/go-cjdns/cjdns"
	"github.com/3M3RY/go-nodeinfo"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

var (
	Verbose         bool
	ResolveNodeinfo bool
	Admin           *cjdns.Admin
)

var rootCmd = &cobra.Command{Use: os.Args[0]}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&ResolveNodeinfo, "nodeinfo", "N", false, "resolve hostnames with nodeinfo")

	rootCmd.AddCommand(GetNodeinfoCmd)
	rootCmd.AddCommand(SetNodeinfoCmd)
	rootCmd.AddCommand(NickCmd)
	rootCmd.AddCommand(PeersCmd)
	rootCmd.AddCommand(PingCmd)
	rootCmd.AddCommand(TraceRoutesCmd)

	var err error
	if Admin, err = cjdns.NewAdmin(nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rootCmd.Execute()
}

var (
	// BUG(emery): oft useless instantiation
	nodeinfoCache = make(map[string]string)
	nodeinfoMutex = new(sync.RWMutex)
)

// TODO actually make this a general reverse lookup fuction
// that can scrap IRC in addition to nodeinfo
// convience fuction for caching nodeinfo hostnames
func NodeinfoReverse(addr string) (hostname string) {
	nodeinfoMutex.RLock()
	var ok bool
	hostname, ok = nodeinfoCache[addr]
	nodeinfoMutex.RUnlock()
	if ok {
		return
	}
	nodeinfoMutex.Lock()
	defer nodeinfoMutex.Unlock()
	hostname, ok = nodeinfoCache[addr]
	if ok {
		return
	}

	if hosts, err := nodeinfo.LookupAddr(addr); err == nil {
		hostname = hosts[0]
		nodeinfoCache[addr] = hostname
	} else if Verbose {
		fmt.Fprintf(os.Stderr, "failed to resolve %s, %s\n", addr, err)
	}
	return
}
