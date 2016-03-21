package irc

import (
    "testing"
)


var parserInput = []ParsedMessage{
        ParsedMessage{ Message: "PING :tepper.freenode.net", Command: "PING",  Params: []string{":tepper.freenode.net"}},
        ParsedMessage{ Message: ":tepper.freenode.net 332 goirctest #mainehackerclub :MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM\\", 
            Prefix: ":tepper.freenode.net", Command: "332", Server: "tepper.freenode.net",
            Params: []string{"goirctest", "#mainehackerclub", ":MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM\\"}},
        ParsedMessage{ Message: ":tepper.freenode.net 353 goirctest = #mainehackerclub :goirctest +ubuntuguru " +
	        /*   */"+FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman " +
	        /*   */"@Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 " +
	        /*   */"+zgrep lambdabot +jrvc", Prefix: ":tepper.freenode.net", Command: "353", Server: "tepper.freenode.net",
            Params: []string{"goirctest", "=", "#mainehackerclub", ":goirctest +ubuntuguru " +
	        /*   */"+FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman " +
	        /*   */"@Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 " +
            /*   */"+zgrep lambdabot +jrvc",
        }},
        ParsedMessage{Message: ":wallyworld!~quassel@1-2-3-4.static.tpgi.com.au  QUIT :Remote host closed the connection",
            Prefix: ":wallyworld!~quassel@1-2-3-4.static.tpgi.com.au", Command: "QUIT", Nick: "wallyworld", User: "~quassel", 
            Host: "1-2-3-4.static.tpgi.com.au", Params: []string{":Remote host closed the connection"}},
        ParsedMessage{Message: ":KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald PRIVMSG #go-nuts :https://golang.org/pkg/time/#Time.String",
            Prefix: ":KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald", Nick: "KirkMcDonald", User: "~Kirk", 
            Host: "python/site-packages/KirkMcDonald",  Command: "PRIVMSG", Params: []string{"#go-nuts", ":https://golang.org/pkg/time/#Time.String"}},
        ParsedMessage{Message: ":somenick!~@5-6-7-8.static.bgth.bz  QUIT", Prefix: ":somenick!~@5-6-7-8.static.bgth.bz", 
            Nick: "somenick", User: "~", Host: "5-6-7-8.static.bgth.bz", Command: "QUIT"},
}

func TestParseString(t *testing.T){
    
    
    for j, expected := range parserInput {
        actual := ParseString(expected.Message)
        
        if actual.Message != expected.Message {
            t.Errorf("input[%d]: Message field not parsed correctly. Expected: %s. Received: %s", j, expected.Message, actual.Message)
        }
        
        if actual.Prefix != expected.Prefix {
            t.Errorf("input[%d]: Prefix field not parsed correctly. Expected: %s. Received: %s", j, expected.Prefix, actual.Prefix)
        }
        
        if actual.Command != expected.Command {
            t.Errorf("input[%d]: Command field not parsed correctly. Expected: %s. Received: %s", j, expected.Command, actual.Command)
        }       
        
        if len(actual.Params) != len(expected.Params) {
            t.Errorf("input[%d]: Unequal number of params. Expected: %+v. Received: %+v", j, expected.Params, actual.Params)
        } else {
            for k := range actual.Params {
                if actual.Params[k] != expected.Params[k] {
                    t.Errorf("input[%d]: Param[%d] is not equal. Expected: %s. Received: %s", j, k, expected.Params[k], actual.Params[k])
                }
            }
        }
        
        if actual.Nick != expected.Nick {
            t.Errorf("input[%d]: Nick field not parsed correctly. Expected: %s. Received: %s", j, expected.Nick, actual.Nick)
        }       
        
        if actual.User != expected.User {
            t.Errorf("input[%d]: User field not parsed correctly. Expected: %s. Received: %s", j, expected.User, actual.User)
        }
        
        if actual.Host != expected.Host {
            t.Errorf("input[%d]: Host field not parsed correctly. Expected: %s. Received: %s", j, expected.Host, actual.Host)
        }              
        
        if actual.Server != expected.Server {
            t.Errorf("input[%d]: Server field not parsed correctly. Expected: %s. Received: %s", j, expected.Server, actual.Server)
        }       
    }
    
}