package main

const (
	ItemError ItemType = iota

	ItemEOF

	ItemMetaKey
	ItemMetaValue

	ItemOverview

	// keywords: headers
	ItemTitleLevel1
	ItemTitleLevel2
	ItemTitleLevel3
	ItemTitleLevel4
	ItemTitleLevel5
	ItemTitleLevel6

	ItemGroup
	ItemRequired
	ItemURI

	// Data structures section
	ItemDataStructures // Section Title
	ItemModel
	ItemPropertyName
	ItemPropertyType
	ItemPropertyDesc

	// keywords: HTTP methods
	ItemCONNECT
	ItemDEL
	ItemGET
	ItemHEAD
	ItemOPTIONS
	ItemPOST
	ItemPUT
	ItemTRACE

	// keywords: list
	ItemAttributes
	ItemBody
	ItemParameters
	ItemRelation
	ItemRequest
	ItemResponse
	ItemSchema
	ItemValues
)

// LexMetaKey scans the Meta Section for the key in a key-value pair.
func LexMetaKey(l *Lexer) StateFn {
	if l.Peek() == '#' {
		return LexSectionTitle
	}

	l.AcceptUntil(":\r\n")
	l.Emit(ItemMetaKey)

	if l.Accept(":") {
		l.Ignore()
		return LexMetaValue
	}

	return nil
}

// LexMetaValue scans the Meta Section for the value in a key-value pair.
func LexMetaValue(l *Lexer) StateFn {
	// Ignore leading WS
	l.AcceptRun(" ")
	l.Ignore()

	// Everything after initial WS to EOL is the value.
	l.AcceptUntil("\n\r")
	l.Emit(ItemMetaValue)

	// Collapse WS to start MetaKey.
	l.AcceptRun("\r\n\t ")
	l.Ignore()

	return LexMetaKey
}

// LexSectionTitle lexes the section title.
func LexSectionTitle(l *Lexer) StateFn {
	start := l.Pos()
	l.AcceptRun("#")
	diff := l.Pos() - start

	// consume WS
	l.AcceptRun(" ")
	l.Ignore()

	// consume to EOL
	l.AcceptUntil("\r\n")

	switch {
	case l.HasPrefix("Data Structures"):
		l.Emit(ItemDataStructures)
		l.AcceptClasses(Whitespace)
		l.Ignore()
		return LexModel
	}

	titles := []ItemType{
		ItemTitleLevel1,
		ItemTitleLevel2,
		ItemTitleLevel3,
		ItemTitleLevel4,
		ItemTitleLevel5,
		ItemTitleLevel6,
	}
	d := titles[diff-1]

	l.Emit(d)

	l.AcceptClasses(Whitespace)
	l.Ignore()

	return LexOverview
}

func LexModel(l *Lexer) StateFn {
	l.AcceptRun("#")

	// consume WS
	l.AcceptRun(" ")
	l.Ignore()

	// consume to EOL, WS or right-parenthesis.
	l.AcceptUntil("\r\n (")

	l.Emit(ItemModel)
	l.AcceptClasses(Whitespace)
	l.Ignore()

	return LexPropertyName
}

// LexOverview scans the for the Overview body.
func LexOverview(l *Lexer) StateFn {
	for {
		// consume to EOL
		l.AcceptUntil("\n")

		// consume WS
		l.AcceptRun("\r\n ")

		// peek if next line is a section title
		if l.Peek() == '#' {
			break
		}
	}
	l.Emit(ItemOverview)

	return LexSectionTitle
}

func LexPropertyName(l *Lexer) StateFn {
	// consume + and WS
	l.Accept("+")
	l.AcceptRun("\t ")
	l.Ignore()

	// capture everything until WS or :
	l.AcceptClasses(Letter, Number)
	r := l.Peek()
	if !(r == ':' || r == ' ') {
		l.Errorf("unexpected character 0x%v for property name", r)
		return nil
	}
	l.Emit(ItemPropertyName)

	// consume trailing whitespace
	l.AcceptRun(" \t")
	l.Ignore()

	return LexPropertyType
}

func LexPropertyType(l *Lexer) StateFn {
	// consume and ignore (
	l.Accept("(")
	l.Ignore()

	// capture letters
	l.AcceptClasses(Letter)
	r := l.Peek()
	if r == ',' || r == ')' {
		l.Emit(ItemPropertyType)
	} else {
		l.Errorf("unexpected character 0x%v for property type", l.Peek())
		return nil
	}

	// consume boundary and ignore WS
	l.Accept(",)")
	l.AcceptClasses(Whitespace)
	l.Ignore()

	r = l.Peek()
	if r == '#' {
		return LexModel
	} else if r == '-' {
		return LexPropertyDesc
	}

	return LexPropertyName
}

func LexPropertyRequired(l *Lexer) StateFn {
	l.Accept(", ")
	l.Ignore()
	return nil
}

func LexPropertyDesc(l *Lexer) StateFn {
	l.AcceptUntil("\r\n")

	l.Emit(ItemPropertyDesc)
	l.AcceptClasses(Whitespace)
	l.Ignore()

	return LexPropertyName
}

func Whitespace(ch rune) bool {
	if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		return true
	}

	return false
}

func Letter(ch rune) bool {
	if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' {
		return true
	}

	return false
}

func Number(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}

	return false
}
