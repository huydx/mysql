// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
    l.items <- item{t, l.start, l.input[l.start:l.pos]}
    l.start = l.pos
}

func lex(name, input, left, right string) *lexer {
    if left == "" {
        left = leftDelim
    }
    if right == "" {
        right = rightDelim
    }
    l := &lexer{
        name:       name,
        input:      input,
        leftDelim:  left,
        rightDelim: right,
        items:      make(chan item),
    }
    go l.run()
    return l
}
