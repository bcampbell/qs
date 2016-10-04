//  Copyright (c) 2014 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// NOTE:
// This file is a copy of bleve/search/query/query_string_parser_test.go
// with minimal changes to make it easier to track

package qs

import (
	"reflect"
	"testing"
	"time"

	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"
)

// "import" some stuff to make it easier to track bleve/search/query/query_string_parser_test.go
// although, in practice... so many other changes that it's hardly worth it...
type Query query.Query

var (
	NewMatchQuery                 = query.NewMatchQuery
	NewMatchPhraseQuery           = query.NewMatchPhraseQuery
	NewNumericRangeInclusiveQuery = query.NewNumericRangeInclusiveQuery
	NewDateRangeInclusiveQuery    = query.NewDateRangeInclusiveQuery
	NewRegexpQuery                = query.NewRegexpQuery
	NewWildcardQuery              = query.NewWildcardQuery
	NewFuzzyQuery                 = query.NewFuzzyQuery
)

func NewBooleanQuery(must []Query, should []Query, mustNot []Query) *query.BooleanQuery {
	must2 := make([]query.Query, len(must))
	should2 := make([]query.Query, len(should))
	mustNot2 := make([]query.Query, len(mustNot))

	for i := 0; i < len(must); i++ {
		must2[i] = must[i]
	}
	for i := 0; i < len(should); i++ {
		should2[i] = should[i]
	}
	for i := 0; i < len(mustNot); i++ {
		mustNot2[i] = mustNot[i]
	}
	return query.NewBooleanQuery(must2, should2, mustNot2)
}

func TestQuerySyntaxParserValid(t *testing.T) {
	fivePointOh := 5.0
	onePointOh := 1.0
	theTruth := true
	theFalsehood := false
	jan_01_2015 := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan_02_2015 := time.Date(2015, time.January, 2, 0, 0, 0, 0, time.UTC)
	mar_15_2015 := time.Date(2015, time.March, 15, 0, 0, 0, 0, time.UTC)
	mar_16_2015 := time.Date(2015, time.March, 16, 0, 0, 0, 0, time.UTC)

	/*theDate, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	*/
	tests := []struct {
		input   string
		result  Query
		mapping mapping.IndexMapping
	}{
		{
			input:   "test",
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("test"),
		},
		{
			input:   `"test phrase 1"`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("test phrase 1"),
		},
		{
			input:   "field:test",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("test")
				q.SetField("field")
				return q
			}(),
		},
		// - is allowed inside a term, just not the start
		{
			input:   "field:t-est",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("t-est")
				q.SetField("field")
				return q
			}(),
		},
		// + is allowed inside a term, just not the start
		{
			input:   "field:t+est",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("t+est")
				q.SetField("field")
				return q
			}(),
		},
		// > is allowed inside a term, just not the start
		{
			input:   "field:t>est",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("t>est")
				q.SetField("field")
				return q
			}(),
		},
		// < is allowed inside a term, just not the start
		{
			input:   "field:t<est",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("t<est")
				q.SetField("field")
				return q
			}(),
		},
		// = is allowed inside a term, just not the start
		{
			input:   "field:t=est",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("t=est")
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   "+field1:test1",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("test1")
				q.SetField("field1")
				return q
			}(),
		},
		{
			input:   "-field2:test2",
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				nil,
				[]Query{
					func() Query {
						q := NewMatchPhraseQuery("test2")
						q.SetField("field2")
						return q
					}(),
				}),
		},
		{
			input:   `field3:"test phrase 2"`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("test phrase 2")
				q.SetField("field3")
				return q
			}(),
		},
		{
			input:   `+field4:"test phrase 1"`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("test phrase 1")
				q.SetField("field4")
				return q
			}(),
		},
		{
			input:   `-field5:"test phrase 2"`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				nil,
				[]Query{
					func() Query {
						q := NewMatchPhraseQuery("test phrase 2")
						q.SetField("field5")
						return q
					}(),
				}),
		},
		{
			input:   `+field6:test3 -field7:test4 field8:test5`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				[]Query{
					func() Query {
						q := NewMatchPhraseQuery("test3")
						q.SetField("field6")
						return q
					}(),
				},
				[]Query{
					func() Query {
						q := NewMatchPhraseQuery("test5")
						q.SetField("field8")
						return q
					}(),
				},
				[]Query{
					func() Query {
						q := NewMatchPhraseQuery("test4")
						q.SetField("field7")
						return q
					}(),
				}),
		},
		{
			input:   "test^3",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("test")
				q.SetBoost(3.0)
				return q
			}(),
		},
		{
			input:   "test^3 other^6",
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewMatchPhraseQuery("test")
						q.SetBoost(3.0)
						return q
					}(),
					func() Query {
						q := NewMatchPhraseQuery("other")
						q.SetBoost(6.0)
						return q
					}(),
				},
				nil),
		},
		{
			input:   "33",
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("33"),
		},
		{
			input:   "field:33",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("33")
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   "cat-dog",
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("cat-dog"),
		},
		{
			input:   "watex~",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewFuzzyQuery("watex")
				q.SetFuzziness(1)
				return q
			}(),
		},
		{
			input:   "watex~2",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewFuzzyQuery("watex")
				q.SetFuzziness(2)
				return q
			}(),
		},

		{
			input:   "watex~ 2",
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewFuzzyQuery("watex")
						q.SetFuzziness(1)
						return q
					}(),
					NewMatchPhraseQuery("2"),
				},
				nil),
		},
		{
			input:   "field:watex~",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewFuzzyQuery("watex")
				q.SetFuzziness(1)
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   "field:watex~2",
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewFuzzyQuery("watex")
				q.SetFuzziness(2)
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   `field:555c3bb06f7a127cda000005`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("555c3bb06f7a127cda000005")
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   `field:>5`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewNumericRangeInclusiveQuery(&fivePointOh, nil, &theFalsehood, nil)
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   `field:>=5`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewNumericRangeInclusiveQuery(&fivePointOh, nil, &theTruth, nil)
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   `field:<5`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewNumericRangeInclusiveQuery(nil, &fivePointOh, nil, &theFalsehood)
				q.SetField("field")
				return q
			}(),
		},
		{
			input:   `field:<=5`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewNumericRangeInclusiveQuery(nil, &fivePointOh, nil, &theTruth)
				q.SetField("field")
				return q
			}(),
		},
		/* XYZZY - TODO: support use of bleve dateparser
		{
			input:   `field:>"2006-01-02T15:04:05Z"`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewDateRangeInclusiveQuery(theDate, time.Time{}, &theFalsehood, nil)
						q.SetField("field")
						return q
					}(),
				},
				nil),
		},
		{
			input:   `field:>="2006-01-02T15:04:05Z"`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewDateRangeInclusiveQuery(theDate, time.Time{}, &theTruth, nil)
						q.SetField("field")
						return q
					}(),
				},
				nil),
		},
		{
			input:   `field:<"2006-01-02T15:04:05Z"`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewDateRangeInclusiveQuery(time.Time{}, theDate, nil, &theFalsehood)
						q.SetField("field")
						return q
					}(),
				},
				nil),
		},
		{
			input:   `field:<="2006-01-02T15:04:05Z"`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewDateRangeInclusiveQuery(time.Time{}, theDate, nil, &theTruth)
						q.SetField("field")
						return q
					}(),
				},
				nil),
		},
		*/

		/* XYZZY - TODO: regexp support
		{
			input:   `/mar.*ty/`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					NewRegexpQuery("mar.*ty"),
				},
				nil),
		},
		{
			input:   `name:/mar.*ty/`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					func() Query {
						q := NewRegexpQuery("mar.*ty")
						q.SetField("name")
						return q
					}(),
				},
				nil),
		},
		*/
		{
			input:   `mart*`,
			mapping: mapping.NewIndexMapping(),
			result:  NewWildcardQuery("mart*"),
		},
		{
			input:   `name:mart*`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewWildcardQuery("mart*")
				q.SetField("name")
				return q
			}(),
		},

		// tests for escaping

		/* XYZZY
		// escape : as field delimeter
		{
			input:   `name\:marty`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("name:marty"),
		},
		// first colon delimiter, second escaped
		{
			input:   `name:marty\:couchbase`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery("marty:couchbase")
				q.SetField("name")
				return q
			}(),
		},
		// escape space, single arguemnt to match query
		{
			input:   `marty\ couchbase`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchQuery("marty couchbase"),
		},
		// escape leading plus, not a must clause
		{
			input:   `\+marty`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("+marty"),
		},
		// escape leading minus, not a must not clause
		{
			input:   `\-marty`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery("-marty"),
		},
		// escape quote inside of phrase
		{
			input:   `"what does \"quote\" mean"`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery(`what does "quote" mean`),
		},
		// escaping an unsupported character retains backslash
		{
			input:   `can\ i\ escap\e`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchQuery(`can i escap\e`),
		},
		*/
		// leading spaces
		{
			input:   `   what`,
			mapping: mapping.NewIndexMapping(),
			result:  NewMatchPhraseQuery(`what`),
		},
		// no boost value defaults to 1
		{
			input:   `term^`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := NewMatchPhraseQuery(`term`)
				q.SetBoost(1.0)
				return q
			}(),
		},
		/* XYZZY
		// weird lexer cases, something that starts like a number
		// but contains escape and ends up as string
		{
			input:   `3.0\:`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					NewMatchQuery(`3.0:`),
				},
				nil),
		},
		{
			input:   `3.0\a`,
			mapping: mapping.NewIndexMapping(),
			result: NewBooleanQuery(
				nil,
				[]Query{
					NewMatchQuery(`3.0\a`),
				},
				nil),
		},
		*/

		/*
		* Extra stuff, above what querystringquery supports
		*
		 */
		{
			input:   `grapefruit AND lemon`,
			mapping: mapping.NewIndexMapping(),
			result: query.NewConjunctionQuery([]query.Query{
				query.NewMatchPhraseQuery("grapefruit"),
				query.NewMatchPhraseQuery("lemon"),
			}),
		},
		{
			input:   `grapefruit OR lemon`,
			mapping: mapping.NewIndexMapping(),
			result: query.NewDisjunctionQuery([]query.Query{
				query.NewMatchPhraseQuery("grapefruit"),
				query.NewMatchPhraseQuery("lemon"),
			}),
		},
		{
			// default operator is OR
			input:   `grapefruit lemon`,
			mapping: mapping.NewIndexMapping(),
			result: query.NewBooleanQuery(
				nil,
				[]query.Query{
					query.NewMatchPhraseQuery("grapefruit"),
					query.NewMatchPhraseQuery("lemon"),
				},
				nil,
			),
		},
		{
			input:   `grapefruit AND NOT lemon`,
			mapping: mapping.NewIndexMapping(),
			result: query.NewConjunctionQuery([]query.Query{
				query.NewMatchPhraseQuery("grapefruit"),
				query.NewBooleanQuery(nil, nil, []query.Query{query.NewMatchPhraseQuery("lemon")}),
			}),
		},
		{
			input:   `field:(grapefruit AND lemon)`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				t1 := NewMatchPhraseQuery(`grapefruit`)
				t1.SetField("field")
				t2 := NewMatchPhraseQuery(`lemon`)
				t2.SetField("field")
				q := query.NewConjunctionQuery([]query.Query{t1, t2})
				return q
			}(),
		},
		{
			input:   `-field:(grapefruit AND lemon)`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				t1 := NewMatchPhraseQuery(`grapefruit`)
				t1.SetField("field")
				t2 := NewMatchPhraseQuery(`lemon`)
				t2.SetField("field")
				andq := query.NewConjunctionQuery([]query.Query{t1, t2})
				q := query.NewBooleanQuery(nil, nil, []query.Query{andq})
				return q
			}(),
		},
		{
			input:   `shoesize:[1 TO 5]`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewNumericRangeInclusiveQuery(&onePointOh, &fivePointOh, &theTruth, &theTruth)
				q.SetField("shoesize")
				return q
			}(),
		},
		{
			input:   `shoesize:{1 TO 5}`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewNumericRangeInclusiveQuery(&onePointOh, &fivePointOh, &theFalsehood, &theFalsehood)
				q.SetField("shoesize")
				return q
			}(),
		},
		{
			input:   `shoesize:[1 TO 5}`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewNumericRangeInclusiveQuery(&onePointOh, &fivePointOh, &theTruth, &theFalsehood)
				q.SetField("shoesize")
				return q
			}(),
		},
		{
			input:   `shoesize:{1 TO 5]`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewNumericRangeInclusiveQuery(&onePointOh, &fivePointOh, &theFalsehood, &theTruth)
				q.SetField("shoesize")
				return q
			}(),
		},
		{
			input:   `shoesize:[ TO 5]`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewNumericRangeInclusiveQuery(nil, &fivePointOh, nil, &theTruth)
				q.SetField("shoesize")
				return q
			}(),
		},
		{
			input:   `shoesize:[1 TO ]`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewNumericRangeInclusiveQuery(&onePointOh, nil, &theTruth, nil)
				q.SetField("shoesize")
				return q
			}(),
		},
		// date ranges (note that endpoints and inclusivity might be modified by the parser)
		{
			input:   `when:[2015-01-01 TO 2015-03-15]`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewDateRangeInclusiveQuery(jan_01_2015, mar_16_2015, &theTruth, &theFalsehood)
				q.SetField("when")
				return q
			}(),
		},
		{
			input:   `when:{2015-01-01 TO 2015-03-15]`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewDateRangeInclusiveQuery(jan_02_2015, mar_16_2015, &theTruth, &theFalsehood)
				q.SetField("when")
				return q
			}(),
		},
		{
			input:   `when:[2015-01-01 TO 2015-03-15}`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewDateRangeInclusiveQuery(jan_01_2015, mar_15_2015, &theTruth, &theFalsehood)
				q.SetField("when")
				return q
			}(),
		},
		{
			input:   `when:>2015-03-15`,
			mapping: mapping.NewIndexMapping(),
			result: func() Query {
				q := query.NewDateRangeInclusiveQuery(mar_16_2015, time.Time{}, &theTruth, nil)
				q.SetField("when")
				return q
			}(),
		},
		// Wildcards
		{
			input:   `foo*`,
			mapping: mapping.NewIndexMapping(),
			result:  query.NewWildcardQuery(`foo*`),
		},
		{
			input:   `f?rt`,
			mapping: mapping.NewIndexMapping(),
			result:  query.NewWildcardQuery(`f?rt`),
		},
	}

	for _, test := range tests {

		q, err := Parse(test.input)
		if err != nil {
			t.Fatalf("`%s`: %s", test.input, err)
		}
		if !reflect.DeepEqual(q, test.result) {
			t.Errorf("Expected %#v, got %#v: for %s", test.result, q, test.input)
		}
	}
}

func TestQuerySyntaxParserInvalid(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"^"},
		{"^5"},
		{"field:-text"},
		{"field:+text"},
		{"field:>text"},
		{"field:>=text"},
		{"field:<text"},
		{"field:<=text"},
		{"field:~text"},
		{"field:^text"},
		{"field::text"},
		{`"this is the time`},
		{`cat^3\:`},
		{`cat^3\0`},
		{`cat~3\:`},
		{`cat~3\0`},
	}

	for _, test := range tests {
		_, err := Parse(test.input)
		if err == nil {
			t.Errorf("expected error, got nil for `%s`", test.input)
		}
	}
}

func BenchmarkLexer(b *testing.B) {

	for n := 0; n < b.N; n++ {

		tokens := lex(`+field4:"test phrase 1"`)
		for i, _ := range tokens {
			tokens[i].val = ""
			tokens[i].pos = 0
		}
	}

}
