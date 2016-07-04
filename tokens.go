package main

const (
	ItemError ItemType = iota

	ItemEOF

	ItemMetaKey
	ItemMetaValue

	ItemApiName
	ItemOverview

	// keywords: headers
	ItemSectionTitle
	ItemDataStructures
	ItemGroup
	ItemRequired
	ItemURI
	ItemStructureName

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
	ItemHeaders
	ItemModel
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
		return LexApiName
	}

	l.AcceptUntil(":\r\n")
	l.Emit(ItemMetaKey)

	if l.Accept(":") {
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

// LexApiName
func LexApiName(l *Lexer) StateFn {
	start := l.Pos()
	l.AcceptRun("#")
	diff := l.Pos() - start

	// check number of hashes is what we want
	if diff != 1 {
		return l.Errorf("API name should have 1 hash but was %v.", diff)
	}

	// consume WS
	l.AcceptRun(" ")
	l.Ignore()

	// consume to EOL
	l.AcceptUntil("\r\n")
	l.Emit(ItemApiName)

	return LexOverview
}

// LexSectionTitle lexes the section title.
func LexSectionTitle(l *Lexer) StateFn {
	l.AcceptRun("#")

	// consume WS
	l.AcceptRun(" ")
	l.Ignore()

	// consume to EOL
	l.AcceptUntil("\r\n")

	switch {
	case l.HasPrefix("Data Structures"):
		l.Emit(ItemDataStructures)
		return nil
	}

	l.Emit(ItemSectionTitle)
	return nil
}

func LexDataStructure(l *Lexer) StateFn {
	l.AcceptRun("#")

	// consume WS
	l.AcceptRun(" ")
	l.Ignore()

	// consume to EOL, WS or right-parenthesis.
	l.AcceptUntil("\r\n (")

	l.Emit(ItemStructureName)
	return nil
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
