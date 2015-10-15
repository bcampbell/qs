// bleve_queryparser - parse and dump a query string to json
//
// The output is a little terse, so easier to read if you pipe it
// through a json prettifier first, eg:
//
//   $ ./bleve_queryparser 'citrus NOT grapefruit' | jq .
//   {
//     "boost": 1,
//     "should": {
//       "min": 1,
//       "boost": 1,
//       "disjuncts": [
//         {
//           "boost": 1,
//           "match_phrase": "citrus"
//         },
//         {
//           "boost": 1,
//           "must_not": {
//             "min": 0,
//             "boost": 1,
//             "disjuncts": [
//               {
//                "boost": 1,
//                "match_phrase": "grapefruit"
//              }
//            ]
//          }
//        }
//      ]
//    }
//  }
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bcampbell/qs"
	"os"
	"strings"
)

var defaultAND bool

func main() {

	flag.BoolVar(&defaultAND, "a", false, `Require all terms to match (implied AND rather than OR in queries like "foo bar")`)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-a] query...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parses a query string and dumps the output to stdout.\noptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	queryString := strings.Join(flag.Args(), " ")

	parser := qs.Parser{}
	if defaultAND {
		parser.DefaultOp = qs.AND
	}

	q, err := parser.Parse(queryString)
	if err != nil {
		pe := err.(qs.ParseError)
		fmt.Fprintf(os.Stderr, "%s\n%s^\n",
			queryString,
			strings.Repeat("-", pe.Pos))
		fmt.Fprintf(os.Stderr, "ERR: %s\n", pe)
		os.Exit(2)
	}

	enc, err := json.Marshal(q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "json ERR: %s\n", err)
		os.Exit(2)
	}

	fmt.Println(string(enc))

}
