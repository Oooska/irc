package irc

import (
	"errors"
	"sort"
	"sync"
)

//ErrChannelDNE is returned when the specified channel does not exist
var ErrChannelDNE = errors.New("Channel Does Not Exist")

//Channels represents the list of channels the client is currently
//connected to, and the users present in those channels
type Channels interface {
	Users(channel string) (users []string, err error)
	ChannelNames() (channels []string)
	NumChannels() int
}

//userList represents a list of users in a channel
//The key is the username, the value is there mode.
type userList map[string]string

type channels struct {
	m     map[string]userList
	mLock *sync.RWMutex
}

func newChannels() channels {
	return channels{m: make(map[string]userList), mLock: new(sync.RWMutex)}
}

//Creates an empty channel.
func (c channels) Add(channel string) {
	c.mLock.Lock()
	c.m[channel] = make(userList)
	c.mLock.Unlock()
}

//Removes a channel
func (c channels) Remove(channel string) {
	c.mLock.Lock()
	delete(c.m, channel)
	c.mLock.Unlock()
}

//Adds the specified user to the specified channel.
//Returns ErrChannelDNE if channel does not exist
func (c channels) UserJoins(channel string, users ...string) error {
	c.mLock.Lock()
	defer c.mLock.Unlock()
	ul, ok := c.m[channel]
	if ok {
		for _, user := range users {
			ul[user] = ""
		}
		return nil
	}

	return ErrChannelDNE
}

//Removes the specified user from the specified channel.
//Returns ErrChannelDNE if room does not exist
func (c channels) UserParts(channel, user string) error {
	c.mLock.Lock()
	defer c.mLock.Unlock()
	ul, ok := c.m[channel]
	if ok {
		delete(ul, user)
		return nil
	}

	return ErrChannelDNE
}

//Removes the specified user from all channels
func (c channels) UserQuits(user string) {
	c.mLock.Lock()
	for _, ul := range c.m {
		delete(ul, user)
	}
	c.mLock.Unlock()
}

//Returns a sorted slice containing the users in a given channel.
//Returns an empty slice if no channel exists
//The bool value is true if the room exists, false otherwise
func (c channels) Users(channel string) ([]string, error) {
	c.mLock.RLock()
	defer c.mLock.RUnlock()
	ch, ok := c.m[channel]
	if ok {
		users := make([]string, len(ch))
		k := 0
		for user := range ch {
			users[k] = user
			k++
		}
		sort.Strings(users)
		return users, nil
	}
	return []string{}, ErrChannelDNE
}

//Returns the number of open channels
func (c channels) NumChannels() int {
	c.mLock.RLock()
	l := len(c.m)
	c.mLock.RUnlock()
	return l
}

//Returns a sorted list of channels
func (c channels) ChannelNames() []string {
	c.mLock.RLock()
	channels := make([]string, len(c.m))
	k := 0
	for key := range c.m {
		channels[k] = key
		k++
	}
	c.mLock.RUnlock()
	sort.Strings(channels)
	return channels
}
