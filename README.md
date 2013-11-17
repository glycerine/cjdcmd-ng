cjdcmd-ng
======

Cjdcmd-ng is a command line tool for interfacing with [cjdns](https://github.com/cjdelisle/cjdns), a mesh network routing engine designed for security, scalability, speed, and ease of use. Its intent is to allow easy debugging of node and network problems as well as make it easier to work with the cjdns program itself.

Cjdcmd-ng is licensed under the GPL version 3 license, the full text of which is
available in `GPLv3.md`.

What's New
----------

Installation
------------

cjdcmd-ng is written in Go. If you already have Go installed, skip the next section.

### Install Go

To install go, check out the [install instructions](http://golang.org/doc/install) or, if you're lucky, you can do a quick install with the information below:

#### Ubuntu (and other Debian variants)

Ubuntu and other Debian variants do not have go 1.1 in their repositories. The easiest way to get it is to use `godeb`. Follow the instructions [here](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go), or:

    # 64 bit:
    wget https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz
    # or 32 bit:
    wget https://godeb.s3.amazonaws.com/godeb-386.tar.gz
    
    # untar it:
    tar xzf godeb-*.tar.gz
    
    # execute it:
    ./godeb install
    
    # Clean up:
    rm godeb-*.tar.gz godeb

#### Mac OSX

If you have [Homebrew](http://mxcl.github.com/homebrew/) installed:

    brew install go

If your system isn't listed above, perhaps you can [download pre-compiled binaries](http://code.google.com/p/go/downloads), otherwise you will have to [install from source](http://golang.org/doc/install/source).

#### Configure your folders (optional)

Next you should set up a special directory for all your go sourcecode and compiled programs. This is optional, however I recommend it because this will prevent you from having to use `sudo` every time you need to install an updated package or program. I will be giving a shortened version of the information found on [the official site](http://golang.org/doc/code.html#tmp_2)

First, make the folder where you want everything to be stored. I use the /home/inhies/projects/go but you may use whatever you like:

    $ mkdir -p $HOME/projects/go 

Next we need to tell our system where to look for Go packages and compiled programs. If you changed the folder name in the previous example, then make sure you change them here. You will want to add this to your `~/.profile` file or similar. On Ubuntu 12.10, I had to add it to my `~/.bashrc`:

	export GOPATH=$HOME/projects/go
	export PATH=$PATH:$HOME/projects/go/bin

Now, to make the changes take effect immediately, you can either run `$ source ~/.bashrc` or just paste those two lines you added to the file on your command line. You should now be setup! Try typing `$ echo $GOPATH` and make sure you see the folder you specified. If you have any problems, try re-reading the [official documentation](http://golang.org/doc/code.html#tmp_2).

### Install cjdcmd-ng

Once you have Go installed, installing new programs and packages couldn't be easier. Simply run the following command to have cjdcmd-ng download, build, and install:

    go get github.com/3M3RY/cjdcmd-ng
	
If f you see no output from that command then everything worked with no errors. To verify that it was successful, run `cjdcmd` and see if it displays some information about the program. If it does you are done! cjdcmd-ng has been downloaded, compiled, and installed. You amy now use it by typing `cjdcmd`.
	
**NOTE:** You may have to be root (use `sudo`) to install Go and cjdcmd-ng.

Updating cjdcmd-ng
---------------

To update your install of cjdcmd-ng, simply run `go get -u github.com/3M3RY/cjdcmd-ng` and it will automatically update, build, and install it. Just like when you initially installed cjdcmd-ng, if you see no output from that command then everything worked with no errors.
