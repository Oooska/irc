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


func channelHandler(client *fullClient){
    cul := channelManager(client)
    client.channelUserList = cul
}

//Channel manager keeps track of which rooms you're in, and who else is in those channel
//Returns a channel user list, a message handler, and list of commands
//the handler should operate on in both directions.
//TODO: Keep track of modes / other pertinent data
//TODO: Listen for nick changes
func channelManager(client Client) channelUserList {
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
          case "KICK": 
          if msg.User.Nick != "" && len(msg.Params) > 0 {
              cul.UserParts(msg.Params[0], msg.User.Nick)
              //TODO: Determine if it was the client that got kicked
          }
          case "QUIT":
          if msg.User.Nick == ""{
              //Client is quitting, empty channel list
              for _, channel := range cul.Channels() {
                  cul.Remove(channel)
              }
          } else {
              //User is quitting
              cul.UserQuits(msg.User.Nick)
          }
      }
  }
  client.AddHandler(Both, handler, "JOIN", "PART", "KICK", "QUIT")
  return cul
}
