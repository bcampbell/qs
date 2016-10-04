# Query String Syntax

## Basics

Single words are treated as simple terms. The results will include
documents containing those terms, eg:

    embiggen
    cromulent

Terms are separated by spaces, so:

    navel orange

Will match both documents containing "navel" and documents containing "orange",
which is probably not what was intended. The user was probably looking for
documents about "navel oranges" rather than those about "orange navels".

To search for phrases - multiple words, matched in order - enclose them in double quotes.

For example:

    "navel orange"

## Boolean Operators

### `OR`

The `OR` operator returns documents matched by terms on either side of it.
It is the default operator, meaning that the following two queries are
considered to be equivalent:

    orange OR "navel orange"
    orange "navel orange"

Note that the query parser can be set up to use `AND` as the default operator instead of `OR`.

### `AND`

'`AND`' returns documents which match the terms on *both* sides.

    orange AND citrus

### `+`

'`+`' is the 'required' operator. A term prefixed with '`+`' must exist in matching
documents.

For example, to return documents definitely containing "orange" and maybe also
containing "fruit":

    fruit +orange


### `NOT`

`NOT` excludes documents that match the following term.

For example, to match documents containing "orange" but not "paint":

    orange NOT paint


### `-`

Prefixing a term with '`-`' will "prohibit" it.

For example, to match documents containing "orange" but not "paint":

    orange -paint


### Precedence

`-`, `+` and `NOT` take precedence over `AND`, which takes precedence over `OR`.

For example,

    lemon OR orange AND citrus

is treated as:

    lemon OR (orange AND citrus)



## Field Scoping

You can control which fields are matched by prefixing the name of a field separated by a colon.

Examples:

    genus:citrus
    headline:"How to Make the Perfect Negroni"
    tags:(fruit OR paint)

If no field is specified, the default `_all` is used, unless overridden in the index mapping.


## Grouping

Parentheses can be used to group sub queries.
For example:

    (shaddock OR pomelo OR pamplemousse) AND (family:rutaceae AND NOT genus:fortunella) AND colour:(greenish OR yellowish)

## Wildcards

Within an individual term, partial matches can be described using wildcard characters:

`?` to match any single character

`*` To match and sequence of zero or more characters

For example:

     qu?ck bro*


## Fuzziness

A fuzzy query is a term query that matches terms within a given [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance).
The edit distance is the number of single-character edits (insertions, deletions or substitutions) allowed.

To specify a fuzzy query, use the tilde sign (`~`), optionally followed by the edit distance you'll accept.


For example,

    colour~1

to match "`colour`" or "`color`" (or "`zolour`" or "`colours`"... but not "`colors`", as the `~1` allows only a single character change).

If no number is specified, the default value is 1, so the following are equivalent:

    wibble~
    wibble~1



## Boosting

Boosting a term doesn't affect the set of matching documents, but it does
indicate that that term should be assigned greater importance when the
results are being scored by relevance.

To boost a term, you use the boost operator '`^`' followed by a number
indicating the relative importance of the term.

For example, in the following query, "grapefruit" should be considered twice
as important as "lime":

    grapefruit^2 lime

You can boost more complex sub queries:

    (grapefruit OR orange)^2 "navel orange"^4 genus:citrus



## Ranges

Inclusive ranges can be described with square braces and `TO`, eg:

    num:[1 TO 5]
    date:[2010-01-01 TO 2010-01-31]

Dates must be in `YYYY-MM-DD` form.

Exclusive ranges are supported using curly braces. So, these are equivalent:

    score:{0 TO 10}
    score:[1 TO 9]



You can mix inclusive and exclusive endpoints, eg:

    pubdate:[2010-01-01 TO 2011-01-01}
    shoesize:{0 TO 16]
    byte:[0 TO 256}

You can have open ranges by leaving off either endpoint:

    pubdate:[2000-01-01 TO ]
    temp:[TO 100}


NOTE: ranges currently work only on numeric and date fields.


## Relational Operators

You can perform numeric comparisons using the `>`, `>=`, `<`, and `<=` operators.
These are equivalent to using the above range syntax with unbounded ranges.

For example:

    score:>=100
    score:[ TO 100]


