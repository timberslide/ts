package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/timberslide/gotimberslide"
)

var (
	configFile = ".timberslide/config"
)

// usage is displayed if the user tries to do something invalid
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s <action> [<topic>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Actions: send, get\n")
	flag.PrintDefaults()
}

func main() {
	var err error

	client, err := ts.NewClient(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	flag.Usage = usage
	flag.Parse()
	// We expect 2 arguments: an action and a topic
	if flag.Arg(1) == "" {
		usage()
		os.Exit(1)
	}
	action := flag.Arg(0)
	topic := flag.Arg(1)

	switch action {
	case "send":
		err = client.Send(topic)
	case "get":
		err = client.Get(topic)
	default:
		err = errors.New(fmt.Sprintln("Unknown command:", flag.Arg(0)))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
