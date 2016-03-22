package irc

import (
    "strings"
)

//ParseString string takes a raw irc command and parses it
//into a ParsedMessage
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
    if prefix[0] != ':' {
        return false
    }
    pm.Prefix = prefix
    
    //Check for the '!' in the host
    i := strings.Index(prefix, "!")
    if (i < 0){//If not present, this is the server name
        pm.Host = prefix[1:]
        return true
    }
    
    pm.Nick = prefix[1:i]
    
    iat := strings.Index(prefix, "@")
    if(iat < 0){
        //No host provided, just nick!user
        pm.User = prefix[i+1:]
        return true
    }
    pm.User = prefix[i+1:iat]
    pm.Host = prefix[iat+1:]
    return true
}


/*  TODO: Implement a parser for the lexer. */
