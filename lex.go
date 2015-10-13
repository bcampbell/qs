package qs

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TODO: handle escaping special characters:

type tokType int

const (
	tEOF tokType = iota
	tERROR
	tLITERAL
	tQUOTED
	tPLUS
	tMINUS
	tCOLON
	tEQUAL
	tGREATER
	tLESS
	tOR
	tAND
	tNOT
	tLPAREN
	tRPAREN
	tLSQUARE
	tRSQUARE
	tLBRACE
	tRBRACE
	tTO
	tBOOST
	tFUZZY
)

func (t tokType) String() string {

	tokTypes := map[tokType]string{
		tEOF:     "tEOF",
		tERROR:   "tERROR",
		tLITERAL: "tLITERAL",
		tQUOTED:  "tQUOTED",
		tOR:      "tOR",
		tAND:     "tAND",
		tNOT:     "tNOT",
		tTO:      "tTO",
		tPLUS:    "tPLUS",
		tMINUS:   "tMINUS",
		tCOLON:   "tCOLON",
		tEQUAL:   "tEQUAL",
		tGREATER: "tGREATER",
		tLESS:    "tLESS",
		tLPAREN:  "tLPAREN",
		tRPAREN:  "tRPAREN",
		tLSQUARE: "tLSQUARE",
		tRSQUARE: "tRSQUARE",
		tLBRACE:  "tLBRACE",
		tRBRACE:  "tRBRACE",
		tBOOST:   "tBOOST",
		tFUZZY:   "tFUZZY",
	}
	return tokTypes[t]
}

// some single-rune tokens
var singles = map[rune]tokType{
	'(': tLPAREN,
	')': tRPAREN,
	'[': tLSQUARE,
	']': tRSQUARE,
	'{': tLBRACE,
	'}': tRBRACE,
	':': tCOLON,
	'+': tPLUS,
	'-': tMINUS,
	'=': tEQUAL,
	'>': tGREATER,
	'<': tLESS,
}

type token struct {
	typ tokType
	val string
	pos int
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input   string
	tokens  []token
	pos     int
	prevpos int
	start   int
}

// lex takes an input string and breaks it up into an array of tokens.
// The last token will be an tEOF, unless an error occurs, in which case
// it will be a tERROR.
func lex(input string) []token {
	l := &lexer{
		input:  input,
		tokens: []token{},
	}
	// run state machine - each state returns the next state, or nil when finished
	for state := lexDefault; state != nil; {
		state = state(l)
	}
	return l.tokens
}

func (l *lexer) next() rune {
	l.prevpos = l.pos
	if l.eof() {
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	return r
}
func (l *lexer) eof() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) ignore() {
	l.start = l.pos
}
func (l *lexer) backup() {
	l.pos = l.prevpos
}
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) emit(t tokType) {
	l.tokens = append(l.tokens, token{t, l.input[l.start:l.pos], l.start})
	l.start = l.pos
}

func (l *lexer) emitError(msg string) {
	l.tokens = append(l.tokens, token{tERROR, msg, l.start})
	l.start = l.pos
}

func lexDefault(l *lexer) stateFn {
	// skip space
	for {
		if l.eof() {
			l.emit(tEOF)
			return nil
		}
		r := l.next()
		if !unicode.IsSpace(r) {
			l.backup()
			l.ignore()
			if r == '~' || r == '^' {
				return lexSuffix
			}
			if typ, got := singles[r]; got {
				l.next()
				l.emit(typ)
				return lexDefault
			}

			if r == '"' || r == '\'' {
				return lexQuoted
			}

			return lexText
		}
	}
}

func lexText(l *lexer) stateFn {
	// non-space characters which terminate a literal
	//stopChars := `&|:(){}[]^~`
	stopChars := `:(){}[]^~`

	for {
		if l.eof() {
			break
		}
		r := l.next()
		if unicode.IsSpace(r) || strings.ContainsRune(stopChars, r) {
			l.backup()
			break
		}
	}

	switch l.input[l.start:l.pos] {
	case "OR":
		l.emit(tOR)
	case "AND":
		l.emit(tAND)
	case "NOT":
		l.emit(tNOT)
	case "TO":
		l.emit(tTO)
	default:
		l.emit(tLITERAL)
	}

	return lexDefault
}

func lexQuoted(l *lexer) stateFn {
	q := l.next()
	for {
		if l.eof() {
			l.emitError("unclosed quote")
			return nil
		}
		r := l.next()
		if r == q {
			break
		}
	}
	l.emit(tQUOTED)
	return lexDefault
}

func lexSuffix(l *lexer) stateFn {
	kind := l.next() // '^' or '~'

	// number (optional)
	for {
		if l.eof() {
			break
		}
		r := l.next()
		if unicode.IsSpace(r) {
			l.backup()
			break
		}
		if !strings.ContainsRune("0123456789.", r) {
			l.emitError("bad number")
		}
	}

	switch kind {
	case '~':
		l.emit(tFUZZY)
	case '^':
		l.emit(tBOOST)
	}
	return lexDefault
}
