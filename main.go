package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Channel is not specified")
		os.Exit(1)
	}

	// Get username and token from ENV
	tu, ok := os.LookupEnv("TCR_USERNAME")
	if !ok {
		fmt.Println("Env TCR_USERNAME is not present")
		os.Exit(1)
	}

	tt, ok := os.LookupEnv("TCR_TOKEN")
	if !ok {
		fmt.Println("Env TCR_TOKEN is not present")
		os.Exit(1)
	}

	// Initialize the Twitch IRC client
	client := NewClient()

	// Closes connection when CTRL+C is pressed
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nDisconnecting...")
		client.Close()
		fmt.Println("bye.")
		os.Exit(1)
	}()

	// Login with twitch username and oauth token
	client.Connect(tu, tt)

	// Join the twitch channel passed in from cmdline
	// and listen for messages
	client.Join(args[1])
	client.Listen()
}
