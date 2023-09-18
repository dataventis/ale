package lang

const (
	structure    = `(){}\[\]\s\"`
	prefixChar   = "`,~@"
	nonStructure = `[^` + structure + `]`
	idStart      = `[^` + structure + prefixChar + `]`
	idCont       = nonStructure + `*`
	numTail      = idStart + `*`

	Keyword = `:` + nonStructure + `+`
	ID      = idStart + idCont

	Comment     = `;[^\n]*([\n]|$)`
	NewLine     = `(\r\n|[\n\r])`
	Whitespace  = `[\t\f ]+`
	ListStart   = `\(`
	ListEnd     = `\)`
	VectorStart = `\[`
	VectorEnd   = `]`
	ObjectStart = `{`
	ObjectEnd   = `}`
	Quote       = `'`
	SyntaxQuote = "`"
	Splice      = `,@`
	Unquote     = `,`
	Pattern     = `~`
	Dot         = `\.`

	String = `(")(?P<s>(\\\\|\\"|\\[^\\"]|[^"\\])*)("?)`

	Ratio = `[+-]?(0|[1-9]\d*)/[1-9]\d*` + numTail

	Float = `[+-]?((0|[1-9]\d*)\.\d+([eE][+-]?\d+)?|` +
		`(0|[1-9]\d*)(\.\d+)?[eE][+-]?\d+)` + numTail

	Integer = `[+-]?(0[bB]\d+|0[xX][\dA-Fa-f]+|0\d*|[1-9]\d*)` + numTail

	AnyChar = "."
)
