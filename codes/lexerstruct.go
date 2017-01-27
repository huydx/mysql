type lexer struct {
    name       string    // the name of the input; used only for error reports
    input      string    // the string being scanned
    leftDelim  string    // start of action
    rightDelim string    // end of action
    state      stateFn   // the next lexing function to enter
    pos        Pos       // current position in the input
    start      Pos       // start position of this item
    width      Pos       // width of last rune read from input
    lastPos    Pos       // position of most recent item returned by nextItem
    items      chan item // channel of scanned items
    parenDepth int       // nesting depth of ( ) exprs
}

