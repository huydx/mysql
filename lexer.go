package main

import (
	"fmt"
	_ "github.com/k0kubun/pp"
	"strings"
	"unicode/utf8"
)

type itemType int
type Pos int

// item represents a token or text string returned from the scanner.
type item struct {
	typ itemType // The type of this item.
	pos Pos      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}


func (i item) String() string {
	return fmt.Sprintf("%s, %d, %s", i.typ, i.pos, i.val)
}

const eof = -1

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        Pos       // current position in the input
	start      Pos       // start position of this item
	width      Pos       // width of last rune read from input
	lastPos    Pos       // position of most recent item returned by nextItem
	items      chan item // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
}


// lex creates a new scanner for the input string.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	i := item{t, l.start, l.input[l.start:l.pos]}
	fmt.Println(i)
	l.items <- i
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexStmt; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

/////////////////
const (
	itemError itemType = iota
	itemSelect
	itemText
	itemComma
	itemFrom
	itemEOF
	itemAsterisk
)

var itemTypeStr = map[itemType]string{
	itemSelect:   "select",
	itemText:     "text",
	itemComma:    "comma",
	itemFrom:     "from",
	itemEOF:      "eof",
	itemAsterisk: "*",
}

func (i itemType) String() string {
	if s, ok := itemTypeStr[i]; ok {
		return s
	}
	return fmt.Sprintf("<item %d>", i)
}

var key = map[string]itemType{
	// Operators.
	"select": itemSelect,
	"from":   itemFrom,
	"*":      itemAsterisk,
	",":      itemComma,
}


func lexStmt(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case isSpace(r):
		return lexSpace
	default:
		return lexKeywordOrText
	}
	return nil
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexKeywordOrText
}



func lexKeywordOrText(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isSpace(r):
			return lexSpace
		case r == eof:
			l.emit(itemEOF)
			return nil
		case r == '*':
			l.emit(itemAsterisk)
			break
		case r == ',':
			l.emit(itemComma)
			break
		default:
			n := l.peek()
			word := l.input[l.start:l.pos]
			if kw, ok := key[strings.ToLower(word)]; ok {
				l.emit(kw)
				break
			}

			if n == ',' || isSpace(n) || n == eof {
				l.emit(itemText)
				break
			}
		}
	}
	return lexStmt
}
