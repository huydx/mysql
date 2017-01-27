// parse is the top-level parser for a template, essentially the same
// as itemList except it also parses {{define}} actions.
// It runs to EOF.
func (t *Tree) parse() (next Node) {
    t.Root = t.newList(t.peek().pos)
    for t.peek().typ != itemEOF {
        if t.peek().typ == itemLeftDelim {
            delim := t.next()
            pp.Println(fmt.Sprintf("token after next delim: %v", t.token))
            if t.nextNonSpace().typ == itemDefine {
                newT := New("definition") // name will be updated once we know it.
                newT.text = t.text
                newT.ParseName = t.ParseName
                newT.startParse(t.funcs, t.lex, t.treeSet)
                newT.parseDefinition()
                continue
            }
            t.backup2(delim)
        }
        switch n := t.textOrAction(); n.Type() {
        case nodeEnd, nodeElse:
            t.errorf("unexpected %s", n)
        default:
            t.Root.append(n)
        }
    }
    return nil
}
