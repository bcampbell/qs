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
	"fmt"
	"github.com/bcampbell/qs"
	"os"
	"strings"
)

func main() {

	queryString := strings.Join(os.Args[1:], " ")

	q, err := qs.Parse(queryString)
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
