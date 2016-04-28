package irc

import (
	"fmt"
	"time"
)

//NewMessage takes a string representing a command and parses it
//Timestamp is set to time.Now()
func NewMessage(msg string) Message {
	pmsg := parseString(msg)
	pmsg.timestamp = time.Now()
	return pmsg
}

//MessageWithTimestamp returns a Message with the specified timestamp.
func MessageWithTimestamp(msg string, ts time.Time) Message {
	pmsg := parseString(msg)
	pmsg.timestamp = ts
	return pmsg
}

//UserMessage returns a parsed User message
func UserMessage(username, addr, servername, realname string) Message {
	return NewMessage(fmt.Sprintf("USER %s %s %s %s", username, addr, servername, realname))
}

//NickMessage returns a parsed Nick message
func NickMessage(nick string) Message {
	return NewMessage("NICK " + nick)
}

//PrivMessage returns a parsed PRIVMSG command
func PrivMessage(channel, msg string) Message {
	return NewMessage(fmt.Sprintf("PRIVMSG %s :%s", channel, msg))
}

//JoinMessage returns a parsed JOIN command
func JoinMessage(channel string) Message {
	return NewMessage("JOIN " + channel)
}

//Message represents a Message sent between the client and server
type Message interface {
	Message() string
	Prefix() string
	Nick() string
	User() string
	Host() string
	Server() string
	Command() string
	Params() []string
	Trailing() string
	String() string
	Timestamp() time.Time
}

/*A Message represents a message sent to or from the IRC server.

The various components of the message are parsed for convenience.

Pseudo-BNF from: https://tools.ietf.org/html/rfc1459#section-2.3.1


    <message>  ::= [':' <prefix> <SPACE> ] <command> <params> <crlf>
    <prefix>   ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
    <command>  ::= <letter> { <letter> } | <number> <number> <number>

    <params>   ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
    <middle>   ::= <Any *non-empty* sequence of octets not including SPACE
                   or NUL or CR or LF, the first of which may not be ':'>
    <trailing> ::= <Any, possibly *empty*, sequence of octets not including
                   NUL or CR or LF>
*/
type message struct {
	message string //The raw, unparsed message

	prefix string //includes the ':' character
	nick   string
	user   string
	host   string
	server string

	command  string
	params   []string //Includes Trailing as the final argument with the '
	trailing string   //excludes the ':'

	timestamp time.Time
	parsed    bool
}

//Message returns the entire message
func (m message) Message() string {
	return m.message
}

//Prefix eturns the prefix (including the preceeding colon), or an emtpy string if not present
func (m message) Prefix() string {
	return m.prefix
}

//Nick returns the nick from the prefix, or an emtpy string if not present
func (m message) Nick() string {
	return m.nick
}

//User returns the user from the prefix, or an emtpy string if not present
func (m message) User() string {
	return m.user
}

//Host returns the host from the prefix, or an emtpy string if not present
func (m message) Host() string {
	return m.host
}

//Server returns the server from the prefix, or an emtpy string if not present
func (m message) Server() string {
	return m.server
}

//Command returns the command from the message
func (m message) Command() string {
	return m.command
}

//Params returns the parameters from the message. The final param will be the
//trailing value, including the colon. An empty slice is returned if not present
func (m message) Params() []string {
	return m.params
}

//Trailing returns the last argument, excluding the colon, or an emtpy string if not present
func (m message) Trailing() string {
	return m.trailing
}

//Timestamp returns the timestamp the message was parsed
//TODO: Allow messages to be created with custome timestamps
func (m message) Timestamp() time.Time {
	return m.timestamp
}

//String returns the entire message (identical to calling Message())
func (m message) String() string {
	return m.Message()
}
