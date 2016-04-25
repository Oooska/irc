package irc

import (
	"testing"
)

const userlist = ":goirctest +ubuntuguru " +
	/*   */ "+FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman " +
	/*   */ "@Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 " +
	/*   */ "+zgrep lambdabot +jrvc"

var parserInput = []message{
	message{message: "PING :tepper.freenode.net", command: "PING", params: []string{":tepper.freenode.net"}, trailing: "tepper.freenode.net"},
	message{message: ":tepper.freenode.net 332 goirctest #mainehackerclub :MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM\\",
		prefix: ":tepper.freenode.net", command: "332", server: "tepper.freenode.net",
		params: []string{"goirctest", "#mainehackerclub", ":MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM\\"}, trailing: "MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM\\"},
	message{message: ":tepper.freenode.net 353 goirctest = #mainehackerclub " + userlist, prefix: ":tepper.freenode.net", command: "353", server: "tepper.freenode.net",
		params: []string{"goirctest", "=", "#mainehackerclub", userlist}, trailing: userlist[1:]},
	message{message: ":wallyworld!~quassel@1-2-3-4.static.tpgi.com.au  QUIT :Remote host closed the connection", prefix: ":wallyworld!~quassel@1-2-3-4.static.tpgi.com.au",
		command: "QUIT", nick: "wallyworld", user: "~quassel", host: "1-2-3-4.static.tpgi.com.au", params: []string{":Remote host closed the connection"},
		trailing: "Remote host closed the connection"},
	message{message: ":KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald PRIVMSG #go-nuts :https://golang.org/pkg/time/#Time.String",
		prefix: ":KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald", nick: "KirkMcDonald", user: "~Kirk", host: "python/site-packages/KirkMcDonald",
		command: "PRIVMSG", params: []string{"#go-nuts", ":https://golang.org/pkg/time/#Time.String"}, trailing: "https://golang.org/pkg/time/#Time.String"},
	message{message: ":somenick!~@5-6-7-8.static.bgth.bz  QUIT", prefix: ":somenick!~@5-6-7-8.static.bgth.bz",
		nick: "somenick", user: "~", host: "5-6-7-8.static.bgth.bz", command: "QUIT"},
}

func TestParseString(t *testing.T) {
	for j, expected := range parserInput {
		actual := NewMessage(expected.Message())

		if actual.Message() != expected.Message() {
			t.Errorf("input[%d]: Message field not parsed correctly. Expected: %s. Received: %s", j, expected.Message(), actual.Message())
		}

		if actual.Prefix() != expected.Prefix() {
			t.Errorf("input[%d]: Prefix field not parsed correctly. Expected: %s. Received: %s", j, expected.Prefix(), actual.Prefix())
		}

		if actual.Command() != expected.Command() {
			t.Errorf("input[%d]: Command field not parsed correctly. Expected: %s. Received: %s", j, expected.Command(), actual.Command())
		}

		if len(actual.Params()) != len(expected.Params()) {
			t.Errorf("input[%d]: Unequal number of params. Expected: %+v. Received: %+v", j, expected.Params(), actual.Params())
		} else {
			for k := range actual.Params() {
				if actual.Params()[k] != expected.Params()[k] {
					t.Errorf("input[%d]: Param[%d] is not equal. Expected: %s. Received: %s", j, k, expected.Params()[k], actual.Params()[k])
				}
			}
		}

		if actual.Nick() != expected.Nick() {
			t.Errorf("input[%d]: Nick field not parsed correctly. Expected: %s. Received: %s", j, expected.Nick(), actual.Nick())
		}

		if actual.User() != expected.User() {
			t.Errorf("input[%d]: User field not parsed correctly. Expected: %s. Received: %s", j, expected.User(), actual.User())
		}

		if actual.Host() != expected.Host() {
			t.Errorf("input[%d]: Host field not parsed correctly. Expected: %s. Received: %s", j, expected.Host(), actual.Host())
		}

		if actual.Trailing() != expected.Trailing() {
			t.Errorf("input[%d]: Host field not parsed correctly. Expected: %s. Received: %s", j, expected.Trailing(), actual.Trailing())
		}
	}

}
