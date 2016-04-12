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

//String returns the full IRC command string
func (msg Message) String() string {
	return msg.Message
}

//NewMessage takes a string representing a command and parses it
func NewMessage(msg string) Message {
	return ParseString(msg)
}

//UserMessage returns a parsed User message 
func UserMessage(username, addr, servername, realname string) Message {
	return NewMessage(fmt.Sprintf("USER %s %s %s %s", username, addr, servername, realname))
}

//NickMessage returns a parsed Nick message
func NickMessage(nick string) Message {
	return NewMessage("NICK " + nick)
}

//PrivMessage returned a parsed PRIVMSG command
func PrivMessage(channel, msg string) Message {
	return NewMessage(fmt.Sprintf("PRIVMSG %s :%s", channel, msg))
}

func JoinMessage(channel string) Message {
    return NewMessage("JOIN "+channel)
}