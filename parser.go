package irc

import (
    "strings"
)

//ParseString string takes a raw irc command and parses it
//into a ParsedMessage
func ParseString(message string) (pm ParsedMessage) {
    
    tokens := strings.Split(strings.TrimSpace(message), " ")

    wordCt := 0
    
    pm.Message = message
    
    //Check for prefix
    if wordCt < len(tokens) && parsePrefix(tokens[wordCt], &pm) {
        wordCt++
    }
    
    //Parse command, ignoring empty tokens
    for ;wordCt < len(tokens); wordCt++{
        if tokens[wordCt] != "" {
            pm.Command = tokens[wordCt]
            wordCt++
            break
        }
    }
    
    //Check params, ignoring empty tokens.
    //The last argument will include
    //the ':' if present
    for ; wordCt < len(tokens); wordCt++ {
        if tokens[wordCt] == "" {
            continue
        } else if tokens[wordCt][0] == ':' {
            //Grab the rest of the string
            s := strings.SplitAfterN(message[1:], ":", 2)
            if len(s) > 1 {
                pm.Params = append(pm.Params, ":"+s[1])
            }
            return
        } 
        
        pm.Params = append(pm.Params, tokens[wordCt])
    }
    
    
    return
}


//parses a prefix, and updates the parsedMEssage fields. Returns true if the string is a prefix
func parsePrefix(prefix string, pm *ParsedMessage) bool {
    if prefix[0] != ':' {
        return false
    }
    pm.Prefix = prefix
    
    //Check for the '!' in the host
    i := strings.Index(prefix, "!")
    if (i < 0){//If not present, this is the server name
        pm.Server = prefix[1:]
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


/*  TODO: Implement a parser for the lexer:
type parser struct {
	l   *lexer
	buf bytes.Buffer
}

//Returns a new parser
func newParser(l *lexer) parser {
	p := parser{l: l}
	return p
}

//Grabs the next token & literal, writes the literal to the buffer, and returns both
func (p *parser) nextToxen() (Token, string) {
	token, literal := p.l.NextItem()
	p.buf.Write(literal)
	return token, string(literal)
}

/*func (p *parser) ParseMessage() (parsedMessage, error) {
    
}*/
