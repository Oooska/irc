package irc

import (
	"log"
)

/* Client Handlers are functions that add useful functionality to
   a irc.Client. They do this by attaching MessageHandlers to the
   client.
*/

//ClientHandler is a function that attaches MessageHandlers to a client
type ClientHandler func(Client)

//LogHandler logs all messages to the default logger
func LogHandler(client Client) {
	handler := func(msg Message) {
		log.Printf(msg.Message)
	}
	client.AddHandler(Both, handler)
}

//PingHandler registers a handler to respond to pings
func pingHandler(client Client) {
	handler := func(msg Message) {
		resp := "PONG"
		if len(msg.Params) > 0 {
			resp += " " + msg.Params[0]
		}

		client.Send(NewMessage(resp))
	}

	client.AddHandler(Incoming, handler, "PING")
}

//Registers the channels handler to a fullclient, and sets the
//channels object.
func channelHandler(client *clientImpl) {
	ch := RegisterChannelsHandler(client)
	client.Channels = ch
}

//RegisterChannelsHandler keeps track of which rooms you're in, and who else is in those channels
//Returns a Channels object.
//TODO: Keep track of modes / other pertinent data
//TODO: Listen for nick changes
func RegisterChannelsHandler(c Conn) Channels {
	cul := newChannelUserList()
	handler := func(msg Message) {
		switch msg.Command {
		case "JOIN":
			if len(msg.Params) > 0 {
				if msg.Nick == "" {
					//JOIN #room
					cul.Add(msg.Params[0])
				} else {
					//nick JOIN #room
					cul.UserJoins(msg.Params[0], msg.Nick)
				}
			} //else malformed request - ignoring
		case "PART":
			if len(msg.Params) > 0 {
				if msg.Nick == "" {
					//PART #room
					cul.Remove(msg.Params[0])
				} else {
					//nick PART #channel :reason
					cul.UserParts(msg.Params[0], msg.Nick)
				}
			} //else malformed request - ignoring
		case "KICK":
			if msg.Nick != "" && len(msg.Params) > 0 {
				cul.UserParts(msg.Params[0], msg.Nick)
				//TODO: Determine if it was the client that got kicked
			}
		case "QUIT":
			if msg.Nick == "" {
				//Client is quitting, empty channel list
				for _, channel := range cul.ChannelNames() {
					cul.Remove(channel)
				}
			} else {
				//User is quitting
				cul.UserQuits(msg.Nick)
			}
		}
	}
	c.AddHandler(Both, handler, "JOIN", "PART", "KICK", "QUIT")
	return Channels(cul)
}
