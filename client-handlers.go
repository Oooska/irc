package irc


//PingHandler registers a handler to respond to pings
func PingHandler(client Client){
    handler := func(msg Message){
      resp := "PONG "
      if len(msg.Params) > 0{
        resp += msg.Params[0]
      } 
      
      client.Send(NewMessage(resp))
    }
    
    client.AddHandler(Incoming, handler, "PING")
}



//Channel manager keeps track of which rooms
//you're in, and who else is in those channels
func ChannelManager(client Client){
  //Will use channelUserList in client-structs to implement
}
