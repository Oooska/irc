package irc

import (
	"log"
	"strings"
	"sync"
)

/* Client Handlers are functions that add useful functionality to
   a irc.Client. They do this by attaching MessageHandlers to the
   client.
*/

//ClientHandler is a function that attaches MessageHandlers to a client
type ClientHandler func(Client)

const (
	rplName       = "353"
	rplEndofNames = "366"
)

//LogHandler logs all messages to the default logger
func LogHandler(client Client) {
	handler := func(msg Message) {
		log.Printf(msg.Message())
	}
	client.AddHandler(Both, handler)
}

func conversationHandler(client *clientImpl) {
	convo := RegisterConversationsHandler(client)
	client.Conversations = convo
}

//RegisterConversationsHandler registers the conversation handler
//with the connection and returns a Conversations object to access
//captured data.
func RegisterConversationsHandler(c Conn) Conversations {
	convos := newConversations(1024)
	handler := func(msg Message) {
		convos.Add(msg.Params()[0], msg.Message())
	}
	c.AddHandler(Both, handler, "PRIVMSG")
	return convos
}

//PingHandler registers a handler to respond to pings
func pingHandler(client Client) {
	handler := func(msg Message) {
		resp := "PONG"
		if len(msg.Params()) > 0 {
			resp += " " + msg.Params()[0]
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
	cul := newChannels()
	namesUpdating := make(map[string]bool) //Keeps track of rplName/rplEndofNames
	namesUpdatingLock := new(sync.Mutex)
	handler := func(msg Message) {
		switch msg.Command() {
		case "JOIN":
			if len(msg.Params()) > 0 {
				if msg.Nick() == "" {
					//JOIN #room
					cul.Add(msg.Params()[0])

					//ANeed to note that a /names update is coming for this channel
					namesUpdatingLock.Lock()
					namesUpdating[msg.Params()[0]] = true
					namesUpdatingLock.Unlock()
				} else {
					//nick JOIN #room
					cul.UserJoins(msg.Params()[0], msg.Nick())
				}
			} //else malformed request - ignoring
		case "PART":
			if len(msg.Params()) > 0 {
				if msg.Nick() == "" {
					//PART #room
					cul.Remove(msg.Params()[0])
				} else {
					//nick PART #channel :reason
					cul.UserParts(msg.Params()[0], msg.Nick())
				}
			} //else malformed request - ignoring
		case "KICK":
			if msg.Nick() != "" && len(msg.Params()) > 0 {
				cul.UserParts(msg.Params()[0], msg.Nick())
				//TODO: Determine if it was the client that got kicked
			}
		case "QUIT":
			if msg.Nick() == "" {
				//Client is quitting, empty channel list
				for _, channel := range cul.ChannelNames() {
					cul.Remove(channel)
				}
			} else {
				//User is quitting
				cul.UserQuits(msg.Nick())
			}
		case "NAMES":
			if len(msg.Params()) >= 1 {
				//ANeed to node that a /names update is coming for this channel
				namesUpdatingLock.Lock()
				namesUpdating[msg.Params()[0]] = true
				namesUpdatingLock.Unlock()
			} //else /names will report list of ALL public channels.
			//TODO: Show user list of all public channels.
		case rplName: //List of nicks in the specified channel
			//:tepper.freenode.net 353 nick @ #gotest :goirctest @Oooska
			namesUpdatingLock.Lock()
			defer namesUpdatingLock.Unlock()
			if len(msg.Params()) >= 3 {
				ch := msg.Params()[2]
				updating, _ := namesUpdating[ch]
				if updating {
					//Only update names if we're requesting the info
					//from a /names #channel or /join command
					names := strings.Split(msg.Trailing(), " ")
					cul.UserJoins(ch, names...)
				}

			}
		case rplEndofNames:
			//:tepper.freenode.net 366 goirctest #gotest :End of /NAMES list.
			namesUpdatingLock.Lock()
			defer namesUpdatingLock.Unlock()
			if len(msg.Params()) >= 2 {
				ch := msg.Params()[1]
				delete(namesUpdating, ch)
			}
		}
	}
	c.AddHandler(Both, handler, "JOIN", "PART", "KICK", "QUIT", "NAMES")
	c.AddHandler(Incoming, handler, rplName, rplEndofNames)
	return Channels(cul)
}
