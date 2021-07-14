package main

import (
	"fmt"
	"strings"
	"time"
)

type Color string

const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)

// Message is the PRIVMSG sent from server
type Message struct {
	Username  string
	MsgType   string
	Channel   string
	Content   string
	Timestamp time.Time
}

// Parse parses the PRIVMSG body
func (m *Message) Parse(body string) {

	parts := strings.Split(body, " ")

	m.Username = m.parseUsername(parts[0])
	m.MsgType = parts[1]
	m.Channel = parts[2]
	m.Content = strings.Join(parts[3:], " ")[1:]
	m.Timestamp = time.Now()

}

func (m *Message) Colorize(color Color, message string) string {
	return fmt.Sprintf("%s%s%s", string(color), message, string(ColorReset))
	// return fmt.Sprintf("%s", message)
}

func (m *Message) parseUsername(username string) string {

	parts := strings.Split(username, "!")
	u := parts[0]

	return u[1:]
}
