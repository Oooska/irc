package irc


//PingHandler registers a handler to respond to pings
func PingHandler(client BasicClient){
    handler := func(msg Message){
      if len(msg.Params) > 0{
        client.Send(NewMessage("PONG "+msg.Params[0]))
      }
    }
    
    client.AddIncomingHandlerOnCmd(handler, "PING")
}

