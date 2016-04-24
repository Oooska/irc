package irc

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"strings"
)

/****The SSL implementation is currently insecure. ***

Conn represents a connection to an IRC server. It provides
methods to read, write and close a connection.

MessageHandlers can be added to the Conn, and will be called anytime
a message is sent or recieved if it matches the specified criteria.

A NewConnectionWrapper method is provided to allow you to provide
your own implementation of net.Conn (e.g. for a websocket)*/
type Conn interface {
	Read() (Message, error)
	Write(Message) error
	Close()

	AddHandler(dir handlerDirection, mh MessageHandler, cmds ...string)
}

//MessageHandler are functions that will be called by a client upon
//recieving a message that matches the supplied criteria.
type MessageHandler func(Message)

//handlerDirection represents the direction a handler should be triggered on
type handlerDirection int

const (
	msgHandlerKey = "*" //key for general handler (triggered on all messages)
)

//NewConnection returns a new IRC Conn object
//TODO: The SSL implementation is currently insecure.
func NewConnection(serverAddress string, useSSL bool) (Conn, error) {
	var c net.Conn
	var err error

	if useSSL {
		conf := tls.Config{ //TODO: Fix insecure SSL connection
			InsecureSkipVerify: true,
		}
		c, err = tls.Dial("tcp", serverAddress, &conf)
	} else {
		c, err = net.Dial("tcp", serverAddress)
	}

	if err != nil {
		return nil, err
	}

	return NewConnectionWrapper(c), nil
}

//NewConnectionWrapper provides a new IRC Conn object using
//the specified input stream. Useful for websockets or other
//connectivity methods
func NewConnectionWrapper(c io.ReadWriteCloser) Conn {
	return &conn{conn: c,
		scanner:          bufio.NewScanner(c),
		incomingHandlers: make(map[string][]MessageHandler),
		outgoingHandlers: make(map[string][]MessageHandler),
	}
}

//A very simple implementation of an IRC client
type conn struct {
	conn    io.ReadWriteCloser
	scanner *bufio.Scanner

	incomingHandlers map[string][]MessageHandler
	outgoingHandlers map[string][]MessageHandler
}

//Read blocks until a new line is available from the server,
//It returns a new Message or returns an error
func (c *conn) Read() (msg Message, err error) {
	ok := c.scanner.Scan()
	if !ok {
		err = c.scanner.Err()
		if err == nil { //Scanner doesn't return EOF
			err = io.EOF
		}
		return
	}
	line := c.scanner.Text()
	msg = NewMessage(line)
	for _, h := range c.incomingHandlers[msgHandlerKey] {
		h(msg)
	}

	for _, h := range c.incomingHandlers[msg.Command] {
		h(msg)
	}

	return
}

//Writes the message to the server.
//Returns an error if one occurs
func (c *conn) Write(msg Message) error {
	_, err := c.conn.Write([]byte(msg.String() + "\r\n"))

	if err == nil {
		for _, h := range c.outgoingHandlers[msgHandlerKey] {
			h(msg)
		}

		for _, h := range c.outgoingHandlers[msg.Command] {
			h(msg)
		}
	}

	return err
}

//Closes the connection to the server. It does not send
//a quit command.
func (c *conn) Close() {
	if c != nil {
		c.conn.Close()
	}
}

//Adds a MessageHandler function to the client. The supplied handler
//will be called for all messages that are going in the specified direction
//(inbound, outbound or both). If commands are specified, the handler will be
//called only on those commands. If no commands are specified, the handler will
//be called for all messages, regardless of the command.
func (c *conn) AddHandler(dir handlerDirection, h MessageHandler, cmds ...string) {
	if len(cmds) < 1 {
		cmds = []string{msgHandlerKey}
	}

	if dir == Incoming || dir == Both {
		for _, cmd := range cmds {
			cmd = strings.ToUpper(cmd)
			handlers := c.incomingHandlers[cmd]
			c.incomingHandlers[cmd] = append(handlers, h)
		}
	}

	if dir == Outgoing || dir == Both {
		for _, cmd := range cmds {
			cmd = strings.ToUpper(cmd)
			handlers := c.outgoingHandlers[cmd]
			c.outgoingHandlers[cmd] = append(handlers, h)
		}
	}
}
