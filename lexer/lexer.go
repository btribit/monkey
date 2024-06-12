package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // skipWhitespace is a helper function

	switch l.ch {
	case '=':
		if l.peekCharacter() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal} // EQ stands for equal
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekCharacter() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal} // NOT_EQ stands for not equal
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)

	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch) // LT stands for less than
	case '>':
		tok = newToken(token.GT, l.ch) // GT stands for greater than
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '@':
		tok = newToken(token.AT, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) { // isLetter is a helper function
			tok.Literal = l.readIdentifier()          // readIdentifier is a helper function
			tok.Type = token.LookupIdent(tok.Literal) // LookupIdent is a helper function
			return tok
		} else if isDigit(l.ch) { // isDigit is a helper function
			return l.readNumber() // readNumber is a helper function
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	tok.Line = l.line

	l.readChar()
	return tok
}

func (l *Lexer) readString() string { // readString is a helper function
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekCharacter() byte { // peekCharacter is a helper function
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() { // skipWhitespace is a helper function
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func (l *Lexer) readNumber() token.Token { // readNumber is a helper function
	var tok token.Token
	position := l.position
	tok.Type = token.INT
	for isDigit(l.ch) || isDecimal(l.ch) { // isDigit is a helper function
		if isDecimal(l.ch) {
			tok.Type = token.FLOAT
		}
		l.readChar()
	}
	tok.Literal = l.input[position:l.position]
	return tok
}

func isDigit(ch byte) bool { // isDigit is a helper function
	return '0' <= ch && ch <= '9'
}

func isDecimal(ch byte) bool { // isDecimal point helper function
	return ch == '.'
}

func isLetter(ch byte) bool { // isLetter is a helper function
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string { // readIdentifier is a helper function
	position := l.position
	for isLetter(l.ch) { // isLetter is a helper function
		l.readChar()
	}
	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
