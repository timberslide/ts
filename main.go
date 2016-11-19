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
	// ErrConfig is returned if we have a proble with the config file
	ErrConfig = iota

	// ErrUsage is returned if we exit because we displayed the usage
	ErrUsage

	// ErrServer is returned for any error because of the server
	ErrServer

	// ErrSystem is returned when a problem with the user's system is encountered
	ErrSystem
)

// usage is displayed if the user tries to do something invalid
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s <action> [options...] <topic>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Actions: send, get\n")
	flag.PrintDefaults()
}

// Get displays all events to stdout
func Get(client ts.Client, topic string, position int64) error {
	for event := range client.Iter(topic, position) {
		fmt.Println(event.Message)
	}
	return nil
}

func main() {
	var err error

	// Configure the Timberslide client
	client, err := ts.NewClientFromFile(configFile)
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
	var quietFlag bool
	sendCmd.BoolVar(&quietFlag, "q", false, "Do not display the stream to stdout")

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	var allFlag bool
	getCmd.BoolVar(&allFlag, "all", false, "Begin from the oldest message")

	topic := os.Args[len(os.Args)-1]

	switch os.Args[1] {
	case "send":
		sendCmd.Parse(os.Args[2:])
		if quietFlag {
			os.Stdout, err = os.Open("/dev/null")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not redirect stdout to /dev/null")
				os.Exit(ErrSystem)
			}
		}
		err = client.Send(topic)
	case "get":
		getCmd.Parse(os.Args[2:])
		position := ts.PositionNewest
		if allFlag {
			position = ts.PositionOldest
		}
		err = Get(client, topic, position)
	default:
		fmt.Fprintln(os.Stderr, os.Args[1], "is not a valid command. Valid commands are `send` and `get`.")
		os.Exit(ErrUsage)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(ErrServer)
	}
}
