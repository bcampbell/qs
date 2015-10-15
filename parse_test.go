package qs

import (
	"github.com/blevesearch/bleve"
	"reflect"
	//"strings"
	"testing"
)

var NewIndexMapping = bleve.NewIndexMapping

func TestQuerySyntaxParserValid(t *testing.T) {
	//fivePointOh := 5.0
	//theTruth := true
	//theFalsehood := false
	tests := []struct {
		input   string
		result  bleve.Query
		mapping *bleve.IndexMapping
	}{
		{
			input:   "test",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("test"),
		},
		{
			input:   `"test phrase 1"`,
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("test phrase 1"),
		},
		{
			input:   "field:test",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("test").SetField("field"),
		},
		// - is allowed inside a term, just not the start
		{
			input:   "field:t-est",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("t-est").SetField("field"),
		},
		// + is allowed inside a term, just not the start
		{
			input:   "field:t+est",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("t+est").SetField("field"),
		},
		/*
			// > is allowed inside a term, just not the start
			{
				input:   "field:t>est",
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewMatchQuery("t>est").SetField("field"),
					},
					nil),
			},
			// < is allowed inside a term, just not the start
			{
				input:   "field:t<est",
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewMatchQuery("t<est").SetField("field"),
					},
					nil),
			},
			// = is allowed inside a term, just not the start
			{
				input:   "field:t=est",
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewMatchQuery("t=est").SetField("field"),
					},
					nil),
			},
		*/
		{
			input:   "+field1:test1",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("test1").SetField("field1"),
		},
		{
			input:   "-field2:test2",
			mapping: NewIndexMapping(),
			result: bleve.NewBooleanQuery(
				nil,
				nil,
				[]bleve.Query{
					bleve.NewMatchPhraseQuery("test2").SetField("field2"),
				}),
		},
		{
			input:   `field3:"test phrase 2"`,
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("test phrase 2").SetField("field3"),
		},
		{
			input:   `+field4:"test phrase 1"`,
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("test phrase 1").SetField("field4"),
		},
		{
			input:   `-field5:"test phrase 2"`,
			mapping: NewIndexMapping(),
			result: bleve.NewBooleanQuery(
				nil,
				nil,
				[]bleve.Query{
					bleve.NewMatchPhraseQuery("test phrase 2").SetField("field5"),
				}),
		},
		{
			input:   `+field6:test3 -field7:test4 field8:test5`,
			mapping: NewIndexMapping(),
			result: bleve.NewBooleanQuery(
				[]bleve.Query{
					bleve.NewMatchPhraseQuery("test3").SetField("field6"),
				},
				[]bleve.Query{
					bleve.NewMatchPhraseQuery("test5").SetField("field8"),
				},
				[]bleve.Query{
					bleve.NewMatchPhraseQuery("test4").SetField("field7"),
				}),
		},
		/*
			{
				input:   "test^3",
				mapping: NewIndexMapping(),
				result:  bleve.NewMatchPhraseQuery("test").SetBoost(3.0),
			},
			{
				input:   "test^3 other^6",
				mapping: NewIndexMapping(),
				result: bleve.NewDisjunctionQuery(
					[]bleve.Query{
						bleve.NewMatchPhraseQuery("test").SetBoost(3.0),
						bleve.NewMatchPhraseQuery("other").SetBoost(6.0),
					},
				),
			},
		*/
		{
			input:   "33",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("33"),
		},
		{
			input:   "field:33",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("33").SetField("field"),
		},
		{
			input:   "cat-dog",
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("cat-dog"),
		},
		/*
			        // TODO: MatchPhraseQuery doesn't handle fuzziness...
					{
						input:   "watex~",
						mapping: NewIndexMapping(),
						result:  bleve.NewMatchPhraseQuery("watex").SetFuzziness(1),
					},

				{
					input:   "watex~2",
					mapping: NewIndexMapping(),
					result: bleve.NewBooleanQuery(
						nil,
						[]bleve.Query{
							bleve.NewMatchQuery("watex").SetFuzziness(2),
						},
						nil),
				},
				{
					input:   "watex~ 2",
					mapping: NewIndexMapping(),
					result: bleve.NewBooleanQuery(
						nil,
						[]bleve.Query{
							bleve.NewMatchQuery("watex").SetFuzziness(1),
							bleve.NewMatchQuery("2"),
						},
						nil),
				},
				{
					input:   "field:watex~",
					mapping: NewIndexMapping(),
					result: bleve.NewBooleanQuery(
						nil,
						[]bleve.Query{
							bleve.NewMatchQuery("watex").SetFuzziness(1).SetField("field"),
						},
						nil),
				},
				{
					input:   "field:watex~2",
					mapping: NewIndexMapping(),
					result: bleve.NewBooleanQuery(
						nil,
						[]bleve.Query{
							bleve.NewMatchQuery("watex").SetFuzziness(2).SetField("field"),
						},
						nil),
				},
		*/
		{
			input:   `field:555c3bb06f7a127cda000005`,
			mapping: NewIndexMapping(),
			result:  bleve.NewMatchPhraseQuery("555c3bb06f7a127cda000005").SetField("field"),
		},
		/*
			{
				input:   `field:>5`,
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewNumericRangeInclusiveQuery(&fivePointOh, nil, &theFalsehood, nil).SetField("field"),
					},
					nil),
			},
			{
				input:   `field:>=5`,
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewNumericRangeInclusiveQuery(&fivePointOh, nil, &theTruth, nil).SetField("field"),
					},
					nil),
			},
			{
				input:   `field:<5`,
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewNumericRangeInclusiveQuery(nil, &fivePointOh, nil, &theFalsehood).SetField("field"),
					},
					nil),
			},
			{
				input:   `field:<=5`,
				mapping: NewIndexMapping(),
				result: bleve.NewBooleanQuery(
					nil,
					[]bleve.Query{
						bleve.NewNumericRangeInclusiveQuery(nil, &fivePointOh, nil, &theTruth).SetField("field"),
					},
					nil),
			},
		*/
		{
			input:   `grapefruit AND lemon`,
			mapping: NewIndexMapping(),
			result: bleve.NewConjunctionQuery([]bleve.Query{
				bleve.NewMatchPhraseQuery("grapefruit"),
				bleve.NewMatchPhraseQuery("lemon"),
			}),
		},
		{
			input:   `grapefruit OR lemon`,
			mapping: NewIndexMapping(),
			result: bleve.NewDisjunctionQuery([]bleve.Query{
				bleve.NewMatchPhraseQuery("grapefruit"),
				bleve.NewMatchPhraseQuery("lemon"),
			}),
		},
		{
			// default operator is OR
			input:   `grapefruit lemon`,
			mapping: NewIndexMapping(),
			result: bleve.NewBooleanQuery(
				nil,
				[]bleve.Query{
					bleve.NewMatchPhraseQuery("grapefruit"),
					bleve.NewMatchPhraseQuery("lemon"),
				},
				nil,
			),
		},
		{
			input:   `grapefruit AND NOT lemon`,
			mapping: NewIndexMapping(),
			result: bleve.NewConjunctionQuery([]bleve.Query{
				bleve.NewMatchPhraseQuery("grapefruit"),
				bleve.NewBooleanQuery(nil, nil, []bleve.Query{bleve.NewMatchPhraseQuery("lemon")}),
			}),
		},
		{
			input:   `field:(grapefruit AND lemon)`,
			mapping: NewIndexMapping(),
			result: bleve.NewConjunctionQuery([]bleve.Query{
				bleve.NewMatchPhraseQuery("grapefruit").SetField("field"),
				bleve.NewMatchPhraseQuery("lemon").SetField("field"),
			}),
		},
		{
			input:   `-field:(grapefruit AND lemon)`,
			mapping: NewIndexMapping(),
			result: bleve.NewBooleanQuery(nil, nil, []bleve.Query{
				bleve.NewConjunctionQuery([]bleve.Query{
					bleve.NewMatchPhraseQuery("grapefruit").SetField("field"),
					bleve.NewMatchPhraseQuery("lemon").SetField("field"),
				}),
			}),
		},
	}

	for _, test := range tests {

		q, err := Parse(test.input)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(q, test.result) {
			t.Errorf("Expected %#v, got %#v: for `%s`", test.result, q, test.input)
			//			t.Errorf("Expected %#v, got %#v: for %s", test.result.(*booleanQuery).Should.(*disjunctionQuery).Disjuncts[0], q.(*booleanQuery).Should.(*disjunctionQuery).Disjuncts[0], test.input)
		}
	}
}

/*
func TestQuerySyntaxParserValid(t *testing.T) {
	tests := []string{
		// fields
		`title:"The Right Way" AND text:go`,
		`title:"Do it right" AND right`,
		// wildcard searches
		`te?t`,
		`test*`,
		`te*t`,
		// fuzzy searches
		`roam~`,
		`roam~0.8`,
		// proximity searches
		`"jakarta apache"~10`,
		// ranges
		`mod_date:[20020101 TO 20030101]`,
		`title:{Aida TO Carmen}`,
		// boosting
		`jakarta^4 apache`,
		`"jakarta apache"^4 "Apache Lucene"`,
		// boolean operators
		`"jakarta apache" jakarta`,
		`"jakarta apache" OR jakarta`,
		`"jakarta apache" AND "Apache Lucene"`,
		`+jakarta lucene`,
		`"jakarta apache" NOT "Apache Lucene"`,
		`NOT "jakarta apache"`,
		// grouping
		`(jakarta OR apache) AND website`,
		// field grouping
		`title:(+return +"pink panther")`,
		// escaping special characters
		`\(1\+1\)\:2`,
	}

	for _, test := range tests {
		_, err := Parse(test)
		if err != nil {
			t.Errorf("failed to parse `%s` (%s)", test, err)
		}
	}
}
*/

func TestQuerySyntaxParserInvalid(t *testing.T) {
	tests := []string{
		"^",
		"^5",
		"field:-text",
		"field:+text",
		"field:>text",
		"field:>=text",
		"field:<text",
		"field:<=text",
		"field:~text",
		"field:^text",
		"field::text",
		// TODO:wildcards not permitted as first char:
		//`?est`,
		//`*iblle`,
	}

	// turn on lexer debugging
	// debugLexer = true
	// logger = log.New(os.Stderr, "bleve", log.LstdFlags)

	for _, test := range tests {
		_, err := Parse(test)
		if err == nil {
			t.Errorf("expected error, got nil for `%s`", test)
		}
	}
}
