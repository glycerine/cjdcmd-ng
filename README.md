cjdcmd-ng
======

```
Usage: 
  ./cjdcmd-ng [command]

Available Commands: 
  getnodeinfo [ADDRESS/HOSTNAME] :: resolves a host or address through NodeInfo
  setnodeinfo HOSTNAME :: sets a NodeInfo hostname for this device
  nick ADDRESS/HOSTNAME :: resolves a host or address to an IRC nick
  peers HOST      :: prints a host's peers
  ping HOST       :: pings a host
  trace HOST [HOST...] :: prints routes to hosts
  help [command]  :: Help about any command
```




## Install from source

To install go, check out the [install instructions](http://golang.org/doc/install).

Run the following command to have cjdcmd-ng download, build, and install:

    go get github.com/3M3RY/cjdcmd-ng

## Install from binary

AMD64: http://urlcloud.net/NiKW
A signature is at cjdcmd-ng.amd64.xz.sig