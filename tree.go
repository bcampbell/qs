package qs

import (
	"fmt"
	"github.com/blevesearch/bleve"
)

type QueryKind int

const (
	MatchNone QueryKind = iota
	Boolean
	Disjunction
	Conjunction
	Wildcard
	Fuzzy
	MatchPhrase
	NumericRangeInclusive
)

func (k QueryKind) String() string {
	switch k {
	case MatchNone:
		return "MatchNone"
	case Boolean:
		return "Boolean"
	case Disjunction:
		return "Disjunction"
	case Conjunction:
		return "Conjunction"
	case Wildcard:
		return "Wildcard"
	case Fuzzy:
		return "Fuzzy"
	case MatchPhrase:
		return "MatchPhrase"
	case NumericRangeInclusive:
		return "NumericRangeInclusive"
	default:
		return ""
	}
}

// Query is an intermediate representation, for constructing a bleve.Query from
// The fields used depend on the type of query.
type Query struct {
	Kind QueryKind
	Txt  string
	// Children holds sets of subqueries.
	// For Boolean queries, these sets are (must,should,mustNot)
	// For Conjunction/Disjunction queries, only the first set is used.
	// Other query types don't have children.
	Children  [3][]*Query
	Boost     float64 // default should be 1.0!
	Field     string
	Fuzziness int

	// icky bits for ranges
	F1, F2                     *float64
	MinInclusive, MaxInclusive *bool
}

func (q *Query) SetField(field string) *Query {
	q.Field = field
	return q
}

func (q *Query) SetBoost(boost float64) *Query {
	q.Boost = boost
	return q
}

func (q *Query) SetFuzziness(fuzziness int) *Query {
	q.Fuzziness = fuzziness
	return q
}

func groupToBleve(grp []*Query) []bleve.Query {
	out := make([]bleve.Query, len(grp))
	for i, child := range grp {
		out[i] = child.ToBleve()
	}
	return out
}

// ToBleve converts the intermediate query tree into a bleve.Query
func (q *Query) ToBleve() bleve.Query {

	var bq bleve.Query
	switch q.Kind {
	case MatchNone:
		bq = bleve.NewMatchNoneQuery()
	case Boolean:
		must := groupToBleve(q.Children[0])
		should := groupToBleve(q.Children[1])
		mustNot := groupToBleve(q.Children[2])
		bq = bleve.NewBooleanQuery(must, should, mustNot)
	case Disjunction:
		bq = bleve.NewDisjunctionQuery(groupToBleve(q.Children[0]))
	case Conjunction:
		bq = bleve.NewConjunctionQuery(groupToBleve(q.Children[0]))
	case Wildcard:
		bq = bleve.NewWildcardQuery(q.Txt)
	case Fuzzy:
		bq = bleve.NewFuzzyQuery(q.Txt).SetFuzziness(q.Fuzziness)
	case MatchPhrase:
		bq = bleve.NewMatchPhraseQuery(q.Txt)
	case NumericRangeInclusive:
		bq = bleve.NewNumericRangeInclusiveQuery(q.F1, q.F2, q.MinInclusive, q.MaxInclusive)
	}

	return bq.SetField(q.Field).SetBoost(q.Boost)
}

func NewMatchNoneQuery() *Query {
	return &Query{Kind: MatchNone, Boost: 1.0}
}

func NewBooleanQuery(must, should, mustNot []*Query) *Query {
	return &Query{Kind: Boolean, Boost: 1.0, Children: [3][]*Query{must, should, mustNot}}
}

func NewDisjunctionQuery(disjuncts []*Query) *Query {
	return &Query{Kind: Disjunction, Boost: 1.0, Children: [3][]*Query{disjuncts, []*Query{}, []*Query{}}}
}

func NewConjunctionQuery(conjuncts []*Query) *Query {

	return &Query{Kind: Conjunction, Boost: 1.0, Children: [3][]*Query{conjuncts, []*Query{}, []*Query{}}}
}

func NewWildcardQuery(wildcard string) *Query {
	return &Query{Kind: Wildcard, Boost: 1.0, Txt: wildcard}
}

func NewFuzzyQuery(term string) *Query {
	return &Query{Kind: Fuzzy, Boost: 1.0, Fuzziness: 2, Txt: term}
}

func NewMatchPhraseQuery(matchPhrase string) *Query {
	return &Query{Kind: MatchPhrase, Boost: 1.0, Txt: matchPhrase}
}

func NewNumericRangeInclusiveQuery(f1, f2 *float64, minInclusive, maxInclusive *bool) *Query {
	return &Query{
		Kind:         NumericRangeInclusive,
		Boost:        1.0,
		F1:           f1,
		F2:           f2,
		MinInclusive: minInclusive,
		MaxInclusive: maxInclusive}
}

func (q *Query) Dump() {
	q.dump("")
}

func (q *Query) dump(indent string) {

	fmt.Printf("%s%s ", indent, q.Kind.String())

	if q.Field != "" {
		fmt.Printf("field: %q ", q.Field)
	}
	if q.Txt != "" {
		fmt.Printf("txt: %q ", q.Txt)
	}
	fmt.Printf("boost: %f", q.Boost)

	fmt.Printf("\n")
	for i := 0; i < 3; i++ {
		if len(q.Children[i]) == 0 {
			continue
		}
		fmt.Printf("%s children[%d]:\n", indent, i)
		for _, child := range q.Children[i] {
			child.dump(indent + "  ")
		}
	}
}
