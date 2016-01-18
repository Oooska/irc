package irc

import (
	"bufio"
	"crypto/tls"
	"net"
)

/* A ridiculously basic IRC client. A simple wrapped around net.Conn, with some
simple helper methods to generate IRC commands.

***The SSL implementation is currently insecure. *** */

type IRCConn interface {
	Read() (Message, error)
	Write(Message) error
	Close()
}

type simpleClient struct {
	conn     net.Conn
	buffconn *bufio.Reader
}

func (client simpleClient) Read() (msg Message, err error) {
	line, err := client.buffconn.ReadString('\n')
	if err == nil {
		msg = Message(line)
	}
	return
}

func (client simpleClient) Write(msg Message) error {
	_, err := client.conn.Write([]byte(msg))
	return err
}

func (client *simpleClient) Close() {
	client.conn.Close()
}

func IRCConnectionWrapper(conn net.Conn) IRCConn {
	return &simpleClient{conn: conn, buffconn: bufio.NewReader(conn)}
}

func NewIRCConnection(serverAddress string, useSSL bool) (IRCConn, error) {
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

	return IRCConnectionWrapper(conn), nil
}
