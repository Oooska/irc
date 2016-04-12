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
//TODO: Keep track of modes / other data
//you're in, and who else is in those channels
func channelManager(client *fullClient){
  //Will use channelUserList in client-structs to implement
  cul := newChannelUserList()
  handler := func(msg Message){
      switch(msg.Command){
          case "JOIN":
          if len(msg.Params) > 0 {
                if msg.User.Nick == "" {
                    //JOIN #room
                    cul.Add(msg.Params[0])
                } else {
                    //nick JOIN #room
                    cul.UserJoins(msg.Params[0], msg.User.Nick)
                }
          } //else malformed request - ignoring
          case "PART":
          if len(msg.Params) > 0 {
              if msg.User.Nick == "" {
                  //PART #room
                  cul.Remove(msg.Params[0])
              } else {
                  //nick PART #channel :reason
                  cul.UserParts(msg.Params[0], msg.User.Nick)
              }
          } //else malformed request - ignoring
          case "QUIT":
          if msg.User.Nick == ""{
              //Client is quitting, empty channel list
              cul = newChannelUserList()
              client.cul = cul
          } else {
              cul.UserQuits(msg.User.Nick)
          }
      }
  }
  client.cul = cul 
  client.AddHandler(Incoming, handler, "JOIN", "PART", "QUIT")
}
