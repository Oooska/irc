package irc

import (
    "errors"
)

type userList map[string]bool
type channelUserList map[string]userList

func newChannelUserList() channelUserList {
    cul := make(map[string]userList)
    return cul
}

//Creates an empty channel.
func (cul channelUserList) Join(channel string){
    cul[channel] = make(map[string]bool)
}

//Deletes a channel
func (cul channelUserList) Part(channel string){
    delete(cul, channel)
}

//Adds the specified user to the specified channel.
//Returns an error if channel does not exist
func (cul channelUserList) UserJoins(channel, user string) error {
    ul, ok := cul[channel]
    if ok {
        ul[user] = true
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

//Returns a slice containing the users in a given channel.
//Returns an empty slice if no channel exists
func (cul channelUserList) Users(channel string) []string{
    ch, ok := cul[channel]
    if ok {
        users := make([]string, len(ch))
        k := 0
        for user := range ch {
            users[k] = user
            k++
        }
        return users    
    }
    return []string{}
}