package enumeration

type WithConstValues interface {
	ConstValues() []IntStringerEnum
}

type IntStringerEnum interface {
	WithConstValues

	TypeName() string
	Int() int
	String() string
	Label() string
}

// Deprecated use IntStringerEnum instead
type Enum = IntStringerEnum

// sql value of enum maybe have offset from value of enum in go
type DriverValueOffset interface {
	Offset() int
}

// Deprecated use DriverValueOffset instead
type EnumDriverValueOffset = DriverValueOffset
