package irc

import "sync"

func newConversations(length int) conversations {
	return conversations{messages: make(map[string][]string), mLock: new(sync.RWMutex), length: length}
}

//Conversations keeps track of the last 'length' privmessages to a channel
type Conversations interface {
	Messages(string) []string
}

type conversations struct {
	messages map[string][]string
	mLock    *sync.RWMutex
	length   int
}

//Adds the specified message to the logs
func (c conversations) Add(ch, message string) {
	c.mLock.Lock()
	messages := c.messages[ch]
	messages = append(messages, message)
	if len(messages) > c.length {
		messages = messages[1:]
	}
	c.messages[ch] = messages
	c.mLock.Unlock()
}

//Returns the current messages logged for the specified channel
func (c conversations) Messages(ch string) []string {
	c.mLock.RLock()
	messages := c.messages[ch]
	c.mLock.RUnlock()
	return messages
}
