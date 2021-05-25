package expressions

type Kind int

const (
	Invalid Kind = iota
	Uint
	Int
	Float
	String
	Bool
)

type Valuer interface {
	Kind() Kind
	Value() interface{}
}

type FloatValue float64

func (FloatValue) Kind() Kind {
	return Float
}

func (f FloatValue) Value() interface{} {
	return float64(f)
}

func Equal(b Valuer, v Valuer) bool {
	if b.Kind() == v.Kind() {
		return b.Value() == v.Value()
	}

	switch b.Kind() {
	case Float:
		switch v.Kind() {
		case Int:
			return b.Value().(float64) == float64(v.Value().(int64))
		case Uint:
			return b.Value().(float64) == float64(v.Value().(uint64))
		}
	case Int:
		if v.Kind() == Uint {
			return b.Value().(int64) == int64(v.Value().(uint64))
		}
	case Uint:
		if v.Kind() == Int {
			return int64(b.Value().(uint64)) == v.Value().(int64)
		}
	}

	return false
}

func GreaterThan(b Valuer, v Valuer) bool {
	switch b.Kind() {
	case Float:
		switch v.Kind() {
		case Float:
			return v.Value().(float64) > b.Value().(float64)
		case Int:
			return float64(v.Value().(int64)) > b.Value().(float64)
		case Uint:
			return float64(v.Value().(uint64)) > b.Value().(float64)
		}
	case Int:
		switch v.Kind() {
		case Float:
			return v.Value().(float64) > float64(b.Value().(int64))
		case Int:
			return v.Value().(int64) > b.Value().(int64)
		case Uint:
			return int64(v.Value().(uint64)) > b.Value().(int64)
		}
	case Uint:
		switch v.Kind() {
		case Float:
			return v.Value().(float64) > float64(b.Value().(uint64))
		case Int:
			return v.Value().(int64) > int64(b.Value().(uint64))
		case Uint:
			return v.Value().(uint64) > b.Value().(uint64)
		}
	case String:
		if v.Kind() == String {
			return v.Value().(string) > b.Value().(string)
		}
	}
	return false
}

func LessThan(b Valuer, v Valuer) bool {
	switch b.Kind() {
	case Float:
		switch v.Kind() {
		case Float:
			return v.Value().(float64) < b.Value().(float64)
		case Int:
			return float64(v.Value().(int64)) < b.Value().(float64)
		case Uint:
			return float64(v.Value().(uint64)) < b.Value().(float64)
		}
	case Int:
		switch v.Kind() {
		case Float:
			return v.Value().(float64) < float64(b.Value().(int64))
		case Int:
			return v.Value().(int64) < b.Value().(int64)
		case Uint:
			return int64(v.Value().(uint64)) < b.Value().(int64)
		}
	case Uint:
		switch v.Kind() {
		case Float:
			return v.Value().(float64) < float64(b.Value().(uint64))
		case Int:
			return v.Value().(int64) < int64(b.Value().(uint64))
		case Uint:
			return v.Value().(uint64) < b.Value().(uint64)
		}
	case String:
		if v.Kind() == String {
			return v.Value().(string) < b.Value().(string)
		}
	}
	return false
}

type IntValue int64

func (IntValue) Kind() Kind {
	return Int
}

func (f IntValue) Value() interface{} {
	return int64(f)
}

type UintValue uint64

func (UintValue) Kind() Kind {
	return Uint
}

func (f UintValue) Value() interface{} {
	return uint64(f)
}

type StringValue string

func (StringValue) Kind() Kind {
	return String
}

func (f StringValue) Value() interface{} {
	return string(f)
}

type BoolValue bool

func (BoolValue) Kind() Kind {
	return Bool
}

func (f BoolValue) Value() interface{} {
	return bool(f)
}

func normalize(in interface{}) Valuer {
	switch v := in.(type) {
	case float64:
		return FloatValue(v)
	case float32:
		return FloatValue(v)
	case int:
		return IntValue(v)
	case int8:
		return IntValue(v)
	case int16:
		return IntValue(v)
	case int32:
		return IntValue(v)
	case int64:
		return IntValue(v)
	case uint:
		return UintValue(v)
	case uint8:
		return UintValue(v)
	case uint16:
		return UintValue(v)
	case uint32:
		return UintValue(v)
	case uint64:
		return UintValue(v)
	case string:
		return StringValue(v)
	case bool:
		return BoolValue(v)
	}
	return nil
}
