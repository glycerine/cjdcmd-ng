cjdcmd-ng
======

```
Usage: 
  cjdcmd [command]

Available Commands: 
  ping <IPv6/DNS>                            :: Preforms a cjdns ping to a specified address.
  route <IPv6/DNS/Path>                      :: Prints all routes to a specific node
  traceroute <IPv6/DNS/Path>                 :: Performs a traceroute on a specific node by pinging each known hop to the target on all known paths
  ip <cjdns public key>                      :: Converts a cjdns public key to its corresponding IPv6 address.
  peers [<IPv6/DNS/Path>]                    :: Displays a list of currently connected peers for a node, if no node is specified your peers are shown.
  host <IPv6/DNS>                            :: Returns a list of all known IP addresses for a specified hostname or the hostname for an address.
  cjdnsadmin <-file /path/to/cjdroute.conf>  :: Generates a .cjdnsadmin file in your home diectory using the specified cjdroute.conf as input
  addpeer '<json peer details>'              :: Adds the peer details to your config file
  addpass [password]                         :: Adds the password to your config file, or generates one and then adds that
  listpass                                   :: ALPHA FEATURE - List currently loaded peering passwords.
  cleanconfig [-file] [-outfile]             :: Strips all comments from the config file and saves it at outfile
  log [--level level] [--file file] [--line] :: Prints cjdns logs to stdout
  passgen [prefix]                           :: Generates a random alphanumeric password between 15 and 50 characters. If you provide [prefix], it will be prepended. This is to help you keep track of your peering passwords
  dump                                       :: Dumps the entire routing table to stdout.
  memory                                     :: Returns the bytes of memory allocated by the router
  nick <IPv6/DNS>                            :: Scrape HyperIRC for nicks using host
  help [command]                             :: Help about any command

 Available Flags:
  -r, --resolve=false: reverse resolve IP addresses
  -v, --verbose=false: verbose output

Use "cjdcmd help [command]" for more information about that command.
```

## Install from source

To install go, check out the [install instructions](http://golang.org/doc/install).

Run the following command to have cjdcmd-ng download, build, and install:

    go get github.com/ehmry/cjdcmd-ng
