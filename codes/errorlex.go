// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
    l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
    return nil
}
