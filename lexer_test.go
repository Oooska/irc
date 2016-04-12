package irc

import (
	"bytes"
	"testing"
)

var input = "PING :tepper.freenode.net\r\n" +
	":tepper.freenode.net 332 goirctest #mainehackerclub :MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM\r\n" +
	":tepper.freenode.net 353 goirctest = #mainehackerclub :goirctest +ubuntuguru " +
	/*  */ "+FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman " +
	/*  */ "@Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 " +
	/*  */ "+zgrep lambdabot +jrvc\r\n" +
	":wallyworld!~quassel@1-2-3-4.static.tpgi.com.au  QUIT :Remote host closed the connection\r\n" +
	":KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald PRIVMSG #go-nuts :https://golang.org/pkg/time/#Time.String\r\n" +
	":somenick!~@5-6-7-8.static.bgth.bz  QUIT\r\n"

var tokenToString = map[ircToken]string{
	tokenIllegal: "ILLEGAL",
	tokenEOF:     "EOF",
	tokenPrefix:  "PREFIX",

	tokenCommand:  "COMMAND",
	tokenParam:    "PARAM",
	tokenTrailing: "TRAILING",

	tokenSpace: "SPACE",

	tokenColon: "COLON",
	tokenEOL:   "EOL",
}

func TestNextItem(t *testing.T) {
	l := newLexer(bytes.NewBufferString(input))

	//PING :tepper.freenode.net
	checkNextItem(t, l, tokenCommand, "PING")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenTrailing, "tepper.freenode.net")
	checkNextItem(t, l, tokenEOL, "\r\n")

	//:tepper.freenode.net 332 goirctest #mainehackerclub :MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenPrefix, "tepper.freenode.net")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenCommand, "332")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenParam, "goirctest")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenParam, "#mainehackerclub")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenTrailing, "MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM")
	checkNextItem(t, l, tokenEOL, "\r\n")

	//:tepper.freenode.net 353 goirctest = #mainehackerclub :goirctest +ubuntuguru +FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman @Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 +zgrep lambdabot +jrvc
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenPrefix, "tepper.freenode.net")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenCommand, "353")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenParam, "goirctest")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenParam, "=")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenParam, "#mainehackerclub")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenTrailing, "goirctest +ubuntuguru +FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman @Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 +zgrep lambdabot +jrvc")
	checkNextItem(t, l, tokenEOL, "\r\n")

	//:wallyworld!~quassel@5-5-5-5.static.tpgi.com.au  QUIT :Remote host closed the connection\r\n
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenPrefix, "wallyworld!~quassel@1-2-3-4.static.tpgi.com.au")
	checkNextItem(t, l, tokenSpace, "  ") //Test extra space
	checkNextItem(t, l, tokenCommand, "QUIT")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenTrailing, "Remote host closed the connection")
	checkNextItem(t, l, tokenEOL, "\r\n")

	//:KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald PRIVMSG #go-nuts :https://golang.org/pkg/time/#Time.String\r\n
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenPrefix, "KirkMcDonald!~Kirk@python/site-packages/KirkMcDonald")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenCommand, "PRIVMSG")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenParam, "#go-nuts")
	checkNextItem(t, l, tokenSpace, " ")
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenTrailing, "https://golang.org/pkg/time/#Time.String")
	checkNextItem(t, l, tokenEOL, "\r\n")

	//":somenick!~@5-6-7-8.static.bgth.bz  QUIT\r\n"
	checkNextItem(t, l, tokenColon, ":")
	checkNextItem(t, l, tokenPrefix, "somenick!~@5-6-7-8.static.bgth.bz")
	checkNextItem(t, l, tokenSpace, "  ")
	checkNextItem(t, l, tokenCommand, "QUIT")
	checkNextItem(t, l, tokenEOL, "\r\n")

}

func checkNextItem(t *testing.T, l *lexer, expectedToken ircToken, expectedLiteral string) {
	actualToken, actualLiteral := l.NextItem()
	stringLiteral := string(actualLiteral)
	if expectedToken != actualToken || expectedLiteral != stringLiteral {
		t.Errorf("Expected: <%s, %s>. Received: <%s, %s>",
			tokenToString[expectedToken], expectedLiteral, tokenToString[actualToken], actualLiteral)
	}
}
