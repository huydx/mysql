// lexRawQuote scans a raw quoted string.
func lexRawQuote(l *lexer) stateFn {
Loop:
    for {
        switch l.next() {
        case eof:
            return l.errorf("unterminated raw quoted string")
        case '`':
            break Loop
        }
    }
    l.emit(itemRawString)
    return lexInsideAction
}
