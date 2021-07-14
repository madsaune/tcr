# Twitch Chat Reader

Read Twitch channel message in your terminal.

This project was created as a learning experience of how to setup a TCP client in Go. There are probably many similar programs that work way better.

## TODO

- [ ] Support TLS connection (port 6697)
- [x] Add color to usernames
- [x] Skip info messages by default. Add switch to show them.
- [ ] Add more colors to pick from

## Retrieve a token

source: [https://dev.twitch.tv/docs/irc/guide](https://dev.twitch.tv/docs/irc/guide)

To quickly get a token for your account, use this [Twitch Chat OAuth Password Generator](https://twitchapps.com/tmi/).

Then set your environment variables:

```bash
export TCR_USERNAME=<username>
export TCR_TOKEN=oauth:<token>
```

## How to run

First you must get a OAuth token as described above, and then set your ENV variables.

```bash
go get -u github.com/madsaune/tcr
tcr esl_csgo # esl_csgo is the channel name
```
