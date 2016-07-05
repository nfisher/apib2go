package main

const (
	ItemError ItemType = iota

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
	ItemPropertyArrayType
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

	l.AcceptUntil(":")
	if l.Peek() == ':' {
		l.Emit(ItemMetaKey)
		l.Accept(":")

		return LexMetaValue
	}

	l.Errorf("not valid meta key.")
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

	if l.Peek() == EOF {
		return nil
	}

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

	r := l.Peek()
	if r == '#' {
		return LexSectionTitle
	}

	return LexOverview
}

// LexOverview scans the for the Overview body.
func LexOverview(l *Lexer) StateFn {
	for {
		// consume to EOL
		l.AcceptUntil("\n")

		// consume WS
		l.AcceptRun("\r\n ")

		// peek if next line is a section title
		r := l.Peek()
		if r == '#' || r == EOF {
			break
		}
	}
	l.Emit(ItemOverview)

	return LexSectionTitle
}

// LexModel scans for a models name.
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

// LexPropertyName scans for a properties name.
func LexPropertyName(l *Lexer) StateFn {
	// consume + and WS
	l.Accept("+")
	l.AcceptRun("\t ")
	l.Ignore()

	// capture everything until WS or :
	l.AcceptClasses(Letter, Number)
	r := l.Peek()
	if !(r == ':' || r == ' ') {
		l.Errorf("unexpected character `%v`:0x%v for property name", string(r), r)
		return nil
	}
	l.Emit(ItemPropertyName)

	if l.Accept(":") {
		return LexPropertyExample
	}

	// consume trailing whitespace
	l.AcceptRun(" \t")
	l.Ignore()

	return LexPropertyType
}

// LexPropertyExample scans for a property example.
func LexPropertyExample(l *Lexer) StateFn {
	l.AcceptUntil("(")
	l.Ignore()
	return LexPropertyType
}

// LexPropertyType scans for a type.
func LexPropertyType(l *Lexer) StateFn {
	// consume and ignore (
	l.Accept("(")
	l.Ignore()

	// capture letters
	l.AcceptClasses(Letter)
	r := l.Peek()
	if r == ',' || r == ')' {
		l.Emit(ItemPropertyType)
	} else if l.Accept("[") {
		l.Ignore()
		l.AcceptClasses(Letter)
		if l.Peek() == ']' {
			l.Emit(ItemPropertyArrayType)
			// consume ]
			l.Next()
		} else {
			l.Errorf("missing closing brace in array type")
			return nil
		}
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
	} else if r == EOF {
		return nil
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
