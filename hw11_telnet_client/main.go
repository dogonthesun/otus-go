package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 0, "connecting timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(
			os.Stderr,
			"Usage of %[1]s: %[1]s [flags] hostname port\nFlags:\n",
			os.Args[0],
		)
		flag.PrintDefaults()
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])

	os.Exit(telnet(address, timeout))
}

func telnet(address string, timeout time.Duration) int {
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	}
	defer client.Close()
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	receiveErrChan, sendErrChan := make(chan error), make(chan error)

	go func() {
		receiveErrChan <- client.Receive()
	}()

	go func() {
		sendErrChan <- client.Send()
	}()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	select {
	case err := <-receiveErrChan:
		if err == nil {
			fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	case err := <-sendErrChan:
		if err == nil {
			fmt.Fprintf(os.Stderr, "...EOF\n")
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	case <-ctx.Done():
	}

	return 0
}
