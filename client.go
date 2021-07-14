package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"strings"
)

type Client struct {
	TwitchHost         string
	Username           string
	Token              string
	ShowStatusMessages bool

	conn net.Conn
}

func NewClient() *Client {
	return &Client{
		TwitchHost:         "irc.chat.twitch.tv:6667",
		ShowStatusMessages: false,
	}
}

func (c *Client) Connect(username, token string) {
	conn, err := net.Dial("tcp", c.TwitchHost)
	if err != nil {
		panic(err)
	}

	c.conn = conn
	c.Username = username
	c.Token = token

	c.Send("PASS " + c.Token)
	c.Send("NICK " + c.Username)
}

func (c *Client) Close() {
	c.Send("QUIT bye.")
	c.conn.Close()
}

func (c *Client) Join(channel string) {
	c.Send("JOIN " + channel)
}

func (c *Client) Pong() {
	c.Send("PONG :tmi.twitch.tv")
}

func (c *Client) Send(message string) {
	fmt.Fprintf(c.conn, "%s\r\n", message)
}

func (c *Client) Listen() {
	// Read whats sent from server
	var users = make(map[string]Color)
	colors := [5]Color{ColorBlack, ColorRed, ColorGreen, ColorYellow, ColorBlue}

	tp := textproto.NewReader(bufio.NewReader(c.conn))

	for {
		status, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}

		// If message starts with PING we must respond with PONG
		if strings.HasPrefix(status, "PING") {
			c.Pong()
		} else if strings.Contains(status, "PRIVMSG") {
			// If the body contains PRIVMSG, it is sent by a user
			// and can be parsed to look prettier
			msg := &Message{}
			msg.Parse(status)

			// Assign a random color for new users
			currentUser := msg.Username
			_, exists := users[currentUser]
			if !exists {
				users[currentUser] = colors[rand.Intn(len(colors))]
			}
			fmt.Printf("[%s] ", msg.Timestamp.Format("15:04:05"))
			fmt.Printf("%25s", msg.Colorize(users[currentUser], msg.Username))
			fmt.Printf(" | %s\n", msg.Content)
		} else {
			// Prints all messages that are not PING og PRIVMSG
			// like connection information etc
			if c.ShowStatusMessages {
				fmt.Println(status)
			}
		}
	}
}
