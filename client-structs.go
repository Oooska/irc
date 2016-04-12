package irc

import (
    "errors"
    "sort"
)


//userList represents a list of users in a channel
//The key is the username, the value is there mode.
//TODO: Implement modes
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
//Returns an error if channel does not exist
func (cul channelUserList) UserJoins(channel, user string) error {
    ul, ok := cul[channel]
    if ok {
        ul[user] = ""
        return nil
    }
    
    return errors.New("No channel exists.") 
}

//Removes the specified user from the specified channel
func (cul channelUserList) UserParts(channel, user string) error {
    ul, ok := cul[channel]
    if ok {
        delete(ul, user)
        return nil
    }
    return errors.New("No channel exists.")  
}

//Removes the specified user from all channels
func (cul channelUserList) UserQuits(user string){
    for _, ul := range cul {
        delete(ul, user)
    }
}

//Returns a sorted slice containing the users in a given channel.
//The bool value is true if the room exists
//Returns an empty slice if no channel exists
func (cul channelUserList) Users(channel string) ([]string, bool) {
    ch, ok := cul[channel]
    if ok {
        users := make([]string, len(ch))
        k := 0
        for user := range ch {
            users[k] = user
            k++
        }
        sort.Strings(users)
        return users, ok    
    }
    return []string{}, ok
}

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