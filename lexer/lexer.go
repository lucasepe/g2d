package lexer

// Package lexer implements the lexical analysis that is used to
// transform the source code input into a stream of tokens for
// parsing by the parser.

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/lucasepe/g2d/token"
)

// Lexer represents the lexer and contains the source input and internal state
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	prevCh       byte // previous char read

	// map of input line boundaries used by linePosition() for error location
	lineMap [][2]int // array of [begin, end] pairs: [[0,12], [13,22], [23,33] ... ]
}

// New returns a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// map the input line boundaries for CurrentLine()
	l.buildLineMap()

	l.readChar()
	return l
}

// buildLineMap creates map of input line boundaries used by LinePosition() for error location
func (l *Lexer) buildLineMap() {
	begin := 0
	idx := 0
	for i, ch := range l.input {
		idx = i
		if ch == '\n' {
			l.lineMap = append(l.lineMap, [2]int{begin, idx})
			begin = idx + 1
		}
	}
	// last line
	l.lineMap = append(l.lineMap, [2]int{begin, idx + 1})
}

// CurrentPosition returns l.position
func (l *Lexer) CurrentPosition() int {
	return l.position
}

// linePosition (pos) returns lineNum, begin, end
func (l *Lexer) linePosition(pos int) (int, int, int) {
	idx := 0
	begin := 0
	end := 0
	for i, tuple := range l.lineMap {
		idx = i
		begin, end = tuple[0], tuple[1]
		if pos >= begin && pos <= end {
			break
		}
	}
	lineNum := idx + 1
	return lineNum, begin, end
}

// ErrorLine (pos) returns lineNum, column, errorLine
func (l *Lexer) ErrorLine(pos int) (int, int, string) {
	lineNum, begin, end := l.linePosition(pos)
	errorLine := l.input[begin:end]
	column := pos - begin + 1
	return lineNum, column, string(errorLine)
}

func (l *Lexer) readChar() {
	l.prevCh = l.ch
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// NextToken returns the next token read from the input stream
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '#':
		tok = token.New(token.COMMENT, l.position, l.readLine())
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.EQ, l.position, literal)
		} else {
			tok = token.New(token.ASSIGN, l.position, string(l.ch))
		}
	case '+':
		tok = token.New(token.PLUS, l.position, string(l.ch))
	case '-':
		tok = token.New(token.MINUS, l.position, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.NEQ, l.position, literal)
		} else {
			tok = token.New(token.NOT, l.position, string(l.ch))
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar() // skip over the '/'
			tok = token.New(token.COMMENT, l.position, l.readLine())
		} else {
			tok = token.New(token.DIVIDE, l.position, string(l.ch))
		}
	case '*':
		tok = token.New(token.MULTIPLY, l.position, string(l.ch))
	case '%':
		tok = token.New(token.MODULO, l.position, string(l.ch))
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.AND, l.position, literal)
		} else {
			tok = token.New(token.ILLEGAL, l.position, string(l.ch))
			//tok = token.New(token.BitwiseAND, l.position, string(l.ch))
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.OR, l.position, literal)
		} else {
			tok = token.New(token.ILLEGAL, l.position, string(l.ch))
			//tok = token.New(token.BitwiseOR, l.position, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.LTE, l.position, "<=")
		} else {
			tok = token.New(token.LT, l.position, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.GTE, l.position, ">=")
		} else {
			tok = token.New(token.GT, l.position, string(l.ch))
		}
	case ';':
		tok = token.New(token.SEMICOLON, l.position, string(l.ch))
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.New(token.BIND, l.position, literal)
		} else {
			tok = token.New(token.COLON, l.position, string(l.ch))
		}
	case ',':
		tok = token.New(token.COMMA, l.position, string(l.ch))
	case '.':
		tok = token.New(token.DOT, l.position, string(l.ch))
	case '(':
		tok = token.New(token.LPAREN, l.position, string(l.ch))
	case ')':
		tok = token.New(token.RPAREN, l.position, string(l.ch))
	case '{':
		tok = token.New(token.LBRACE, l.position, string(l.ch))
	case '}':
		tok = token.New(token.RBRACE, l.position, string(l.ch))
	case '[':
		tok = token.New(token.LBRACKET, l.position, string(l.ch))
	case ']':
		tok = token.New(token.RBRACKET, l.position, string(l.ch))
	case 0:
		tok = token.New(token.EOF, l.position, "")
	case '"':
		str, err := l.readString()
		if err != nil {
			tok = token.New(token.ILLEGAL, l.position, string(l.prevCh))
		} else {
			tok = token.New(token.STRING, l.position, str)
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Position = l.position
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			integer := l.readNumber()
			if l.ch == '.' && isDigit(l.peekChar()) {
				// OK here we think we've got a float.
				l.readChar()
				fraction := l.readNumber()

				tok.Type = token.FLOAT
				tok.Position = l.position
				tok.Literal = fmt.Sprintf("%s.%s", integer, fraction)
			} else {
				tok.Type = token.INT
				tok.Position = l.position
				tok.Literal = integer
			}
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.position, string(l.ch))
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readLine() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\r' || l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, error) {
	b := &strings.Builder{}
	for {
		l.readChar()

		// Support some basic escapes like \"
		if l.ch == '\\' {
			switch l.peekChar() {
			case '"':
				b.WriteByte('"')
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			case 'x':
				// Skip over the the '\\', 'x' and the next two bytes (hex)
				l.readChar()
				l.readChar()
				l.readChar()
				src := string([]byte{l.prevCh, l.ch})
				dst, err := hex.DecodeString(src)
				if err != nil {
					return "", err
				}
				b.Write(dst)
				continue
			}

			// Skip over the '\\' and the matched single escape char
			l.readChar()
			continue
		} else {
			if l.ch == '"' || l.ch == 0 {
				break
			}
		}

		b.WriteByte(l.ch)
	}

	return b.String(), nil
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
