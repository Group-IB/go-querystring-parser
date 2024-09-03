package querystring

import (
	"github.com/gobwas/glob"
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

func TestParserWildcardPhase(t *testing.T) {
	cond1, err := Parse(`a: "\\/*/abc/"`)
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	assert.Equal(t, &WildcardCondition{
		Field: "a",
		Value: glob.MustCompile("/*/abc/"),
	}, cond1)

	cond2, err := Parse(`b: "\\/abc/?/"`)
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	assert.Equal(t, &WildcardCondition{
		Field: "b",
		Value: glob.MustCompile("/abc/?/"),
	}, cond2)

	cond3, err := Parse(`c: "\\/abc/[!a-z]/[!1-5]"`)
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	assert.Equal(t, &WildcardCondition{
		Field: "c",
		Value: glob.MustCompile("/abc/[!a-z]/[!1-5]"),
	}, cond3)

	cond4, err := Parse(`a: "*/AbC/CrAzY/CaSE/*"`, WithOptions(&Options{
		LowerCaseWildcard: true,
	}))
	if err != nil {
		t.Errorf("parse return error, %s", err)
		return
	}

	assert.Equal(t, &WildcardCondition{
		Field: "a",
		Value: glob.MustCompile("*/abc/crazy/case/*"),
	}, cond4)
}

func TestParserRegexPhase(t *testing.T) {
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
