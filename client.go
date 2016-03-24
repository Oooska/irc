package irc

const (
    msgHandlerKey = "*" //key for basicClient.handlers for general handler
)


//MessageHandler are functions that will be called by a client upon
//recieving a message that matches the supplied criteria. 
type MessageHandler func(Message)


//BasicClient interface provides a bareboned client to send/recieved
//messages, and allows you to add MessageHandlers
type BasicClient interface {
    Next() (Message, error)
    Send(...Message) (int, error)
    Close() 
    
    AddHandler(MessageHandler)                          //Called on all incoming/outgoing messages
    
    AddIncomingHandler(MessageHandler)                  //Called for all incoming messages
    AddIncomingHandlerOnCmd(MessageHandler, ...string)  //Called on specified commands
    
    AddOutgoingHandler(MessageHandler)                  //Called for all incoming messages
    AddOutgoingHandlerOnCmd(MessageHandler, ...string)  //Called on specified commands
}


//NewBasicClient returns a basic IRC client interface with
//the ability to Read and Write messages to a server,
//as well as add 
func NewBasicClient(serverAddress string, useSSL bool) (BasicClient, error) {
    conn, err := NewConnection(serverAddress, useSSL)
    if err != nil {
        return nil, err
    }

    return &basicClient{ conn: conn }, nil
}


//basicClient implements the BasicClient interface
type basicClient struct {
    conn    Conn
    incomingHandlers map[string][]MessageHandler
    outgoingHandlers map[string][]MessageHandler
}

//Next reads the next message, calling all handlers
//Returns the next Message, and an error if one occured.
func (c basicClient) Next() (Message, error){
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
func (c basicClient) Send(msgs ...Message) (int, error) {
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
func (c *basicClient) Close() {
    if c != nil {
        c.conn.Close()
    }
}

//Adds a MessageHandler function to the client. The supplied
//function will be called for every new message sent or
//recieved by the client
func (c basicClient) AddHandler(h MessageHandler){
    c.AddIncomingHandler(h)
    c.AddOutgoingHandler(h)
}


func (c basicClient) AddIncomingHandler(h MessageHandler){
    c.AddIncomingHandlerOnCmd(h, msgHandlerKey)
}

func (c basicClient) AddIncomingHandlerOnCmd(h MessageHandler, cmds ...string){
    for _, cmd := range cmds {
        handlers := c.incomingHandlers[cmd]
        c.incomingHandlers[cmd] = append(handlers, h)  
    }    
}

func (c basicClient) AddOutgoingHandler(h MessageHandler){
    c.AddOutgoingHandlerOnCmd(h, msgHandlerKey)
}

func (c basicClient) AddOutgoingHandlerOnCmd(h MessageHandler, cmds ...string){
    for _, cmd := range cmds {
        handlers := c.outgoingHandlers[cmd]
        c.outgoingHandlers[cmd] = append(handlers, h)  
    }    
}