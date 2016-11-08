package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/BurntSushi/toml"
	"github.com/timberslide/gotimberslide"
)

var (
	configFile = ".timberslide/config"
)

// Client contains our client configuration
type Client struct {
	Host  string
	Token string
}

// usage is displayed if the user tries to do something invalid
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s <action> [<topic>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Actions: send, get\n")
	flag.PrintDefaults()
}

// newClient creates a new client from a default config file
func newClient() (Client, error) {
	var client Client

	usr, err := user.Current()
	if err != nil {
		return client, err
	}

	rcLoc := fmt.Sprintf("%s/%s", usr.HomeDir, configFile)
	b, err := ioutil.ReadFile(rcLoc)
	if err != nil {
		return client, fmt.Errorf("Could not read configuration at %s", rcLoc)
	}
	if _, err = toml.Decode(string(b), &client); err != nil {
		return client, err
	}

	return client, nil
}

// send takes stdin and sends it to the topic
func (c *Client) send(topic string) error {
	fmt.Fprintf(os.Stderr, "ts: connecting...\n")
	host, _, err := net.SplitHostPort(c.Host)
	if err != nil {
		return err
	}
	creds := credentials.NewClientTLSFromCert(nil, host)
	conn, err := grpc.Dial(c.Host, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := ts.NewIngestClient(conn)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.Token))
	ctx := metadata.NewContext(context.Background(), md)
	stream, err := client.StreamEvents(ctx)
	if err != nil {
		return err
	}

	msg := &ts.Event{Topic: topic}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg.Message = scanner.Text()
		// TODO make printing to stdout optional
		fmt.Println(msg.Message)
		stream.Send(msg)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// We finished normally, so...
	// set the done flag and block until we hear something from the server
	stream.Send(&ts.Event{Topic: topic, Done: true})
	fmt.Fprintf(os.Stderr, "ts: done, waiting for server...\n")
	stream.Recv()
	stream.CloseSend()

	return nil
}

// get receives the stream from the topic and writes it to stdout
func (c *Client) get(topic string) error {
	host, _, err := net.SplitHostPort(c.Host)
	if err != nil {
		return err
	}
	creds := credentials.NewClientTLSFromCert(nil, host)
	conn, err := grpc.Dial(c.Host, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	client := ts.NewStreamerClient(conn)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", c.Token))
	ctx := metadata.NewContext(context.Background(), md)
	stream, err := client.GetStream(ctx, &ts.Topic{Topic: topic})
	if err != nil {
		return err
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Println(event.Message)
	}
}

func main() {
	var err error

	client, err := newClient()
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
		err = client.send(topic)
	case "get":
		err = client.get(topic)
	default:
		err = errors.New(fmt.Sprintln("Unknown command:", flag.Arg(0)))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
