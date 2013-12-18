package main

import (
	"errors"
	"fmt"
	"github.com/3M3RY/go-nodeinfo"
	"github.com/3M3RY/go-cjdns/admin"
	"github.com/spf13/cobra"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

var ipRegexp = regexp.MustCompile("^fc[a-f0-9]{1,2}:([a-f0-9]{0,4}:){2,6}[a-f0-9]{1,4}$")

var (
	ConfFile        string
	NmapOutput      bool
	ResolveNodeinfo bool
	ResolveSystem   bool
	Verbose         bool
	Admin           *admin.Conn
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
	if Admin, err = admin.Connect(nil); err != nil {
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

func resolve(host string) (hostname string, ip net.IP, err error) {
	if ipRegexp.MatchString(host) {
		ip = net.ParseIP(host)
		addr := ip.String()

		var ok bool
		rMutex.RLock()
		hostname, ok = rCache[addr]
		rMutex.RUnlock()
		if ok {
			return
		}
		rMutex.Lock()
		defer rMutex.Unlock()
		hostname, ok = rCache[addr]
		if ok {
			return
		}

		if ResolveNodeinfo {
			if hostnames, err := nodeinfo.LookupAddr(addr); err == nil {
				for _, h := range hostnames {
					hostname = hostname + h + " "
				}
			} else if Verbose {
				fmt.Fprintf(os.Stderr, "failed to resolve %s, %s\n", addr, err)
			}
		}

		if ResolveSystem {
			if hostnames, err := net.LookupAddr(addr); err == nil {
				for _, h := range hostnames {
					hostname = hostname + h + " "
				}
			} else if Verbose {
				fmt.Fprintf(os.Stderr, "failed to resolve %s, %s\n", addr, err)
			}
		}
		if len(hostname) > 0 {
			rCache[addr] = hostname
		}
		return
	} else {
		hostname = host
		var addrs []string
		if addrs, err = net.LookupHost(host); err != nil {
			return
		}

		for _, addr := range addrs {
			if ipRegexp.MatchString(addr) {
				ip = net.ParseIP(addr)
				return
			}
		}
		err = errors.New("no fc::/8 address found")
	}
	return
}

func padIPv6(truncated string) (full string) {
	if len(truncated) == 39 {
		return truncated
	}
	full = truncated[:4]
	for _, couplet := range strings.SplitN(truncated[5:], ":", 7) {
		if len(couplet) == 4 {
			full = full + ":" + couplet
		} else {
			full = full + fmt.Sprintf(":%04s", couplet)
		}
	}
	return
}
