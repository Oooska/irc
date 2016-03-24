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

//Conn represents a basic connection to an IRC client. 
type Conn interface {
	Read() (Message, error)
	Write(Message) error
	Close()
}

//NewConnection returns a new IRC Conn object
//TODO: The SSL implementation is currently insecure.
func NewConnection(serverAddress string, useSSL bool) (Conn, error) {
	var conn net.Conn
	var err error

	if useSSL {
		conf := tls.Config{ //TODO: Fix insecure SSL connection
			InsecureSkipVerify: true,
		}

		conn, err = tls.Dial("tcp", serverAddress, &conf)
	} else {
		conn, err = net.Dial("tcp", serverAddress)
	}

	if err != nil {
		return nil, err
	}

	return &simpleConn{conn: conn, buffconn: bufio.NewReader(conn)}, nil
}

//A very simple implementation of an IRC client
type simpleConn struct {
	conn     net.Conn
	buffconn *bufio.Reader
}

//Read blocks until a new line is available from the server,
//It returns a new Message or returns an error
func (client simpleConn) Read() (msg Message, err error) {
	line, err := client.buffconn.ReadString('\n')
	if err == nil {
		msg = NewMessage(line)
	}
    return
}

//Writes the message to the server.
//Returns an error if one occurs
func (client simpleConn) Write(msg Message) error {
	_, err := client.conn.Write([]byte(msg.String()+"\r\n"))
	return err
}

//Closes the connection to the server. It does not send  
//a quit command.
func (client *simpleConn) Close() {
    if client != nil{
	    client.conn.Close()
    }
}

//NewConnectionWrapper provides a new IRC Conn object using the supplied net.Conn object
func NewConnectionWrapper(conn net.Conn) Conn {
    return &simpleConn{conn: conn, buffconn: bufio.NewReader(conn)}
}