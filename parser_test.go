package querystring

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	cond, err := Parse("message: test\\ value AND datetime: [\"2020-01-01T00:00:00\" TO \"2020-12-31T00:00:00\"]")
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	expected := Condition(&AndCondition{
		Left: &MatchCondition{
			Field: "message",
			Value: "test value",
		},
		Right: &TimeRangeCondition{
			Field:        "datetime",
			Start:        sptr("2020-01-01T00:00:00"),
			End:          sptr("2020-12-31T00:00:00"),
			IncludeStart: true,
			IncludeEnd:   true,
		},
	})

	assert.Equal(t, expected, cond)
}

func TestParserPhase(t *testing.T) {
	cond, err := Parse(`a: "/.*/"`)
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	assert.Equal(t, &RegexpCondition{
		Field: "a",
		Value: regexp.MustCompile(".*"),
	}, cond)
}

func TestParserMixedCondition(t *testing.T) {
	cond, err := Parse("a: 1 OR (b: 2 and c: 4)")
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	assert.Equal(t, &OrCondition{
		Left: &NumberRangeCondition{
			Field:        "a",
			Start:        iptr(int64(1)),
			End:          iptr(int64(1)),
			IncludeStart: true,
			IncludeEnd:   true,
		},
		Right: &AndCondition{
			Left: &NumberRangeCondition{
				Field:        "b",
				Start:        iptr(int64(2)),
				End:          iptr(int64(2)),
				IncludeStart: true,
				IncludeEnd:   true,
			},
			Right: &NumberRangeCondition{
				Field:        "c",
				Start:        iptr(int64(4)),
				End:          iptr(int64(4)),
				IncludeStart: true,
				IncludeEnd:   true,
			},
		},
	}, cond)
}

func iptr(i int64) *int64 {
	return &i
}

func sptr(s string) *string {
	return &s
}
