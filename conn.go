package irc

import (
	"bufio"
	"crypto/tls"
	"net"
)

/* Conn represents a connection to an IRC server. It provides
methods to read, write and close a connection. 

A NewConnectionWrapper method is provided to allow you to provide
your own implementation of net.Conn (e.g. for a websocket)

***The SSL implementation is currently insecure. *** */
type Conn interface {
	Read() (Message, error)
	Write(Message) error
	Close()
}

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

	return &conn{conn: c, buffconn: bufio.NewReader(c)}, nil
} 

//A very simple implementation of an IRC client
type conn struct {
	conn     net.Conn
	buffconn *bufio.Reader
}

//Read blocks until a new line is available from the server,
//It returns a new Message or returns an error
func (c *conn) Read() (msg Message, err error) {
	line, err := c.buffconn.ReadString('\n')
	if err == nil {
		msg = NewMessage(line)
	}
    return
}

//Writes the message to the server.
//Returns an error if one occurs
func (c *conn) Write(msg Message) error {
	_, err := c.conn.Write([]byte(msg.String()+"\r\n"))
	return err
}

//Closes the connection to the server. It does not send  
//a quit command.
func (c *conn) Close() {
    if c != nil{
	    c.conn.Close()
    }
}

//NewConnectionWrapper provides a new IRC Conn object using the supplied net.Conn object
func NewConnectionWrapper(c net.Conn) Conn {
    return &conn{conn: c, buffconn: bufio.NewReader(c)}
}