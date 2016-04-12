package irc

import (
    "errors"
    "sort"
)


 var   ErrChannelDNE = errors.New("Room Does Not Exist")


//Channels represents the list of channels the client is currently
//connected to, and the use
type Channels interface {
    Add(channel string)
    Remove(channel string)
    UserJoins(channel, nick string) error
    UserParts(channel, nick string) error
    UserQuits(nick string)
    Users(channel string) (users []string, err error)
    Channels() (channels []string)
    NumChannels() int
}

//userList represents a list of users in a channel
//The key is the username, the value is there mode.
type userList map[string]string
//channelUserList is a map of room names to userLists
type channelUserList map[string]userList

func newChannelUserList() channelUserList {
    cul := make(channelUserList)
    return cul
}

//Creates an empty channel.
func (cul channelUserList) Add(channel string){
    cul[channel] = make(userList)
}

//Removes a channel
func (cul channelUserList) Remove(channel string){
    delete(cul, channel)
}

//Adds the specified user to the specified channel.
//Returns ErrChannelDNE if channel does not exist
func (cul channelUserList) UserJoins(channel, user string) error {
    ul, ok := cul[channel]
    if ok {
        ul[user] = ""
        return nil
    }
    
    return ErrChannelDNE
}

//Removes the specified user from the specified channel.
//Returns ErrChannelDNE if room does not exist
func (cul channelUserList) UserParts(channel, user string) error {
    ul, ok := cul[channel]
    if ok {
        delete(ul, user)
        return nil
    }
    return ErrChannelDNE 
}

//Removes the specified user from all channels
func (cul channelUserList) UserQuits(user string){
    for _, ul := range cul {
        delete(ul, user)
    }
}

//Returns a sorted slice containing the users in a given channel.
//Returns an empty slice if no channel exists
//The bool value is true if the room exists, false otherwise
func (cul channelUserList) Users(channel string) ([]string, error) {
    ch, ok := cul[channel]
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
func (cul channelUserList) NumChannels() int {
    return len(cul)
}

//Returns a sorted list of channels
func (cul channelUserList) Channels() []string {
    channels := make([]string, len(cul))
    k := 0
    for key := range cul {
        channels[k] = key
        k++
    }
    sort.Strings(channels)
    return channels
}