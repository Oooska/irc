package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/oooska/irc"
)

var (
	address  = flag.String("address", "irc.freenode.net:6667", "IRC server address")
	ssl      = flag.Bool("ssl", false, "Use SSL")
	nick     = flag.String("nick", "go_test_client", "User nick")
	username = flag.String("username", "go_name", "User name")
)

//Code adapted from https://github.com/husio/irc/blob/master/examples/echobot.go
//A barebones IRC 'client' in the loosest sense of the word.
//Takes input from console. If command starts with a '/', everything after is sent as a raw IRC command.
//Otherwise the first argument is considered the channel/username, and the rest of the line is the message to send
// as a privmsg.
func main() {
	flag.Parse()

	fmt.Println("Simple Text-Based IRC Test Client\n")


    fmt.Printf("Connecting to %s . . . \n", *address)
    
    //LogHandler will handle printing out to stdio unless we change the default logger
    client, err := irc.NewClient(*address, *ssl, irc.LogHandler)
    
	if err != nil {
		panic(err)
	}
    fmt.Printf("Connected.\n\n")
    
    client.Send(irc.UserMessage(*username, "host", "domain", "realname"))
    client.Send(irc.NickMessage(*nick))
    client.Send(irc.JoinMessage("#go_test"))


    go func(){
        reader := bufio.NewReader(os.Stdin)
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                log.Fatalf("Cannot read from stdin: %s", err)
            }

            line = strings.TrimSpace(line)
            if len(line) == 0 {
                continue
            }

            msg, err := parseLine(line)
            if err != nil {
                log.Println("Err: ", err)
            } else {
                client.Send(msg)
            }
        }
    }()

    for {
        _, err := client.Next()
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            return
        }
    }
}

//parseLine returns an irc.Message object. If the line starts with a forward
//slash, everything after the '/' is converted directly to a server command
//If there is no slash, the first word is taken to be the channel or user to
//send a PRIVMSG to
func parseLine(line string) (msg irc.Message, err error) {
	if line[0] == '/' {
		msg = irc.NewMessage(line[1:])
	} else {
		splitlines := strings.SplitN(line, " ", 2)
		if len(splitlines) > 1 {
			msg = irc.PrivMessage(splitlines[0], splitlines[1])
		} else {
			err = errors.New("Unable to parse input")
		}
	}
	return
}
