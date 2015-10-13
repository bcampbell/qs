/*
Package qs is a query language parser for Bleve (http://www.blevesearch.com)

Example:

    import "github.com/bcampbell/qs"

    query, err := qs.Parse("grapefruit lemon orange lime")

Or, to create a parse which uses AND as the default operator:

    p := qs.Parser{DefaultOp: qs.AND}
    query,err := p.Parse("grapefruit lemon orange lime")


*/
package qs
