package main_test

import (
	"testing"

	. "github.com/nfisher/apib2go"
)

func Test_rune_classifiers(t *testing.T) {
	dataTable := [][]interface{}{
		// rune, whitespace, letter
		{'A', false, true},
		{'Z', false, true},
		{'a', false, true},
		{'z', false, true},
		{'0', false, false},
		{'9', false, false},
		{' ', true, false},
		{'\t', true, false},
		{'\n', true, false},
		{'\r', true, false},
	}

	for i, td := range dataTable {
		input := td[0].(rune)
		actual := Whitespace(input)
		expected := td[1].(bool)

		if actual != expected {
			t.Errorf("[%v] want Whitespace(0x%X) = %v, but was %v", i, input, actual, expected)
		}

		actual = Letter(input)
		expected = td[2].(bool)

		if actual != expected {
			t.Errorf("[%v] want Letter(0x%X) = %v, but was %v", i, input, actual, expected)
		}
	}
}

func Test_Lex_functions(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"FORMAT: 1A9", 7, ItemMetaKey, "FORMAT", StateFn(LexMetaKey)},
		{"FORMAT 2B: 1A9", 10, ItemMetaKey, "FORMAT 2B", StateFn(LexMetaKey)},
		{" 1A9\n\nAUTHOR:", 6, ItemMetaValue, "1A9", StateFn(LexMetaValue)},
		{"# Policy API\n\n", 12, ItemApiName, "Policy API", StateFn(LexApiName)},
		{"## Too Many\n\n", 2, ItemError, "API name should have 1 hash but was 2.", StateFn(LexApiName)},
		{"Overview\n#", 9, ItemOverview, "Overview\n", StateFn(LexOverview)},
		{"Overview #\n#", 11, ItemOverview, "Overview #\n", StateFn(LexOverview)},
		{"## Data Structures\n", 18, ItemDataStructures, "Data Structures", StateFn(LexSectionTitle)},
		{"### Author\n", 10, ItemStructureName, "Author", StateFn(LexDataStructure)},
		{"### Author (object)\n", 10, ItemStructureName, "Author", StateFn(LexDataStructure)},
		{"### Author(object)\n", 10, ItemStructureName, "Author", StateFn(LexDataStructure)},
	}

	for i, td := range dataTable {
		l := New("meta.apib", td[0].(string))
		go func() {
			fn, ok := td[4].(StateFn)
			if ok {
				fn(l)
			} else {
				t.Errorf("[%v] can't cast value to StateFn.", i)
			}
		}()

		item := <-l.Items

		expPos := td[1].(int)
		if expPos != l.Pos() {
			t.Errorf("[%v] want l.Pos() = %v, got %v", i, expPos, l.Pos())
		}

		expected := td[2].(ItemType)
		if item.Type != expected {
			t.Errorf("[%v] want item.Type = %v, got %v", i, expected, item.Type)
		}

		expVal := td[3].(string)
		if item.Value != expVal {
			t.Errorf("[%v] want item.Value = `%v`, got `%v`", i, expVal, item.Value)
		}
	}
}

/*
func Test_simple_document(t *testing.T) {
	t.Parallel()

	var doc = `Version: 1A9

# Simple API
Simple API with only Data Structures defined.

## Data Structures

### Dimension (object)
+ radius (number)
+ length (number)

### Produce (object)

+ colour (string) - What colour is it?
+ dimesions (Dimension)
+ fruit (boolean) - Is it fruit?`

	if false {
		t.Fatal(doc)
	}
}
*/
