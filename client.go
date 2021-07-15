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
	colors := [12]Color{
		ColorRed,
		ColorRedBright,
		ColorGreen,
		ColorGreenBright,
		ColorYellow,
		ColorYellowBright,
		ColorBlue,
		ColorBlueBright,
		ColorMagenta,
		ColorMagentaBright,
		ColorCyan,
		ColorCyanBright,
	}

	tp := textproto.NewReader(bufio.NewReader(c.conn))

	for {
		status, err := tp.ReadLine()
		if err != nil {
			fmt.Println(err)
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

			// TODO: Check width of terminal and push wrapped lines
			// to fit the first line.
			//
			// Resources:
			//	- https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
			//  - https://pkg.go.dev/golang.org/x/term
			//
			// width, height, err := term.GetSize(0)
			//
			// Example:
			// [10:37:37]        tommyluco | this is also why the finding waldo games are
			//                             | so "hard" because you cant refer to memory
			//
			//

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
