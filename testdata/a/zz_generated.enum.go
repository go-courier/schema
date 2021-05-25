package a

import (
	"bytes"
	driver "database/sql/driver"

	enumeration "github.com/go-courier/schema/pkg/enumeration"
	errors "github.com/pkg/errors"
)

var InvalidProtocol = errors.New("invalid Protocol")

func ParseProtocolFromString(s string) (Protocol, error) {
	switch s {
	case "HTTPS":
		return PROTOCOL__HTTPS, nil
	case "TCP":
		return PROTOCOL__TCP, nil
	case "HTTP":
		return PROTOCOL__HTTP, nil
	}
	return PROTOCOL_UNKNOWN, InvalidProtocol
}

func ParseProtocolFromLabelString(s string) (Protocol, error) {
	switch s {
	case "https":
		return PROTOCOL__HTTPS, nil
	case "tcp":
		return PROTOCOL__TCP, nil
	case "http":
		return PROTOCOL__HTTP, nil
	}
	return PROTOCOL_UNKNOWN, InvalidProtocol
}

func (Protocol) TypeName() string {
	return "github.com/go-courier/schema/testdata/a.Protocol"
}

func (v Protocol) String() string {
	switch v {
	case PROTOCOL__HTTPS:
		return "HTTPS"
	case PROTOCOL__TCP:
		return "TCP"
	case PROTOCOL__HTTP:
		return "HTTP"
	}
	return "UNKNOWN"
}

func (v Protocol) Label() string {
	switch v {
	case PROTOCOL__HTTPS:
		return "https"
	case PROTOCOL__TCP:
		return "tcp"
	case PROTOCOL__HTTP:
		return "http"
	}
	return "UNKNOWN"
}

func (v Protocol) Int() int {
	return int(v)
}

func (Protocol) ConstValues() []enumeration.IntStringerEnum {
	return []enumeration.IntStringerEnum{
		PROTOCOL__HTTPS,
		PROTOCOL__TCP,
		PROTOCOL__HTTP,
	}
}

func (v Protocol) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidProtocol
	}
	return []byte(str), nil
}

func (v *Protocol) UnmarshalText(data []byte) (err error) {
	*v, err = ParseProtocolFromString(string(bytes.ToUpper(data)))
	return
}

func (v Protocol) Value() (driver.Value, error) {
	offset := 0
	if o, ok := (interface{})(v).(enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}
	return int64(v) + int64(offset), nil
}

func (v *Protocol) Scan(src interface{}) error {
	offset := 0
	if o, ok := (interface{})(v).(enumeration.DriverValueOffset); ok {
		offset = o.Offset()
	}

	i, err := enumeration.ScanIntEnumStringer(src, offset)
	if err != nil {
		return err
	}
	*v = Protocol(i)
	return nil
}

func (Protocol) OpenAPISchemaEnum() []interface{} {
	return []interface{}{
		PROTOCOL__HTTPS,
		PROTOCOL__TCP,
		PROTOCOL__HTTP,
	}
}
