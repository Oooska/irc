package irc

import (
	"bufio"
	"bytes"
	"io"
)

//Based on: https://golang.org/src/text/template/parse/lex.go

/*
The lexer scans from a stream of irc messages, and breaks them into
component parts. It parses one line at a time, returning a
token, representing what part of the message it is, and a string
literal.

The tokens are as follows:
BNF From: tools.ietf.org/html/rfc1459#section-2.3.1
			MESSAGE ::= [':' <prefix> <SPACE> ] <command> <params> <crlf>
			PREFIX  ::= <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
			COMMAND ::= <letter> { <letter> } | <number> <number> <number>
			PARAMS  ::= <SPACE> [ ':' <trailing> | <middle> <params> ]
			MIDDLE  ::= <Any *non-empty* sequence of octets not including SPACE
			//           or NUL or CR or LF, the first of which may not be ':'>
			TRAILING::= <Any, possibly *empty*, sequence of octets not including
			//           NUL or CR or LF>


TODO: The lexer currently does not signal trouble by returning an illegal token
when illegal characters are found
*/
type Token int

const (
	tokenIllegal Token = iota
	tokenEOF

	tokenPrefix
	tokenCommand
	tokenParam
	tokenTrailing

	tokenSpace //One or more space characters (no tabs/etc)
	tokenColon //:
	tokenEOL   //End of line (\r{\n})
)

const (
	byteSPACE = byte(0x20)
	byteNUL   = byte(0x00)
	byteCR    = '\r'
	byteLF    = '\n'
	byteBELL  = byte(0x07)
)

type lexer struct {
	s *bufio.Reader

	bufTokens  []Token
	bufLiteral []string
}

func newLexer(r io.Reader) *lexer {
	return &lexer{s: bufio.NewReader(r)}
}

//Scan() returns the next token, and the actual string value
func (l *lexer) NextItem() (token Token, literal string) {
	//ch, _, _ := l.next()
	if len(l.bufTokens) == 0 {
		l.tokenizeNextMessage()
	}

	if len(l.bufTokens) > 0 {
		token = l.bufTokens[0]
		literal = l.bufLiteral[0]

		l.bufTokens = l.bufTokens[1:]
		l.bufLiteral = l.bufLiteral[1:]
	} else {
		//No tokens to read from. Return nul byte
		token = tokenIllegal
		literal = string(byteNUL)
	}
	return
}

func (l *lexer) next() byte {
	ch, _ := l.s.ReadByte()
	return ch
}

func (l *lexer) unread() {
	l.s.UnreadByte()
}

func (l *lexer) peak() byte {
	ch, _ := l.s.ReadByte()
	l.s.UnreadByte()
	return ch
}

func (l *lexer) addToken(token Token, literal string) {
	l.bufTokens = append(l.bufTokens, token)
	l.bufLiteral = append(l.bufLiteral, literal)
}

//Tokenize methods actually go through and tokenize the next irc message to be read
//TODO: Tje
func (l *lexer) tokenizeNextMessage() {
	//MESSAGE ::= [':' <prefix> <SPACE> ] <command> <params> <crlf>
	if l.peak() == ':' {
		l.addToken(tokenColon, string(l.next()))
		l.addToken(tokenPrefix, l.scanChstring())
		l.addToken(l.scanSpaces())
	}

	//Parse message
	l.addToken(tokenCommand, l.scanChstring())
	l.addToken(l.scanSpaces())

	//Parse params
	for {
		if l.peak() == ':' {
			l.addToken(tokenColon, string(l.next()))
			l.addToken(l.scanTrailing())
			break
		} else if isEOL(l.peak()) {
			break
		} else {
			param := l.scanWord()
			l.addToken(tokenParam, param)
			l.addToken(l.scanSpaces())
		}
	}

	l.addToken(l.scanEOL())
}

//Reads a chstring
//<chstring> ::= <any 8bit code except SPACE, BELL, NUL, CR, LF and (',')>)
func (l *lexer) scanChstring() string {
	var buf bytes.Buffer
	for {
		if ch := l.next(); isChchar(ch) {
			buf.WriteByte(ch)
		} else {
			l.unread()
			return buf.String()
		}
	}
}

//Scans until encountering CRLF, or NUL
func (l *lexer) scanTrailing() (Token, string) {
	var buf bytes.Buffer
	for {
		if ch := l.next(); ch != byteCR && ch != byteLF && ch != byteNUL {
			buf.WriteByte(ch)
		} else {
			l.unread()
			return tokenTrailing, buf.String()
		}
	}
}

//Scans the next string until whitespace is encountered
func (l *lexer) scanWord() string {
	var buf bytes.Buffer
	for {
		if isNonwhite(l.peak()) {
			buf.WriteByte(l.next())
		} else {
			return buf.String()
		}
	}
}

//Consumes spaces (0x20) until the next non-space character
//Returns tokenSpace, and a string containing the same
//number of space characters consumed.
func (l *lexer) scanSpaces() (Token, string) {
	var buf bytes.Buffer
	for {
		if l.peak() == ' ' {
			buf.WriteByte(l.next())
		} else {
			return tokenSpace, buf.String()
		}
	}
}

//Scans EOL characters (\r\n)
func (l *lexer) scanEOL() (Token, string) {
	var buf bytes.Buffer
	if isEOL(l.peak()) {
		buf.WriteByte(l.next())
		if isEOL(l.peak()) {
			buf.WriteByte(l.peak())
		}
	}
	return tokenEOL, buf.String()
}

/*BNF From: tools.ietf.org/html/rfc1459#section-2.3.1
	   <chstring>   ::= <any 8bit code except SPACE, BELL, NUL, CR, LF and
	                     comma (',')>

	   Other parameter syntaxes are:

	   <user>       ::= <nonwhite> { <nonwhite> }
	   <letter>     ::= 'a' ... 'z' | 'A' ... 'Z'
	   <number>     ::= '0' ... '9'
	   <special>    ::= '-' | '[' | ']' | '\' | '`' | '^' | '{' | '}'
   	   <nonwhite>   ::= <any 8bit code except SPACE (0x20), NUL (0x0), CR
                     (0xd), and LF (0xa)>
*/
func isLetter(r byte) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z')
}

func isNonwhite(r byte) bool {
	return r != byteSPACE && r != byteNUL &&
		r != byteCR && r != byteLF
}

func isSpecial(r byte) bool {
	return r == '-' || r == '[' || r == ']' ||
		r == '\\' || r == '`' || r == '^' ||
		r == '{' || r == '}'
}

func isNumber(r byte) bool {
	return '0' <= r && r <= '9'
}

//An individual character of chstring
func isChchar(r byte) bool {
	return isNonwhite(r) && r != byteBELL && r != byte(',')
}

func isEOL(r byte) bool {
	return r == '\r' || r == '\n'
}
