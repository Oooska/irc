package irc

//Client maintains common state elements used by an IRC client.
//Unlike the Conn, it keeps track of which channels you are in,
//who else is in those channels, modes for users, etc.
type Client interface {
	Send(...Message) (int, error)
	Conn
	Channels
}

const (
	Incoming handlerDirection = iota
	Outgoing
	Both
)

//NewClient returns a basic IRC client interface with
//the ability to Read and Write messages to a server,
//as well as add
func NewClient(serverAddress string, useSSL bool, handlers ...ClientHandler) (Client, error) {
	conn, err := NewConnection(serverAddress, useSSL)
	if err != nil {
		return nil, err
	}

	c := clientImpl{
		Conn: conn,
	}
	channelHandler(&c)
	pingHandler(&c)

	for _, h := range handlers {
		h(&c)
	}

	return &c, nil
}

//LiteClient implements the LiteClient interface
type clientImpl struct {
	Conn
	Channels
}

//Send sends all of the supplied messages to the server.
//Returns the number of successfully sent messages. It
//stops sending and returns the first error message recieved
func (c clientImpl) Send(msgs ...Message) (int, error) {
	var k int
	var msg Message
	for k, msg = range msgs {
		err := c.Write(msg)
		if err != nil {
			return k - 1, err
		}
	}
	return k, nil
}

//Closes the connection to the iRC server. It does not
//send a QUIT message.
func (c *clientImpl) Close() {
	if c != nil {
		c.Close()
	}
}
