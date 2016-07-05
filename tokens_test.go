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

func Test_Lex_functions(t *testing.T) {
	t.Parallel()

	// doc, pos, item, value
	dataTable := [][]interface{}{
		{"FORMAT:", 7, ItemMetaKey, "FORMAT", StateFn(LexMetaKey)},                                           //0
		{"FORMAT 2B:", 10, ItemMetaKey, "FORMAT 2B", StateFn(LexMetaKey)},                                    //1
		{" 1A9\n\nAUTHOR:", 6, ItemMetaValue, "1A9", StateFn(LexMetaValue)},                                  //2
		{" nf@jbx\n\n# API", 9, ItemMetaValue, "nf@jbx", StateFn(LexMetaValue)},                              //3
		{"# Level1\nab", 9, ItemTitleLevel1, "Level1", StateFn(LexSectionTitle)},                             //4
		{"# Level1\n\na", 10, ItemTitleLevel1, "Level1", StateFn(LexSectionTitle)},                           //5
		{"## Level2\na", 10, ItemTitleLevel2, "Level2", StateFn(LexSectionTitle)},                            //6
		{"### Level3\na", 11, ItemTitleLevel3, "Level3", StateFn(LexSectionTitle)},                           //7
		{"#### Level4\na", 12, ItemTitleLevel4, "Level4", StateFn(LexSectionTitle)},                          //8
		{"##### Level5\na", 13, ItemTitleLevel5, "Level5", StateFn(LexSectionTitle)},                         //9
		{"###### Level6\na", 14, ItemTitleLevel6, "Level6", StateFn(LexSectionTitle)},                        //10
		{"Overview\n#", 9, ItemOverview, "Overview\n", StateFn(LexOverview)},                                 //11
		{"Overview #\n#", 11, ItemOverview, "Overview #\n", StateFn(LexOverview)},                            //12
		{"## Data Structures\n\n", 20, ItemDataStructures, "Data Structures", StateFn(LexSectionTitle)},      //13
		{"### Author\n", 11, ItemModel, "Author", StateFn(LexModel)},                                         //14
		{"### Author (object)\n", 11, ItemModel, "Author", StateFn(LexModel)},                                //15
		{"### Author(object)\n", 10, ItemModel, "Author", StateFn(LexModel)},                                 //16
		{"+ email:", 7, ItemPropertyName, "email", StateFn(LexPropertyName)},                                 //17
		{"+ email ", 8, ItemPropertyName, "email", StateFn(LexPropertyName)},                                 //18
		{"+ email* ", 7, ItemError, "unexpected character 0x42 for property name", StateFn(LexPropertyName)}, //19
		{"(number\n", 7, ItemError, "unexpected character 0x10 for property type", StateFn(LexPropertyType)}, //20
		{"(number) ", 9, ItemPropertyType, "number", StateFn(LexPropertyType)},                               //21
		{"(number,required)", 8, ItemPropertyType, "number", StateFn(LexPropertyType)},                       //22
		{"(number) ", 9, ItemPropertyType, "number", StateFn(LexPropertyType)},                               //23
		{"(number)\n", 9, ItemPropertyType, "number", StateFn(LexPropertyType)},                              //24
		{"- desc\n", 7, ItemPropertyDesc, "- desc", StateFn(LexPropertyDesc)},                                //25
		{"- desc\r\n", 8, ItemPropertyDesc, "- desc", StateFn(LexPropertyDesc)},                              //26
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
			t.Errorf("[%v] l.Pos() = %v, want %v", i, l.Pos(), expPos)
		}

		expected := Item{td[2].(ItemType), td[3].(string)}
		if item != expected {
			t.Errorf("[%v] item.Type = %v, want %v", i, item, expected)
		}
	}
}

func Test_simple_document(t *testing.T) {
	t.Parallel()

	var doc = `Version: 1A9

# Simple API
Overview

## Data Structures

### Dimension
+ radius (number)
+ length (number)

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
		Item{ItemTitleLevel1, "Simple API"},
		Item{ItemOverview, "Overview\n\n"},
		Item{ItemDataStructures, "Data Structures"},
		Item{ItemModel, "Dimension"},
		Item{ItemPropertyName, "radius"},
		Item{ItemPropertyType, "number"},
		Item{ItemPropertyName, "length"},
		Item{ItemPropertyType, "number"},
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
