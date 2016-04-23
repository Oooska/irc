package irc

import (
	"bufio"
	"log"
	"net"
	"testing"
)

func getListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Unable to listen on port 8080: %s", err.Error())
	}
	return l
}

//TODO: Test SSL
func TestNewConnection(t *testing.T) {
	l := getListener()
	go func() {
		_, err := l.Accept()
		if err != nil {
			t.Errorf("Unable to accept connection from IRC client: %s", err.Error())
		}
	}()

	ircConn, err := NewConnection("127.0.0.1:8080", false)
	if err != nil {
		t.Errorf("Unable to connect to IRC server.")
	}
	ircConn.Close()
	l.Close()

	//Listener is now off. NewConnection should return an error
	ircConn, err = NewConnection("127.0.0.1:8080", false)
	if err == nil {
		t.Error("Connection failed, but no error returned.")
	}
}

func TestRead(t *testing.T) {
	l := getListener()
	defer l.Close()
	go func() {
		lconn, err := l.Accept()
		if err != nil {
			t.Fatalf("Unable to accept connection from IRC client: %s", err.Error())
		}
		lconn.Write([]byte("Message 1\r\n"))
		lconn.Write([]byte("Message 2\r\n"))
		lconn.Close()
	}()

	ircConn, err := NewConnection("127.0.0.1:8080", false)
	if err != nil {
		t.Errorf("Unable to connect to IRC server.")
	}

	msg, err := ircConn.Read()
	if err != nil {
		t.Errorf("Error returned while reading from server: %s", err.Error())
	}
	if msg.Message != "Message 1" {
		t.Errorf(`Read() did not return the expected message. Expected: "Message 1", Received: "%s"`, msg.Message)
	}

	msg, err = ircConn.Read()
	if err != nil {
		t.Errorf("Error returned while reading from server: %s", err.Error())
	}
	if msg.Message != "Message 2" {
		t.Errorf(`Read() did not return the expected message. Expected: "Message 2", Received: "%s"`, msg.Message)
	}

	msg, err = ircConn.Read()
	if err == nil {
		t.Errorf("Read from a closed connection. No error returned where one was expected.")
	}
}

func TestWrite(t *testing.T) {
	l := getListener()
	defer l.Close()
	go func() {
		lconn, err := l.Accept()
		if err != nil {
			t.Fatalf("Unable to accept connection from IRC client: %s", err.Error())
		}
		s := bufio.NewScanner(lconn)
		ok := s.Scan()
		if !ok {
			t.Fatalf("Unable to read data from irc.Conn. Error: %s", s.Err())
		}
		line := s.Text()
		if line != "Message 1" {
			t.Errorf(`Write() did not send the expected value. Expected: "Message 1", Received : "%s"`, line)
		}

		ok = s.Scan()
		if !ok {
			t.Fatalf("Unable to read data from irc.Conn. Error: %s", s.Err())
		}
		line = s.Text()
		if line != "Message 2" {
			t.Errorf(`Write() did not send the expected value. Expected: "Message 2", Received : "%s"`, line)
		}
		lconn.Close()
	}()

	ircConn, err := NewConnection("127.0.0.1:8080", false)
	if err != nil {
		t.Errorf("Unable to connect to IRC server.")
	}
	err = ircConn.Write(NewMessage("Message 1"))
	if err != nil {
		t.Errorf("Unexpected error when sending message: %s", err.Error())
	}
	err = ircConn.Write(NewMessage("Message 2"))
	if err != nil {
		t.Errorf("Unexpected error when sending message: %s", err.Error())
	}

	ircConn.Close()
	err = ircConn.Write(NewMessage("This message should fail"))
	if err == nil {
		t.Error("Sent message on closed connection. No error Received.")
	}
}

func TestHandlers(t *testing.T) {
	incomingOnAll := 0
	incomingOnMessage := 0
	outgoingOnAll := 0
	outgoingOnMessage := 0
	bothOnAll := 0
	bothOnMessage := 0

	l := getListener()
	defer l.Close()

	ircConn, err := NewConnection("127.0.0.1:8080", false)
	if err != nil {
		t.Errorf("Unable to connect to IRC server.")
	}

	go func() {
		lconn, _ := l.Accept()
		lconn.Write([]byte("Response 1\r\n"))
		lconn.Write([]byte("DiffResponse 2\r\n"))
		lconn.Close()
	}()

	ircConn.AddHandler(Incoming, func(msg Message) {
		incomingOnAll++
	})

	ircConn.AddHandler(Incoming, func(msg Message) {
		incomingOnMessage++
	}, "Response")

	ircConn.AddHandler(Outgoing, func(msg Message) {
		outgoingOnAll++
	})

	ircConn.AddHandler(Outgoing, func(msg Message) {
		outgoingOnMessage++
	}, "Message")

	ircConn.AddHandler(Both, func(msg Message) {
		bothOnAll++
	})

	ircConn.AddHandler(Both, func(msg Message) {
		bothOnMessage++
	}, "Message", "Response")

	ircConn.Write(NewMessage("Message 1"))
	ircConn.Read()
	ircConn.Write(NewMessage("DiffMessage 2"))
	ircConn.Read()
	ircConn.Close()
	ircConn.Write(NewMessage("Message 3"))
	ircConn.Read()

	if incomingOnAll != 2 {
		t.Errorf("Error /w incoming handler listening to all messages. Expected: 2 calls, Received: %d calls", incomingOnAll)
	}

	if incomingOnMessage != 1 {
		t.Errorf("Error /w incoming handler listening to specific message. Expected: 1 calls, Received: %d calls", incomingOnMessage)
	}

	if outgoingOnAll != 2 {
		t.Errorf("Error /w outgoing handler listening to all messages. Expected: 2 calls, Received: %d calls", outgoingOnAll)
	}

	if outgoingOnMessage != 1 {
		t.Errorf("Error /w outgoing handler listening to specific message. Expected: 1 calls, Received; %d calls", outgoingOnMessage)
	}

	if bothOnAll != 4 {
		t.Errorf("Error /w handler listening to all message in both directions. Expected: 4 calls, Received: %d calls", bothOnAll)
	}

	if bothOnMessage != 2 {
		t.Errorf("Error /w handler listening to specific messages in both directions. Expected: 2 calls, Received: %d calls", bothOnMessage)
	}
}
