// item represents a token or text string returned from the scanner.
type item struct {
    typ itemType // The type of this item.
    pos Pos      // The starting position, in bytes, of this item in the input string.
    val string   // The value of this item.
}
