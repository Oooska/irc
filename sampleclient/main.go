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

//A barebones IRC 'client' in the loosest sense of the word.
//Takes input from console. If command starts with a '/', everything after is sent as a raw IRC command.
//Otherwise the first argument is considered the channel/username, and the rest of the line is the message to send
// as a privmsg.
func main() {
	flag.Parse()

	fmt.Print("Simple Text-Based IRC Test Client\n\n")

	fmt.Printf("Connecting to %s . . . \n", *address)

	//LogClientHandler will handle printing out to stdio unless we change the default logger
	client, err := irc.NewClient(*address, *ssl, irc.LogHandler)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	fmt.Print("Connected.\n\n")

	client.Send(irc.UserMessage(*username, "host", "domain", "realname"))
	client.Send(irc.NickMessage(*nick))
	client.Send(irc.JoinMessage("#go_test"))

	//Listen for input.
	go readInput(client)

	for { //Continuously read from the client until an error occurs
		_, err := client.Read()
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			fmt.Print("Exiting...")
			return
		}
	}
}

//readInput continuously reads line from stdin.
func readInput(client irc.Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Cannot read from stdin: %s", err.Error())
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		msg, err := parseLine(line)
		if err != nil {
			log.Printf("Err: %s\n", err.Error())
		} else {
			client.Send(msg)
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
