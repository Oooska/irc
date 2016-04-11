package irc

import (
	"fmt"
)

//Message contains the parsed components of an IRC message
type Message struct {
	Message  string
	Prefix   string //including :
	Command  string
	Params   []string 
    
    User User
    Server string
}

//Represents an IRC user
type User struct {
    //Parsed parts of prefix
    Nick string 
    User string 
    Host string 
}

func (msg Message) String() string {
	return msg.Message
}

func NewMessage(msg string) Message {
	return ParseString(msg)
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
