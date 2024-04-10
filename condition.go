package querystring

import (
	"github.com/gobwas/glob"
	"regexp"
	"strings"
	"time"
)

// Condition .
type Condition interface {
}

// AndCondition .
type AndCondition struct {
	Left  Condition
	Right Condition
}

// NewAndCondition .
func NewAndCondition(left, right Condition) *AndCondition {
	return &AndCondition{Left: left, Right: right}
}

// OrCondition .
type OrCondition struct {
	Left  Condition
	Right Condition
}

// NewOrCondition .
func NewOrCondition(left, right Condition) *OrCondition {
	return &OrCondition{Left: left, Right: right}
}

// NotCondition .
type NotCondition struct {
	Condition Condition
}

// NewNotCondition .
func NewNotCondition(q Condition) *NotCondition {
	return &NotCondition{Condition: q}
}

// FieldableCondition .
type FieldableCondition interface {
	GetField() string
	SetField(field string)
}

// MatchCondition .
type MatchCondition struct {
	Field string
	Value string
}

// NewMatchCondition .
func NewMatchCondition(s string) *MatchCondition {
	return &MatchCondition{
		Value: s,
	}
}

// GetField .
func (q *MatchCondition) GetField() string {
	return q.Field
}

// SetField .
func (q *MatchCondition) SetField(field string) {
	q.Field = field
}

// RegexpCondition .
type RegexpCondition struct {
	Field string
	Value *regexp.Regexp
}

// MustNewRegexpCondition panics and must be used only inside goyacc, due it recovers panics into goyacc errors
func MustNewRegexpCondition(s string) *RegexpCondition {
	return &RegexpCondition{
		Value: regexp.MustCompile(s),
	}
}

// NewRegexpCondition .
func NewRegexpCondition(s string) (*RegexpCondition, error) {
	r, err := regexp.Compile(s)

	if err != nil {
		return nil, err
	}

	return &RegexpCondition{
		Value: r,
	}, nil
}

// GetField .
func (q *RegexpCondition) GetField() string {
	return q.Field
}

// SetField .
func (q *RegexpCondition) SetField(field string) {
	q.Field = field
}

// WildcardCondition .
type WildcardCondition struct {
	Field string
	Value glob.Glob
}

// MustNewWildcardCondition panics and must be used only inside goyacc, due it recovers panics into goyacc errors
func MustNewWildcardCondition(s string) *WildcardCondition {
	return &WildcardCondition{
		Value: glob.MustCompile(s),
	}
}

// NewWildcardCondition .
func NewWildcardCondition(s string) (*WildcardCondition, error) {
	g, err := glob.Compile(s)

	if err != nil {
		return nil, err
	}
	return &WildcardCondition{
		Value: g,
	}, nil
}

// GetField .
func (q *WildcardCondition) GetField() string {
	return q.Field
}

// SetField .
func (q *WildcardCondition) SetField(field string) {
	q.Field = field
}

// NumberRangeCondition .
type NumberRangeCondition struct {
	Field        string
	Start        *string
	End          *string
	IncludeStart bool
	IncludeEnd   bool
}

// NewNumberRangeCondition .
func NewNumberRangeCondition(start, end *string, includeStart, includeEnd bool) *NumberRangeCondition {
	return &NumberRangeCondition{
		Start:        start,
		End:          end,
		IncludeStart: includeStart,
		IncludeEnd:   includeEnd,
	}
}

// GetField .
func (q *NumberRangeCondition) GetField() string {
	return q.Field
}

// SetField .
func (q *NumberRangeCondition) SetField(field string) {
	q.Field = field
}

// TimeRangeCondition .
type TimeRangeCondition struct {
	Field        string
	Start        *string
	End          *string
	IncludeStart bool
	IncludeEnd   bool
}

// NewTimeRangeCondition .
func NewTimeRangeCondition(start, end *string, includeStart, includeEnd bool) *TimeRangeCondition {
	return &TimeRangeCondition{
		Start:        start,
		End:          end,
		IncludeStart: includeStart,
		IncludeEnd:   includeEnd,
	}
}

// GetField .
func (q *TimeRangeCondition) GetField() string {
	return q.Field
}

// SetField .
func (q *TimeRangeCondition) SetField(field string) {
	q.Field = field
}

func queryTimeFromString(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}

func mustNewStringCondition(str string) FieldableCondition {
	if strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/") {
		return MustNewRegexpCondition(str[1 : len(str)-1])
	}

	if strings.ContainsAny(str, "*?") {
		return MustNewWildcardCondition(str)
	}

	return NewMatchCondition(str)
}

func newStringCondition(str string) (FieldableCondition, error) {
	if strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/") {
		return NewRegexpCondition(str[1 : len(str)-1])
	}

	if strings.ContainsAny(str, "*?") {
		return NewWildcardCondition(str)
	}

	return NewMatchCondition(str), nil
}
