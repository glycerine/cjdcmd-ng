cjdcmd-ng
======

```
Available Commands: 
  ping        Preforms a cjdns ping to a specified address.
  route       Prints all routes to a specific node
  traceroute  Performs a traceroute on a specific node by pinging each known hop to the target on all known paths
  ip          Converts a cjdns public key to its corresponding IPv6 address.
  peers       Displays a list of currently connected peers for a node, if no node is specified your peers are shown.
  host        Returns a list of all known IP addresses for a specified hostname or the hostname for an address.
  cjdnsadmin  Generates a .cjdnsadmin file in your home diectory using the specified cjdroute.conf as input
  addpeer     Adds the peer details to your config file
  addpass     Adds the password to your config file, or generates one and then adds that
  listpass    ALPHA FEATURE - List currently loaded peering passwords.
  cleanconfig Strips all comments from the config file and saves it at outfile
  log         Prints cjdns logs to stdout
  passgen     Generates a random alphanumeric password between 15 and 50 characters. If you provide [prefix], it will be prepended. This is to help you keep track of your peering passwords
  dump        Dumps the entire routing table to stdout.
  nick        Scrape HyperIRC for nicks using host
  convert     Convert key forms
  connect     Connect directly to another node over UDP
  info        Show node information
  traffic     Show traffic statistics
  fingerprint Show public key unicode fingerprint
  help        Help about any command
```

## Install from source

To install go, check out the [install instructions](http://golang.org/doc/install).

Run the following command to have cjdcmd-ng download, build, and install:

    go get github.com/ehmry/cjdcmd-ng


## Configuration

You'll need a to create the file ''~/.cjdnsadmin'' so that cjdcmd-ng knows how
to connect to CJDNS.

```
{
        "addr": "127.0.0.1",
        "port": 11234,
        "password": "xxxxxxxxxxxxxxxxxxxxxxx",
        "config": "/etc/cjdroute.conf"
}
```

If you are using NixOS (like me) the password is stored in /etc/cjdns.keys.
