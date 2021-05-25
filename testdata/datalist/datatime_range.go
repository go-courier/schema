package datalist

import (
	"bytes"
	"encoding"
	"fmt"
	"time"
	_ "time/tzdata"

	"github.com/go-courier/sqlx/v2/builder"
	sqlxDatatypes "github.com/go-courier/sqlx/v2/datatypes"
)

// openapi:strfmt date-time-or-range
type DateTimeOrRange struct {
	From sqlxDatatypes.Timestamp
	To   sqlxDatatypes.Timestamp
	ValueOrRangeOpt
}

func (timeRange DateTimeOrRange) IsZero() bool {
	return timeRange.From.IsZero() && timeRange.To.IsZero()
}

func (timeRange DateTimeOrRange) ConditionFor(c *builder.Column) builder.SqlCondition {
	return timeRange.ValueOrRangeOpt.ConditionFor(c, &timeRange.From, &timeRange.To)
}

func (timeRange DateTimeOrRange) MarshalText() ([]byte, error) {
	return timeRange.ValueOrRangeOpt.MarshalText(&timeRange.From, &timeRange.To)
}

func (timeRange *DateTimeOrRange) UnmarshalText(data []byte) error {
	tr := DateTimeOrRange{}

	err := tr.ValueOrRangeOpt.UnmarshalText(data, &tr.From, &tr.To)
	if err != nil {
		return err
	}

	if !tr.From.IsZero() && !tr.To.IsZero() {
		if time.Time(tr.From).After(time.Time(tr.To)) {
			return fmt.Errorf("time from should not after time to")
		}
	}

	*timeRange = tr
	return nil
}

type Textable interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	IsZero() bool
}

// inspire by https://www.logicbig.com/tutorials/misc/groovy/range-operator.html
type ValueOrRangeOpt struct {
	ExclusiveFrom bool
	ExclusiveTo   bool
	Exactly       bool
}

func (valueOrRange *ValueOrRangeOpt) ConditionFor(c *builder.Column, from Textable, to Textable) builder.SqlCondition {
	where := builder.EmptyCond()

	if valueOrRange != nil {
		if !from.IsZero() {
			if valueOrRange.Exactly {
				return c.Eq(from)
			}

			if valueOrRange.ExclusiveFrom {
				where = where.And(c.Gt(from))
			} else {
				where = where.And(c.Gte(from))
			}
		}

		if !to.IsZero() {
			if valueOrRange.ExclusiveTo {
				where = where.And(c.Lt(to))
			} else {
				where = where.And(c.Lte(to))
			}
		}
	}

	return where
}

func (valueOrRange *ValueOrRangeOpt) UnmarshalText(text []byte, fromValue Textable, toValue Textable) error {
	if len(text) == 0 {
		return nil
	}

	r := ValueOrRangeOpt{}

	spliter := []byte("..")

	r.Exactly = !bytes.Contains(text, spliter)

	fromTo := bytes.Split(text, spliter)

	if len(fromTo) > 0 {
		from := fromTo[0]

		if len(from) > 0 {
			lastChar := from[len(from)-1]
			if lastChar == '!' || lastChar == '<' {
				from = from[:len(from)-1]
				r.ExclusiveFrom = true
			}
		}

		if len(from) > 0 {
			err := fromValue.UnmarshalText(from)
			if err != nil {
				return err
			}
		}
	}

	if len(fromTo) > 1 {
		to := fromTo[1]

		if len(to) > 0 {
			firstChar := to[0]
			if firstChar == '!' || firstChar == '<' {
				to = to[1:]
				r.ExclusiveTo = true
			}
		}

		if len(to) > 0 {
			err := toValue.UnmarshalText(to)
			if err != nil {
				return err
			}
		}
	}

	*valueOrRange = r

	return nil
}

func (valueOrRange ValueOrRangeOpt) MarshalText(fromValue Textable, toValue Textable) (text []byte, err error) {
	buf := bytes.NewBuffer(nil)

	if !fromValue.IsZero() {
		from, err := fromValue.MarshalText()
		if err != nil {
			return nil, err
		}

		buf.Write(from)
		if valueOrRange.ExclusiveFrom {
			buf.WriteByte('<')
		}
	}

	if !valueOrRange.Exactly {
		buf.WriteString("..")

		if !toValue.IsZero() {
			if valueOrRange.ExclusiveTo {
				buf.WriteByte('<')
			}
			to, err := toValue.MarshalText()
			if err != nil {
				return nil, err
			}

			buf.Write(to)
		}
	}

	return buf.Bytes(), nil
}

type RightLikeOrIn []string

func (values RightLikeOrIn) ConditionFor(c *builder.Column) builder.SqlCondition {
	if len(values) == 1 {
		return c.RightLike(values[0])
	}
	return c.In([]string(values))
}
