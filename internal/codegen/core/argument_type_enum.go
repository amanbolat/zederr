// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package core

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

const (
	// ArgumentTypeUnknown is a ArgumentType of type Unknown.
	ArgumentTypeUnknown ArgumentType = iota
	// ArgumentTypeString is a ArgumentType of type String.
	ArgumentTypeString
	// ArgumentTypeInt is a ArgumentType of type Int.
	ArgumentTypeInt
	// ArgumentTypeFloat is a ArgumentType of type Float.
	ArgumentTypeFloat
	// ArgumentTypeBool is a ArgumentType of type Bool.
	ArgumentTypeBool
	// ArgumentTypeTimestamp is a ArgumentType of type Timestamp.
	ArgumentTypeTimestamp
)

var ErrInvalidArgumentType = errors.New("not a valid ArgumentType")

const _ArgumentTypeName = "unknownstringintfloatbooltimestamp"

var _ArgumentTypeMap = map[ArgumentType]string{
	ArgumentTypeUnknown:   _ArgumentTypeName[0:7],
	ArgumentTypeString:    _ArgumentTypeName[7:13],
	ArgumentTypeInt:       _ArgumentTypeName[13:16],
	ArgumentTypeFloat:     _ArgumentTypeName[16:21],
	ArgumentTypeBool:      _ArgumentTypeName[21:25],
	ArgumentTypeTimestamp: _ArgumentTypeName[25:34],
}

// String implements the Stringer interface.
func (x ArgumentType) String() string {
	if str, ok := _ArgumentTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ArgumentType(%d)", x)
}

var _ArgumentTypeValue = map[string]ArgumentType{
	_ArgumentTypeName[0:7]:                    ArgumentTypeUnknown,
	strings.ToLower(_ArgumentTypeName[0:7]):   ArgumentTypeUnknown,
	_ArgumentTypeName[7:13]:                   ArgumentTypeString,
	strings.ToLower(_ArgumentTypeName[7:13]):  ArgumentTypeString,
	_ArgumentTypeName[13:16]:                  ArgumentTypeInt,
	strings.ToLower(_ArgumentTypeName[13:16]): ArgumentTypeInt,
	_ArgumentTypeName[16:21]:                  ArgumentTypeFloat,
	strings.ToLower(_ArgumentTypeName[16:21]): ArgumentTypeFloat,
	_ArgumentTypeName[21:25]:                  ArgumentTypeBool,
	strings.ToLower(_ArgumentTypeName[21:25]): ArgumentTypeBool,
	_ArgumentTypeName[25:34]:                  ArgumentTypeTimestamp,
	strings.ToLower(_ArgumentTypeName[25:34]): ArgumentTypeTimestamp,
}

// ParseArgumentType attempts to convert a string to a ArgumentType.
func ParseArgumentType(name string) (ArgumentType, error) {
	if x, ok := _ArgumentTypeValue[name]; ok {
		return x, nil
	}
	// Case insensitive parse, do a separate lookup to prevent unnecessary cost of lowercasing a string if we don't need to.
	if x, ok := _ArgumentTypeValue[strings.ToLower(name)]; ok {
		return x, nil
	}
	return ArgumentType(0), fmt.Errorf("%s is %w", name, ErrInvalidArgumentType)
}

// MarshalText implements the text marshaller method.
func (x ArgumentType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ArgumentType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseArgumentType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errArgumentTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *ArgumentType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = ArgumentType(0)
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case int64:
		*x = ArgumentType(v)
	case string:
		*x, err = ParseArgumentType(v)
	case []byte:
		*x, err = ParseArgumentType(string(v))
	case ArgumentType:
		*x = v
	case int:
		*x = ArgumentType(v)
	case *ArgumentType:
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x = *v
	case uint:
		*x = ArgumentType(v)
	case uint64:
		*x = ArgumentType(v)
	case *int:
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x = ArgumentType(*v)
	case *int64:
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x = ArgumentType(*v)
	case float64: // json marshals everything as a float64 if it's a number
		*x = ArgumentType(v)
	case *float64: // json marshals everything as a float64 if it's a number
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x = ArgumentType(*v)
	case *uint:
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x = ArgumentType(*v)
	case *uint64:
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x = ArgumentType(*v)
	case *string:
		if v == nil {
			return errArgumentTypeNilPtr
		}
		*x, err = ParseArgumentType(*v)
	}

	return
}

// Value implements the driver Valuer interface.
func (x ArgumentType) Value() (driver.Value, error) {
	return x.String(), nil
}
