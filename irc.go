package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

var (
	HypIRCAddrs = []string{
		"[fcbf:7bbc:32e4:0716:bd00:e936:c927:fc14]:6667",
		"[fc13:6176:aaca:8c7f:9f55:924f:26b3:4b14]:6667",
		"[fca8:2dd7:4987:a9be:c8fc:34d7:05a1:4606]:6667",
		"[fcef:c7a9:792a:45b3:741f:59aa:9adf:4081]:6667",
	}
)

const (
	RPL_WHOISUSER     = "311" // "<nick> <user> <host> * :<real name>"
	RPL_WHOISSERVER   = "312" // "<nick> <server> :<server info>"
	RPL_WHOISOPERATOR = "313" // "<nick> :is an IRC operator"
	RPL_WHOISIDLE     = "317" // "<nick> <integer> :seconds idle"
	RPL_ENDOFWHOIS    = "318" // "<nick> :End of WHOIS list"
	RPL_WHOISCHANNELS = "319" // "<nick> :*( ( "@" / "+" ) <channel> " " )"
	RPL_NAMREPLY      = "353"
	RPL_ENDOFNAMES    = "3666"
)

func nickCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}
	_, ip, err := resolve(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve "+args[0]+".")
		os.Exit(1)
	}
	addr := ip.String()

	rand.Seed(time.Now().UnixNano())

	var conn net.Conn
	for {
		server := HypIRCAddrs[rand.Int31()%4]
		if Verbose {
			fmt.Fprintln(os.Stderr, "connecting to", server)
		}
		if conn, err = net.Dial("tcp6", server); err == nil {
			break
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect to any HypeIRC servers,", err)
		os.Exit(1)
	}

	infoMap := map[string]*ircInfo{
		addr: nil,
	}

	client := newIrcClient(conn)
	client.FindUsers(infoMap)

	if info := infoMap[addr]; info != nil {
		info.name = strings.TrimSuffix(info.name, "]")
		info.server = strings.TrimSuffix(strings.TrimPrefix(info.server, "["), "]")
		info.idle = strings.TrimSuffix(strings.TrimPrefix(info.idle, "["), "]")
		info.channels = strings.TrimSuffix(strings.TrimPrefix(info.channels, "["), "]")

		fmt.Println("\t nick:     ", info.nick)
		fmt.Println("\t name:     ", info.name)
		fmt.Println("\t server:   ", info.server)
		fmt.Println("\t idle:     ", info.idle)
		fmt.Println("\t channels: ", info.channels)
		if info.op {
			fmt.Println("\t ", info.nick, "is an IRC operator")
		}
	} else {
		fmt.Println("not found on HyperIRC")
	}
}

type ircClient struct {
	conn          net.Conn
	w             *textproto.Writer
	namesReplies  chan []string
	whoisRequests chan string
	whoisReplies  chan []string
}

func newIrcClient(conn net.Conn) *ircClient {
	c := &ircClient{
		conn,
		textproto.NewWriter(bufio.NewWriter(conn)),
		make(chan []string, 1024),
		make(chan string, 1),
		make(chan []string, 6),
	}

	go func() {
		var line string
		var fields []string
		var err error
		r := textproto.NewReader(bufio.NewReader(conn))
		for {
			if line, err = r.ReadLine(); err == nil {
				fields = strings.Fields(line)[1:]
				if len(fields) > 3 {
					switch fields[0] {

					case RPL_NAMREPLY:
						c.namesReplies <- fields[2:]

					case RPL_WHOISUSER, RPL_WHOISSERVER, RPL_WHOISOPERATOR, RPL_WHOISIDLE, RPL_ENDOFWHOIS, RPL_WHOISCHANNELS:
						c.whoisReplies <- fields

					case RPL_ENDOFNAMES:

					default:
						if Verbose {
							fmt.Fprintln(os.Stderr, "\t ", line)
						}
					}
				} else if len(fields) >= 2 && fields[0] == "PING" {
					c.w.PrintfLine("PONG %s", fields[1])
				} else {
					if Verbose {
						fmt.Fprintf(os.Stderr, "got some sort of short unhandled msg %q\n", line)
					}
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
				conn.Close()
				break
			}
		}
	}()

	c.w.PrintfLine("NICK fc%x", rand.Int31())
	c.w.PrintfLine("USER %x 8 * :%x", rand.Int31(), rand.Int31())
	return c
}

func (c *ircClient) FindUsers(infoMap map[string]*ircInfo) {
	go c.parseNames()
	c.w.PrintfLine("NAMES")

	c.whoisNames(infoMap)

	c.w.PrintfLine("QUIT :")
	c.conn.Close()
}

func (c *ircClient) parseNames() {
	// map to cache users we've queried
	ircUsers := make(map[string]bool)

	for fields := range c.namesReplies {
		if fields[0] != "=" {
			close(c.whoisRequests)
			return
		}
		for _, nick := range fields[2:] {
			for C := nick[0]; C == ':' || C == '@' || C == '+'; C = nick[0] {
				nick = nick[1:]
			}

			if !ircUsers[nick] {
				ircUsers[nick] = true

				c.whoisRequests <- nick
			}
		}
	}
}

func (c *ircClient) whoisNames(infoMap map[string]*ircInfo) {
	for nick := range c.whoisRequests {
		c.w.PrintfLine("WHOIS %s", nick)

		if end := c.readWhois(infoMap); end {
			return
		}
	}
}

func (c *ircClient) readWhois(infoMap map[string]*ircInfo) (end bool) {
	var info *ircInfo

	for rpl := range c.whoisReplies {

		switch rpl[0] {
		case RPL_WHOISUSER:
			addr := rpl[4]
			if ipRegex.MatchString(addr) {
				if _, ok := infoMap[addr]; ok {
					info = new(ircInfo)
					infoMap[addr] = info
				}
			}
			if info != nil {
				// might need to check the len on the rest
				info.nick = rpl[3]
				info.name = fmt.Sprint(rpl[7:])[1:] // concatenate the rest and strip the leading ':'
			}

		case RPL_WHOISSERVER:
			if info != nil {
				info.server = fmt.Sprint(rpl[3:])
			}

		case RPL_WHOISIDLE:
			if info != nil {
				info.idle = fmt.Sprint(rpl[3:])
			}

		case RPL_WHOISCHANNELS:
			if info != nil {
				info.channels = fmt.Sprint(rpl[3:])
			}

		case RPL_WHOISOPERATOR:
			if info != nil {
				info.op = true
			}

		case RPL_ENDOFWHOIS:
			bail := true
			if info != nil {
				// if all infos are present in the map, bail out
				for _, info := range infoMap {
					if info == nil {
						bail = false
						break
					}
				}
				if bail {
					return true
				}
				info = nil
			}
			return false

		default:
			fmt.Fprintf(os.Stderr, "the following WHOIS message is unhandled:%q\n", rpl)
		}
	}
	return true
}

type ircInfo struct {
	nick, name, server, idle, channels string
	op                                 bool
}
