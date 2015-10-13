package qs

import (
	//	"fmt"
	"reflect"
	//	"strings"
	"testing"
)

/*
func ExampleLex() {

	queries := []string{
		``,
		`wibble`,
		`wibble pibble`,
		`"a phrase"`,
		`tag: citrus`,
		`tag:(lemon mango)`,
		`-tag:("lemon mango")`,
		`tag:(citrus -banana)`,
		`date:[2014-01-01 TO 2014-01-02]`,
		`1234wibble 12`,
	}

	for _, q := range queries {
		toks := lex(q)

		out := []string{}
		for _, tok := range toks {
			out = append(out, tok.typ.String())
		}
		fmt.Printf("%s %s\n", q, strings.Join(out, ","))
	}

	// Output:
	//

}
*/

func TestLexer(t *testing.T) {
	/*
		tokName := func(tok int) string {
			// Cheesy and britle - must be a better way...
			i := tok - yyPrivate - 1
			if i >= 0 && i < len(yyToknames) {
				return yyToknames[i]
			} else {
				return fmt.Sprintf("tok-%d", tok)
			}
		}
	*/
	data := []struct {
		input    string
		expected []tokType
	}{
		{`wibble 123 : < > { } [ ] - + ( ) TO AND OR NOT "this is a phrase"`, []tokType{tLITERAL, tLITERAL, tCOLON, tLESS, tGREATER, tLBRACE, tRBRACE, tLSQUARE, tRSQUARE, tMINUS, tPLUS, tLPAREN, tRPAREN, tTO, tAND, tOR, tNOT, tQUOTED, tEOF}},
		{`1.234`, []tokType{tLITERAL, tEOF}},
		{`     blah    `, []tokType{tLITERAL, tEOF}},
		{`foo:[ 1 TO 5 ]`, []tokType{tLITERAL, tCOLON, tLSQUARE, tLITERAL, tTO, tLITERAL, tRSQUARE, tEOF}},
		{`foo:[1 TO 5]`, []tokType{tLITERAL, tCOLON, tLSQUARE, tLITERAL, tTO, tLITERAL, tRSQUARE, tEOF}},
		{`t<est`, []tokType{tLITERAL, tEOF}},
		{`t>est`, []tokType{tLITERAL, tEOF}},
		{`t=est`, []tokType{tLITERAL, tEOF}},
		{`t+est`, []tokType{tLITERAL, tEOF}},
		{`TOe`, []tokType{tLITERAL, tEOF}},
		{`NOTe`, []tokType{tLITERAL, tEOF}},
		{`sAND`, []tokType{tLITERAL, tEOF}},
		{`g AND t`, []tokType{tLITERAL, tAND, tLITERAL, tEOF}},
		{`shoesize:>10`, []tokType{tLITERAL, tCOLON, tGREATER, tLITERAL, tEOF}},
		{`123four`, []tokType{tLITERAL, tEOF}},
		{`four321`, []tokType{tLITERAL, tEOF}},
		{`123 four`, []tokType{tLITERAL, tLITERAL, tEOF}},
		{`wibble~`, []tokType{tLITERAL, tFUZZY, tEOF}},
		{`wibble~0.1`, []tokType{tLITERAL, tFUZZY, tEOF}},
		{`wibble^5`, []tokType{tLITERAL, tBOOST, tEOF}},
	}

	for _, dat := range data {

		toks := lex(dat.input)
		got := make([]tokType, len(toks))
		for i, tok := range toks {
			got[i] = tok.typ
		}

		if !reflect.DeepEqual(dat.expected, got) {
			t.Errorf("Expected %#v, got %#v: for `%s`", dat.expected, got, dat.input)
		}
	}

}
