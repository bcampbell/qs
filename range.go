package qs

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/numeric_util"
	"strconv"
	"time"
)

// helper stuff for setting up range queries

// returns time, precision
// TODO: support other precisions? "YYYY-MM", "YYYY-MM-DDThh:mm:ss" etc?
func parseTime(in string, loc *time.Location) (time.Time, string) {
	t, err := time.ParseInLocation("2006-01-02", in, loc)
	if err != nil {
		return time.Time{}, ""
	}
	return t, "day"
}

type rangeParams struct {
	min, max                   *string
	minInclusive, maxInclusive *bool
	// need location for parsing times
	loc *time.Location
}

func newRangeParams(minVal, maxVal string, minInc, maxInc bool, loc *time.Location) *rangeParams {
	rp := &rangeParams{}
	if minVal != "" {
		rp.min = &minVal
		rp.minInclusive = &minInc
	}
	if maxVal != "" {
		rp.max = &maxVal
		rp.maxInclusive = &maxInc
	}
	if loc == nil {
		loc = time.UTC
	}
	rp.loc = loc
	return rp
}

func (rp *rangeParams) numericArgs() (bool, *float64, *float64) {
	var f1, f2 *float64
	if rp.min != nil {
		f, err := strconv.ParseFloat(*rp.min, 64)
		if err != nil {
			return false, nil, nil
		}
		f1 = &f
	}
	if rp.max != nil {
		f, err := strconv.ParseFloat(*rp.max, 64)
		if err != nil {
			return false, nil, nil
		}
		f2 = &f
	}
	return true, f1, f2
}

func (rp *rangeParams) dateArgs() (bool, time.Time, time.Time) {
	var truthy bool = true
	var falsey bool = false
	var t1, t2 time.Time
	var prec string
	if rp.min != nil {
		t1, prec = parseTime(*rp.min, rp.loc)
		switch prec {
		case "day":
			if !*rp.minInclusive {
				// add 1 day and make inclusive
				t1 = t1.AddDate(0, 0, 1)
				rp.minInclusive = &truthy
			}
		default:
			return false, t1, t2
		}
	}
	if rp.max != nil {
		t2, prec = parseTime(*rp.max, rp.loc)
		switch prec {
		case "day":
			if *rp.maxInclusive {
				// add 1 day and change to exclusive
				t2 = t2.AddDate(0, 0, 1)
				rp.maxInclusive = &falsey
			}
		default:
			return false, t1, t2
		}
	}

	return true, t1, t2
}

// try and build a query from the given params
func (rp *rangeParams) generate() (bleve.Query, error) {
	if rp.min == nil && rp.max == nil {
		return nil, fmt.Errorf("empty range")
	}
	isNumeric, f1, f2 := rp.numericArgs()
	if isNumeric {
		return bleve.NewNumericRangeInclusiveQuery(f1, f2, rp.minInclusive, rp.maxInclusive), nil
	}

	isDate, t1, t2 := rp.dateArgs()
	if isDate {

		// we'll skip the whole daterange thing and go with the raw timestamp
		// relevant: https://github.com/blevesearch/bleve/issues/251
		var fMin, fMax *float64
		if rp.min != nil {
			foo1 := numeric_util.Int64ToFloat64(t1.UnixNano())
			fMin = &foo1
		}
		if rp.max != nil {
			foo2 := numeric_util.Int64ToFloat64(t2.UnixNano())
			fMax = &foo2
		}
		return bleve.NewNumericRangeInclusiveQuery(fMin, fMax, rp.minInclusive, rp.maxInclusive), nil
	}
	return nil, fmt.Errorf("not numeric")

}
