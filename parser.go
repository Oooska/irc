package irc

import (
    "strings"
)

//ParseString string takes a raw irc command and parses it
//into a ParsedMessage
//:PREFIX COMMAND ARG1 ARG2 :Last arg may have spaces if preceeded by colon
//PREFIX is nick!user@host or servername, and is optional
func ParseString(message string) (pm Message) {
    tokens := strings.Split(strings.TrimSpace(message), " ")
    k := 0
    
    pm.Message = message
    
    //Check for prefix
    if k < len(tokens) && parsePrefix(tokens[k], &pm) {
        k++
    }
    
    //Parse command, ignoring empty tokens
    for ;k < len(tokens); k++{
        if tokens[k] != "" {
            pm.Command = tokens[k]
            k++
            break
        }
    }
    
    //Check params, ignoring empty tokens.
    //The last argument will include
    //the ':' if present
    for ; k < len(tokens); k++ {
        if tokens[k] == "" {
            continue
        } else if tokens[k][0] == ':' {
            //Grab the rest of the string
            s := strings.SplitAfterN(message[1:], ":", 2)
            if len(s) > 1 {
                pm.Params = append(pm.Params, ":"+s[1])
            }
            return
        } 
        
        pm.Params = append(pm.Params, tokens[k])
    }
    return
}


//parses a prefix, and updates the parsedMEssage fields. Returns true if the string is a prefix
func parsePrefix(prefix string, pm *Message) bool {
    if len(prefix) < 1 || prefix[0] != ':' {
        return false
    }
    pm.Prefix = prefix
    
    //Check for the '!' in the prefix
    i := strings.Index(prefix, "!")
    if (i < 0){//If not present, this is the server name
        pm.Server = prefix[1:]
        return true
    }
    
    pm.User = User{}
    pm.User.Nick = prefix[1:i]
    
    iat := strings.Index(prefix, "@")
    if(iat < 0){
        //No host provided, just nick!user
        pm.User.User = prefix[i+1:]
        return true
    }
    pm.User.User = prefix[i+1:iat]
    pm.User.Host = prefix[iat+1:]
    return true
}


/*  TODO: Implement a parser for the lexer. */
