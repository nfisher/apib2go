package main_test

import (
	"testing"

	. "github.com/nfisher/apib2go"
)

func Test_rune_classifiers(t *testing.T) {
	dataTable := [][]interface{}{
		// rune, whitespace, letter, number
		{'A', false, true, false},
		{'Z', false, true, false},
		{'a', false, true, false},
		{'z', false, true, false},
		{'0', false, false, true},
		{'9', false, false, true},
		{' ', true, false, false},
		{'\t', true, false, false},
		{'\n', true, false, false},
		{'\r', true, false, false},
	}

	for i, td := range dataTable {
		input := td[0].(rune)
		actual := Whitespace(input)
		expected := td[1].(bool)

		if actual != expected {
			t.Errorf("[%v] want Whitespace(0x%X) = %v, want %v", i, input, actual, expected)
		}

		actual = Letter(input)
		expected = td[2].(bool)

		if actual != expected {
			t.Errorf("[%v] Letter(0x%X) = %v, want %v", i, input, actual, expected)
		}

		actual = Number(input)
		expected = td[3].(bool)
		if actual != expected {
			t.Errorf("[%v] Number(0x%X) = %v, want %v", i, input, actual, expected)
		}
	}
}

func lexItem(doc string, fn StateFn) (item Item, pos int) {
	l := New("meta.apib", doc)
	go func() {
		fn(l)
	}()

	return <-l.Items, l.Pos()
}

func Test_LexMetaKey(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"FORMAT:", 7, ItemMetaKey, "FORMAT"},
		{"FORMAT\n", 7, ItemError, "not valid meta key."},
		{"FORMAT 2B:", 10, ItemMetaKey, "FORMAT 2B"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexMetaKey)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}
}

func Test_LexMetaValue(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{" 1A9\n\nAUTHOR:", 6, ItemMetaValue, "1A9"},
		{" nf@jbx\n\n# API", 9, ItemMetaValue, "nf@jbx"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexMetaValue)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}
}

func Test_LexSectionTitle(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"# Level1\nab", 9, ItemTitleLevel1, "Level1"},
		{"# Level1\n\na", 10, ItemTitleLevel1, "Level1"},
		{"## Level2\na", 10, ItemTitleLevel2, "Level2"},
		{"### Level3\na", 11, ItemTitleLevel3, "Level3"},
		{"#### Level4\na", 12, ItemTitleLevel4, "Level4"},
		{"##### Level5\na", 13, ItemTitleLevel5, "Level5"},
		{"###### Level6\na", 14, ItemTitleLevel6, "Level6"},
		{"## Data Structures\n\n", 20, ItemDataStructures, "Data Structures"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexSectionTitle)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}
}

func Test_LexOverview(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"Overview\n#", 9, ItemOverview, "Overview\n"},
		{"Overview #\n#", 11, ItemOverview, "Overview #\n"},
		{"Overview", 8, ItemOverview, "Overview"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexOverview)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}

}

func Test_LexModel(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"### Author\n", 11, ItemModel, "Author"},
		{"### Author (object)\n", 11, ItemModel, "Author"},
		{"### Author(object)\n", 10, ItemModel, "Author"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexModel)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}
}

func Test_LexPropertyName(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"+ email:", 8, ItemPropertyName, "email"},
		{"+ email ", 8, ItemPropertyName, "email"},
		{"+ email* ", 7, ItemError, "unexpected character `*`:0x42 for property name"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexPropertyName)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v,\nwant %v", i, item, expected)
		}
	}
}

func Test_LexPropertyType(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"(number\n", 7, ItemError, "unexpected character 0x10 for property type"},
		{"(number) ", 9, ItemPropertyType, "number"},
		{"(number,required)", 8, ItemPropertyType, "number"},
		{"(number) ", 9, ItemPropertyType, "number"},
		{"(number)\n", 9, ItemPropertyType, "number"},
		{"(array[number])\n", 16, ItemPropertyArrayType, "number"},
		{"(array[number)\n", 13, ItemError, "missing closing brace in array type"},
		{"(arraynumber])\n", 12, ItemError, "unexpected character 0x93 for property type"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexPropertyType)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v for %v", i, pos, expPos, item)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}
}

func Test_LexPropertyDesc(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"- desc\n", 7, ItemPropertyDesc, "- desc"},
		{"- desc\r\n", 8, ItemPropertyDesc, "- desc"},
	}

	for i, td := range dataTable {
		item, pos := lexItem(td[0].(string), LexPropertyDesc)

		expPos := td[1].(int)
		if expPos != pos {
			t.Errorf("[%v] pos = %v, want %v", i, pos, expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item = %v, want %v", i, item, expected)
		}
	}
}

func Test_datastructure_apib_document(t *testing.T) {
	t.Parallel()
	var doc = `Version: 1A9

# DS API

## Data Structures

### Dimension
+ radius: 123 (number)
+ length (number) - This is a comment.

### Produce

+ colour (string) - What colour is it?
+ dimensions (Dimension)
+ fruit (boolean) - Is it fruit?`

	l := New("meta.apib", doc)

	go func() {
		l.Run()
	}()

	expected := []Item{
		Item{ItemMetaKey, "Version"},
		Item{ItemMetaValue, "1A9"},
		Item{ItemTitleLevel1, "DS API"},
		Item{ItemDataStructures, "Data Structures"},
		Item{ItemModel, "Dimension"},
		Item{ItemPropertyName, "radius"},
		Item{ItemPropertyType, "number"},
		Item{ItemPropertyName, "length"},
		Item{ItemPropertyType, "number"},
		Item{ItemPropertyDesc, "- This is a comment."},
		Item{ItemModel, "Produce"},
		Item{ItemPropertyName, "colour"},
		Item{ItemPropertyType, "string"},
		Item{ItemPropertyDesc, "- What colour is it?"},
		Item{ItemPropertyName, "dimensions"},
		Item{ItemPropertyType, "Dimension"},
		Item{ItemPropertyName, "fruit"},
		Item{ItemPropertyType, "boolean"},
		Item{ItemPropertyDesc, "- Is it fruit?"},
	}

	for i, ex := range expected {
		item := <-l.Items
		if item != ex {
			t.Errorf("[%v] item = %v, want %v: %s", i, item, ex, string(l.Peek()))
		}
	}
}

func Test_simple_api_document(t *testing.T) {
	t.Parallel()

	var doc = `Version: 1A9

# Simple API
Overview

# Group Health Check

## Ping [/ping]

### Ping-Pong [GET]

+ Request pong

    + Headers

            Accept: text/plain

+ Response 200 (text/plain; charset=utf-8)

        pong`

	req := `+ Request pong

    + Headers

            Accept: text/plain

+ Response 200 (text/plain; charset=utf-8)

        pong`

	l := New("meta.apib", doc)

	go func() {
		l.Run()
	}()

	expected := []Item{
		Item{ItemMetaKey, "Version"},
		Item{ItemMetaValue, "1A9"},
		Item{ItemTitleLevel1, "Simple API"},
		Item{ItemOverview, "Overview\n\n"},
		Item{ItemTitleLevel1, "Group Health Check"},
		Item{ItemTitleLevel2, "Ping [/ping]"},
		Item{ItemTitleLevel3, "Ping-Pong [GET]"},
		Item{ItemOverview, req},
	}

	for i, ex := range expected {
		item := <-l.Items
		if item != ex {
			t.Errorf("[%v] item = %v, want %v: %s", i, item, ex, string(l.Peek()))
		}
	}
}
