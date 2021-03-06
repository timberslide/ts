package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

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

// SendStdin pipes stdin into Timberslide
func SendStdin(client ts.Client, topic string) error {
	var err error
	ch, err := client.CreateChannel(topic)
	if err != nil {
		return err
	}
	defer ch.Close()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		ch.Send(message)
		fmt.Println(message)
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Stream done, closing connection")
	return err
}

// Get displays all events to stdout
func Get(client ts.Client, topic string, position int64) error {
	for event := range client.Iter(topic, position) {
		fmt.Println(strings.TrimRightFunc(event.Message, unicode.IsSpace))
	}
	return nil
}

// List gets a list of topics and displays them to the user
func List(client ts.Client) error {
	topics, err := client.GetTopics()
	if err != nil {
		return err
	}
	for _, topic := range topics {
		fmt.Println(topic.Name)
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

	if len(os.Args) < 2 {
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

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	topic := os.Args[len(os.Args)-1]

	fmt.Fprintln(os.Stderr, "Connecting to Timberslide")
	err = client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ErrServer)
	}
	defer client.Close()

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
		err = SendStdin(client, topic)
	case "get":
		getCmd.Parse(os.Args[2:])
		position := ts.PositionNewest
		if allFlag {
			position = ts.PositionOldest
		}
		err = Get(client, topic, position)
	case "list":
		listCmd.Parse(os.Args[2:])
		err = List(client)
	default:
		fmt.Fprintln(os.Stderr, os.Args[1], "is not a valid command. Valid commands are `send` and `get`.")
		os.Exit(ErrUsage)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(ErrServer)
	}
}
