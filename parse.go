package main

import (
	"fmt"
)

type parser struct {
	lex       *lexer
	token     [3]item
	peekCount int
}

// ParseErr wraps a parsing error with line and position context.
// If the parsing input was a single line, line will be 0 and omitted
// from the error string.
type ParseErr struct {
	Line, Pos int
	Err       error
}

func (e *ParseErr) Error() string {
	if e.Line == 0 {
		return fmt.Sprintf("parse error at char %d: %s", e.Pos, e.Err)
	}
	return fmt.Sprintf("parse error at line %d, char %d: %s", e.Line, e.Pos, e.Err)
}

// newParser returns a new parser.
func newParser(input string) *parser {
	p := &parser{
		lex: lex(input),
	}
	return p
}

func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		t := p.lex.nextItem()
		p.token[0] = t
	}
	return p.token[p.peekCount]
}

// peek returns but does not consume the next token.
func (p *parser) peek() item {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1

	t := p.lex.nextItem()
	p.token[0] = t
	return p.token[0]
}


///////////////////////////

func ParseSql(input string) (sql SelectStmt, e error) {
	p := newParser(input)
	sql = SelectStmt{}
	for p.peek().typ != itemEOF {
		if p.peek().typ == itemSelect {
			p.next()
			fields, err := p.parseFields()
			if err != nil {
				e = err
				return
			}

			sql.Fields = fields

			if p.peek().typ == itemFrom {
				p.next()
				tableName := p.parseTableName()
				sql.TableName = tableName
				return
			}
		}

		/**
		if p.peek().type == itemDelete {
		}
		 */
	}

	e = fmt.Errorf("wrong syntax")
	return
}

func (p *parser) parseFields() ([]string, error) {
	fields := make([]string, 0)
	hasAnyField := true
	endWithComma := false
	for {
		if p.peek().typ == itemEOF || p.peek().typ == itemFrom {
			return fields, nil
		} else {
			if p.peek().typ == itemComma {
				p.next()
				continue
			} else {
				fields = append(fields, p.next().val)
				endWithComma = false
			}
		}
	}

	if endWithComma || !hasAnyField {
		return fields, fmt.Errorf("invalid state of fields")
	}

	return fields, nil
}

func (p *parser) parseTableName() string {
	return p.peek().val
}
