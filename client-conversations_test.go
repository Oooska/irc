package irc

import "testing"

func TestConversations(t *testing.T) {
	convos := newConversations(5)

	convos.Add("#chan1", "Message 1")
	convos.Add("#chan1", "Message 2")
	convos.Add("#chan2", "Message 1 in #chan2")

	messages := convos.Messages("#chan1")
	if len(messages) != 2 {
		t.Errorf("#chan1 was expected to return 2 messages. Returned %d", len(messages))
	}
	if messages[0] != "Message 1" {
		t.Errorf("First message returned should be \"Message1\", Returned: \"%s\"", messages[0])
	}
	if messages[1] != "Message 2" {
		t.Errorf("First message returned should be \"Message1\", Returned: \"%s\"", messages[0])
	}

	messages = convos.Messages("#chan2")
	if len(messages) != 1 {
		t.Errorf("#chan2 was expected to return 1 messages. Returned %d", len(messages))
	}
	if messages[0] != "Message 1 in #chan2" {
		t.Errorf("First message returned should be \"Message 1 in #chan2\", Returned: \"%s\"", messages[0])
	}

	convos.Add("#chan1", "Message 3")
	convos.Add("#chan1", "Message 4")
	convos.Add("#chan1", "Message 5")
	messages = convos.Messages("#chan1")
	if len(messages) != 5 {
		t.Errorf("#chan1 was expected to return 5 messages. Returned %d", len(messages))
	}
	if messages[0] != "Message 1" {
		t.Errorf("First message returned should be \"Message 1\", Returned: \"%s\"", messages[0])
	}

	//Make sure limit is working correctly
	convos.Add("#chan1", "Message 6")
	messages = convos.Messages("#chan1")
	if len(messages) != 5 {
		t.Errorf("#chan1 was expected to return 5 messages. Returned %d", len(messages))
	}
	if messages[0] != "Message 2" {
		t.Errorf("First message returned should be \"Message 2\", Returned: \"%s\"", messages[0])
	}

	convos.Add("#chan1", "Message 7")
	convos.Add("#chan1", "Message 8")
	convos.Add("#chan1", "Message 9")
	convos.Add("#chan1", "Message 10")
	messages = convos.Messages("#chan1")
	if len(messages) != 5 {
		t.Errorf("#chan1 was expected to return 5 messages. Returned %d", len(messages))
	}
	if messages[0] != "Message 6" {
		t.Errorf("First message returned should be \"Message 6\", Returned: \"%s\"", messages[0])
	}

}
