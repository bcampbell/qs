[![GoDoc](https://godoc.org/github.com/bcampbell/qs?status.svg)](https://godoc.org/github.com/bcampbell/qs)

## Overview

qs implements a query language for the [Bleve](http://www.blevesearch.com/)
text indexer library.

The query syntax is aiming to be about the same as those used by
[Lucene](http://lucene.apache.org/core/5_3_1/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description)
and [Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax).


## Usage

    import "github.com/bcampbell/qs"

    query, err := qs.Parse("grapefruit lemon orange lime")

Or, to create a parse which uses AND as the default operator:

    p := qs.Parser{DefaultOp: qs.AND}
	query,err := p.Parse("grapefruit lemon orange lime")



## Quick Examples

Documents containing the term `grapefruit`:

    grapefruit

There is support for boolean AND, OR and NOT:

    grapefruit OR lemon
    grapefruit AND lemon
    grapefruit AND NOT lemon

Phrases are indicated with quotes:

   "navel orange"

The default boolean operator is OR, so these two are equivalent when using the default parser:

    grapefruit lemon
    grapefruit OR lemon


Fields:

    tags:citrus
    headline:"How to Make the Perfect Negroni"

Grouping:

    (lemon AND lime) OR (orange AND grapefruit)

Boosting term relevance:

    fruit grapefruit^2

Inclusive ranges in square brackets, Exclusive ranges in curly brackets:

    date:[2012-01-01 TO 2012-12-31]
    count:{0 TO 100}

Wildcard markers `*` (any character sequence, 0 or more) and `?` (any single character)
eg to match `foot`, `fort`, `fret` etc:

    f??t

eg to match `ant`, `anteater`, `antidisestablishmentarianism` etc:

    ant*



