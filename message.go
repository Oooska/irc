package irc

import (
	"fmt"
)

/* A Message represents a message sent to or from the IRC server.

The various components of the message are parsed for convenience.

Psudo-BNF from: https://tools.ietf.org/html/rfc1459#section-2.3.1


<message>  ::= [':' <prefix> <SPACE> ] <command> <params> <crlf>
<prefix>   ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
<command>  ::= <letter> { <letter> } | <number> <number> <number>

<params>   ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
<middle>   ::= <Any *non-empty* sequence of octets not including SPACE
               or NUL or CR or LF, the first of which may not be ':'>
<trailing> ::= <Any, possibly *empty*, sequence of octets not including
                 NUL or CR or LF>
*/
type Message struct {
	Message  string   //The raw, unparsed message
    
	Prefix   string   //includes the ':' character
    Nick string 
    User string  
    Host string 
    Server string
    
    
	Command  string
	Params   []string //Includes Trailing as the final argument with the '
    Trailing string   //excludes the ':'

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