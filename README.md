cjdcmd-ng
======

```
Usage: 
  cjdcmd-ng [command]

Available Commands: 
  getnodeinfo [ADDRESS/HOSTNAME] :: resolves a host or address through NodeInfo
  setnodeinfo HOSTNAME :: sets a NodeInfo hostname for this device
  nick ADDRESS/HOSTNAME :: resolves a host or address to an IRC nick
  peers HOST      :: prints a host's peers
  ping HOST       :: pings a host
  trace HOST [HOST...] :: prints routes to hosts
  tunnel          :: 
  log             :: 
  help [command]  :: Help about any command
```


### Install Go

To install go, check out the [install instructions](http://golang.org/doc/install) or, if you're lucky, you can do a quick install with the information below:

Run the following command to have cjdcmd-ng download, build, and install:

    go get github.com/3M3RY/cjdcmd-ng
	
Updating cjdcmd-ng
---------------

To update your install of cjdcmd-ng, simply run `go get -u github.com/3M3RY/cjdcmd-ng` and it will automatically update, build, and install it. Just like when you initially installed cjdcmd-ng, if you see no output from that command then everything worked with no errors.
