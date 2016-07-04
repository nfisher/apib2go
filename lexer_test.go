package main_test

import "testing"

import . "github.com/nfisher/apib2go"

func Test_Lexer_Next_should_increment_through_characters_in_string(t *testing.T) {
	l := New("meta.apib", `ab`)

	for i, r := range []rune{'a', 'b', EOF} {
		ch := l.Next()
		if ch != r {
			t.Errorf("[%v] got ch = %v, want %v", i, ch, r)
		}

	}
}

func Test_Lexer_Errorf_should_emit_error_with_description(t *testing.T) {
	l := New("meta.apib", `func`)
	go func() {
		l.Errorf("Because")
	}()

	e := <-l.Items
	if e.Type != ItemError {
		t.Fatalf("got e.Type = %v, want %v", e.Type, ItemError)
	}

	if e.Value != "Because" {
		t.Fatalf("got e.Value = %v, want %v", e.Value, "Because")
	}
}

func Test_Lexer_Accept(t *testing.T) {
	l := New("meta.apib", `bo`)
	dataTable := [][]interface{}{
		{"zb", true, "should accept first character"},
		{"b", false, "should not accept second character"},
		{"bo", true, "should accept second character"},
		{"o", false, "should not accept EOF"},
	}

	for i, td := range dataTable {
		valid := td[0].(string)
		actual := l.Accept(valid)
		expected := td[1].(bool)
		if actual != expected {
			t.Errorf("[%v] got l.Accept(%v) = %v, want %v. %v", i, valid, actual, expected, td[2].(string))
		}
	}
}

func Test_Lexer_AcceptRun(t *testing.T) {
	// valid run runes, expected position, expectation desc.
	dataTable := [][]interface{}{
		{"012", 2, "should move position forward 2 runes"},
		{"0123456789.", 11, "should move position forward 11 runes"},
		{"0123456789.!eha", 18, "should move position forward to eof"},
		{".!eha", 0, "should not move position forward"},
	}

	for i, td := range dataTable {
		valid := td[0].(string)
		l := New("meta.apib", `123456789.0e!hahah`)
		l.AcceptRun(valid)
		actual := l.Pos()
		expected := td[1].(int)

		if actual != expected {
			desc := td[2].(string)
			t.Errorf("[%v] got l.AcceptRun(%v) = %v, want %v. %v", i, valid, actual, expected, desc)
		}
	}
}

func Test_Lexer_AcceptUntil(t *testing.T) {
	dataTable := [][]interface{}{
		{"\n\r", "12345\r\n", 5, "12345", "should stop at carriage return."},
		{"\r\n", "12345\n67890", 5, "12345", "should stop at new line."},
		{"\r\n", "12345", 5, "12345", "should stop at EOF."},
		{":", "FORMAT: 1A9", 6, "FORMAT", "should stop at colon."},
	}

	for i, td := range dataTable {
		stop := td[0].(string)
		l := New("meta.apib", td[1].(string))
		l.AcceptUntil(stop)
		pos := l.Pos()
		expPos := td[2].(int)
		if pos != expPos {
			t.Errorf("[%v] want l.Pos() = %v, got %v", i, expPos, pos)
		}

		var pants ItemType = 100
		go func() {
			l.Emit(pants)
		}()

		exp := td[3].(string)
		item := <-l.Items
		if item.Value != exp {
			t.Errorf("[%v] want item.Value = %v, go %v", i, exp, item.Value)
		}
	}
}

func Test_Lexer_Emit(t *testing.T) {
	l := New("meta.apib", `func stuff()`)
	l.AcceptRun("func")
	var pants ItemType = 100
	go func() {
		l.Emit(pants)
	}()

	item := <-l.Items
	if item.Value != "func" {
		t.Errorf("got item.Value = %v, want %v", item.Value, "func")
	}

	if item.Type != pants {
		t.Errorf("got item.Type = %v, want %v", item.Type, pants)
	}
}
