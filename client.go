package irc

const (
    msgHandlerKey = "*" //key for LiteClient.handlers for general handler
)

//MessageHandler are functions that will be called by a client upon
//recieving a message that matches the supplied criteria. 
type MessageHandler func(Message)

//LiteClient interface provides a barebones client to send/recieved
//messages, and allows you to add MessageHandlers
type Client interface {
    Next() (Message, error)
    Send(...Message) (int, error)
    Close() 
    
    AddHandler(direction handlerDirection, handler MessageHandler, commands ...string) 
}

//Represents the direction a handler should be triggered on
type handlerDirection int
const (
    //handlerDirections
    Incoming = handlerDirection(iota)
    Outgoing
    Both
)


//NewLiteClient returns a basic IRC client interface with
//the ability to Read and Write messages to a server,
//as well as add 
func NewClient(serverAddress string, useSSL bool) (Client, error) {
    conn, err := NewConnection(serverAddress, useSSL)
    if err != nil {
        return nil, err
    }

    return &liteClient{ conn: conn }, nil
}


type fullClient struct {
    liteClient
    Channels
}

//LiteClient implements the LiteClient interface
type liteClient struct {
    conn    Conn
    incomingHandlers map[string][]MessageHandler
    outgoingHandlers map[string][]MessageHandler
}

//Next reads the next message, calling all handlers
//Returns the next Message, and an error if one occured.
func (c liteClient) Next() (Message, error){
    msg, err := c.conn.Read()
    
    if err == nil {
        for _, h := range c.incomingHandlers[msgHandlerKey]{
            h(msg)
        }
        
        for _, h := range c.incomingHandlers[msg.Command]{
            h(msg)
        }
    }
    
    return msg, err
}

//Send sends all of the supplied messages to the server.
//Returns the number of successfully sent messages. It 
//stops and returns the first error message recieved
func (c liteClient) Send(msgs ...Message) (int, error) {
    var k int
    var msg Message
    for k, msg = range msgs {
        err := c.conn.Write(msg)
        if err != nil {
            return k-1, err
        } 
        
        for _, h := range c.outgoingHandlers[msgHandlerKey]{
            h(msg)
        }
        
        for _, h := range c.outgoingHandlers[msg.Command]{
            h(msg)
        }
    
    }
    return k, nil
}

//Closes the connection to the iRC server. It does not
//send a QUIT message.
func (c *liteClient) Close() {
    if c != nil {
        c.conn.Close()
    }
}

//Adds a MessageHandler function to the client. The supplied handler
//will be called for all messages that are going in the specified direction
//(inbound, outbound or both). If commands are specified, the handler will be
//called only on those commands. If no commands are specified, the handler will 
//be called for all messages, regardless of the command.
func (c liteClient) AddHandler(dir handlerDirection, h MessageHandler, cmds ...string){

    if len(cmds) < 1 {
        cmds = []string{msgHandlerKey}
    }
    
    if dir == Incoming || dir == Both{
        for _, cmd := range cmds {
            handlers := c.incomingHandlers[cmd]
            c.incomingHandlers[cmd] = append(handlers, h)
        }
    } 
    
    if dir == Outgoing || dir == Both {
        for _, cmd := range cmds {
            handlers := c.outgoingHandlers[cmd]
            c.outgoingHandlers[cmd] = append(handlers, h)
        }        
    }
}