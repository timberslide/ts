package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/timberslide/gotimberslide"
)

var (
	configFile = ".timberslide/config"
)

const (
	ErrConfig = iota // ErrConfig is returned if we have a proble with the config file
	ErrUsage         // ErrUsage is returned if we exit because we displayed the usage
	ErrServer        // ErrServer is returned for any error because of the server
)

// usage is displayed if the user tries to do something invalid
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s <action> [options...] <topic>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Actions: send, get\n")
	flag.PrintDefaults()
}

func main() {
	var err error

	// Configure the Timberslide client
	client, err := ts.NewClient(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ErrConfig)
	}
	flag.Usage = usage

	if len(os.Args) < 3 {
		usage()
		os.Exit(ErrUsage)
	}

	// Set up our command line flags
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	var all bool
	getCmd.BoolVar(&all, "all", false, "Begin from the oldest message")

	topic := os.Args[len(os.Args)-1]

	switch os.Args[1] {
	case "send":
		sendCmd.Parse(os.Args[2:])
		err = client.Send(topic)
	case "get":
		getCmd.Parse(os.Args[2:])
		position := ts.PositionNewest
		if all {
			position = ts.PositionOldest
		}
		err = client.Get(topic, position)
	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid command. Valid commands are `send` and `get`.\n", os.Args[1])
		os.Exit(ErrUsage)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(ErrServer)
	}
}
