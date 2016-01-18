package irc

import (
	"fmt"
	"strings"
)

type Message string

func (msg Message) String() string {
	return string(msg)
}

func NewMessage(msg string) Message {
	if strings.HasSuffix(msg, "\r\n") {
		return Message(msg)
	}
	return Message(msg + "\r\n")
}

func UserMessage(username, addr, servername, realname string) Message {
	return NewMessage(fmt.Sprintf("USER %s %s %s %s", username, addr, servername, realname))
}

func NickMessage(nick string) Message {
	return NewMessage("NICK " + nick)
}

func PrivMessage(channel, msg string) Message {
	return NewMessage(fmt.Sprintf("PRIVMSG %s :%s", channel, msg))
}
