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

	fmt.Println("IRC Test.")

	conn, err := irc.NewIRCConnection(*address, *ssl)
	if err != nil {
		panic(err)
	}

	conn.Write(irc.UserMessage(*username, "hostname", "servrename", "real name"))
	conn.Write(irc.NickMessage(*nick))
	conn.Write(irc.NewMessage("JOIN #go_Test"))

	go func() {
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
				log.Println("YOU: ", msg)
				conn.Write(msg)
			}
		}
	}()

	var msg irc.Message

	// handle incomming messages
	for {

		msg, err = conn.Read()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		if len(msg) >= 4 && msg[0:4] == "PING" {
			var end string = msg[4:].String()
			conn.Write(irc.NewMessage("PONG" + end))
		}

		fmt.Println(msg)
	}
}

//parseLine returns an irc.Message object. If the line starts with a forward
//slash, everything after the '/' is converted directly to a server command
//If there is no slash, the first word is taken to be the channel or user to
//send a PRIVMSG to
func parseLine(line string) (msg irc.Message, err error) {
	if line[0] == '/' {
		msg = irc.NewMessage(line[1:]) //TODO Parse actual command
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
