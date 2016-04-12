package irc

import (
    "testing"
)


//TODO: Fix up these tests to use table-driven tests
//e.g. #4 https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742

//Test adding/removing channels using channelUserList
func TestChannelUserList_channels(t *testing.T){
    cul := newChannelUserList()
    
    if cul.NumChannels() != 0 {
        t.Errorf("A newly created channelUserList is suppose to be empty, but has length %d", cul.NumChannels())
    }
    
    channels := cul.Channels()
    if len(channels) != 0 {
        t.Errorf("A new created channelUserList does not return an empty slice of channel names. Found %d", len(channels))
    }
    
    cul.Add("#test")
    
    if cul.NumChannels() != 1 {
        t.Errorf("Adding a channel does not increase the number of channels. Expected: 1, found: %d", cul.NumChannels())
    }
    
    cul.Add("#test2")
    cul.Add("#test3")
    if cul.NumChannels() != 3 {
        t.Errorf("Adding a channel does not increase the number of channels. Expected: 3, found: %d", cul.NumChannels())
    }
    
    channels = cul.Channels()
    if channels[0] != "#test" || channels[1] != "#test2" || channels[2] != "#test3" {
        t.Errorf(`Channels() did not return the correct channel names in alphabetical order. `+
        ` Expected: ["#test", "#test1", "#test2"], Received: %+v`, channels)
    }
    
    cul.Remove("#test2")
    if cul.NumChannels() != 2 {
        t.Errorf("Remove() did not remove the requested item. Expected 2 items, found %d", cul.NumChannels())
    }
    
    channels = cul.Channels()
    if channels[0] != "#test" || channels[1] != "#test3" {
        t.Errorf(`Channels() did not return the correct channel names in alphabetical order. `+
        ` Expected: ["#test", "#test3"], Received: %+v`, channels)
    }
    
    cul.Remove("#test")
    if cul.NumChannels() != 1 {
        t.Errorf("Remove() did not remove the requested item. Expected 2 items, found %d", cul.NumChannels())
    }
    
    channels = cul.Channels()
    if channels[0] != "#test3" {
        t.Errorf(`Channels() did not return the correct channel names in alphabetical order. `+
        ` Expected: [#test3"], Received: %+v`, channels)
    }
    
    cul.Remove("#test3")
    if cul.NumChannels() != 0 {
        t.Errorf("Remove() did not remove the requested item. Expected 0 items, found %d", cul.NumChannels())
    }  
}

//Tests adding/removing users from specified channels
func TestChannelUserList_users(t *testing.T){
    cul := newChannelUserList()
    
    users, err := cul.Users("#no-such-channel-exists")
    if err == nil || len(users) != 0 {
        t.Errorf("Requesting a channel did not return with an empty slice, false. Received: %+v, %s",users, err.Error())
    }
    
    cul.Add("#test")
    cul.Add("#test2")
    
    users, err = cul.Users("#test")
    if err != nil {
        t.Errorf("Users() returned error %s when requesting a list of users for an existing channel.", err.Error())
    }
    if len(users) != 0 {
        t.Errorf("Users()) did not return an empty list when requesting a list of users in an empty channel. Received %v", users)
    }
    
    
    cul.UserJoins("#test", "user")
    
    users, err = cul.Users("#test")
    if err != nil {
        t.Errorf("Users() returned error %s when requesting a list of users for an existing channel", err.Error())
    }
    if len(users) != 1 {
        t.Errorf("Users() did not return the correct number of users. Expected 1, Received: %d", len(users))
    }
    if users[0] != "user" {
        t.Errorf("Users() returns an incorrect user. Expected \"user\", found \"%s\"",users[9])
    }
    
    cul.UserJoins("#test", "captain_planet")
    users, err = cul.Users("#test")
    if len(users) != 2 {
        t.Errorf("Users() did not return the correct number of users. Expected 2, Received: %d", len(users))
    }
    if users[0] != "captain_planet" || users[1] != "user" {
        t.Errorf(`Users() did not return the correct slice of users. Expected {"captain_planet", "user"}", found: %+v`, users)
    } 
    
    
    cul.UserJoins("#test", "captain_america")
    users, err = cul.Users("#test")
    if len(users) != 3 {
        t.Errorf("Users() did not return the correct number of users. Expected 3, Received: %d", len(users))
    }
    if users[0] != "captain_america" || users[1] != "captain_planet" || users[2] != "user" {
        t.Errorf(`Users() did not return the correct slice of users. Expected {"captain_america", captain_planet", "user"}", found: %+v`, users)
    } 
    
    cul.UserJoins("#test2", "user")
    cul.UserJoins("#test2", "captain_america")
    
    users, err = cul.Users("#test2")
    if len(users) != 2 {
        t.Errorf("Users() did not return the correct number of users. Expected 2, Received: %d", len(users))
    }
    if users[0] != "captain_america" || users[1] != "user" {
        t.Errorf(`Users() did not return the correct slice of users. Expected {"captain_america", "user"}", found: %+v`, users)
    } 
    
    cul.UserQuits("captain_america")
    users, err = cul.Users("#test")
    if len(users) != 2 {
        t.Errorf("Users() did not return the correct number of users. Expected 2, Received: %d", len(users))
    }
    if users[0] != "captain_planet" || users[1] != "user" {
        t.Errorf(`Users() returning a user that should have been removed after quitting. `+
          `Expected: {"captain_planet", "user"}, Received: %+v`, users)
    }
    
    users, err = cul.Users("#test2")
    if len(users) != 1 {
        t.Errorf("Users() did not return the correct number of users. Expected 1, Received: %d", len(users))
    }
    if users[0] != "user" {
        t.Errorf(`Users() returning a user that should have been removed after quitting. `+
          `Expected: {"captain_planet", "user"}, Received: %+v`, users)
    }
    
    cul.UserParts("#test", "user")
    users, err = cul.Users("#test")
    if len(users) != 1 {
        t.Errorf("Users() did not return the correct number of users. Expected 1, Received: %d", len(users))
    }
    if users[0] != "captain_planet" {
        t.Errorf(`Users() returning the incorrect list `+
          `Expected: {"captain_planet"}, Received: %+v`, users)
    }
    
}