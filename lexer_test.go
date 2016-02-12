package irc

import (
	"bytes"
	"testing"
)

var input = `PING :tepper.freenode.net
:tepper.freenode.net 332 goirctest #mainehackerclub :MHC@MSF - 3-18-2016 @ Bangor Y 4-8PM
:tepper.freenode.net 353 goirctest = #mainehackerclub :goirctest +ubuntuguru +FatalNIX MaineHackerBot +hpcr2013 +T-800 +infina Oooska +AaronBallman @Derrick[afk] dual +Church- +arschmitz +wrexem +powellc +jeepingben +nh_99 +zgrep lambdabot +jrvc
`

var tokenToString = map[Token]string{
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
	checkNextItem(t, l, tokenEOL, "\n")

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
	checkNextItem(t, l, tokenEOL, "\n")

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

}

func checkNextItem(t *testing.T, l *lexer, expectedToken Token, expectedLiteral string) {
	actualToken, actualLiteral := l.NextItem()
	if expectedToken != actualToken || expectedLiteral != actualLiteral {
		t.Errorf("Expected: <%s, %s>. Received: <%s, %s>",
			tokenToString[expectedToken], expectedLiteral, tokenToString[actualToken], actualLiteral)
	}
}
