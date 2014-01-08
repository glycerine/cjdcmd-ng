package main

import (
	"fmt"
	"github.com/inhies/go-cjdns/admin"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	LogCmd = &cobra.Command{
		Use: "log",
		Run: log,
	}

	level, file string
	line        int
)

func init() {
	LogCmd.Flags().StringVarP(&level, "level", "l", "", `log level:
	KEYS - Not compiled in by default, contains private keys and other secret information.
	DEBUG    - Default level, contains lots of information which is probably not useful unless you are diagnosing an ongoing problem.
	INFO     - Shows starting and stopping of various components and general purpose information.
	WARN     - Generally this means some system has undergone a minor failure, this includes failures due to network disturbance.
	ERROR    - This means there was a (possibly temporary) failure of a system within cjdns.
	CRITICAL - This means something is broken such that the cjdns core will likely have to exit immedietly.
`)
	LogCmd.Flags().StringVarP(&file, "file", "f", "", `The name of the file where the log message came from, eg: "CryptoAuth.c".`)
	LogCmd.Flags().IntVarP(&line, "line", "L", 0, `The line number of the line where the log function was called.`)
}

func log(cmd *cobra.Command, args []string) {
	logChan := make(chan *admin.LogMessage, 32)

	_, err := Admin.AdminLog_subscribe(level, file, line, logChan)
	if err != nil {
		fmt.Println("Error subscribing to logging,", err)
		os.Exit(1)
	}
	for msg := range logChan {
		fmt.Println(time.Unix(msg.Time, 0), msg.Level, msg.File, msg.Message)
	}
}
