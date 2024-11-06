package engine

import (
	"errors"
	"fmt"
)

var (
	ErrDataIsNotArray  = errors.New("data is not an array")
	ErrEmptyArray      = errors.New("empty array")
	ErrDataIsNotObject = errors.New("data is not a object")
)

type columnType string

const (
	NumberType  columnType = "number"
	StringType  columnType = "string"
	BooleanType columnType = "boolean"
	ArrayType   columnType = "array"
	ObjectType  columnType = "object"
	NullType    columnType = "null"
)

type (
	columns map[string]columnType
	rows    []Row
)

type Table struct {
	columns columns
	rows    rows
}

type Row map[string]any

func eqColumns(a, b columns) bool {
	if len(a) != len(b) {
		return false
	}

	for name, aType := range a {
		if bType, ok := b[name]; !ok || aType != bType {
			return false
		}
	}

	return true
}

func parseColumns(data Row) (columns, error) {
	columns := make(columns)

	for key, value := range data {
		switch value.(type) {
		case string:
			columns[key] = StringType
		case float64:
			columns[key] = NumberType
		case int:
			columns[key] = NumberType
		case bool:
			columns[key] = BooleanType
		case []any:
			columns[key] = ArrayType
		case map[string]any:
			columns[key] = ObjectType
		case nil:
			columns[key] = NullType
		default:
			return nil, fmt.Errorf("unknown type: %T", value)
		}
	}

	return columns, nil
}

// TODO: add support for nested arrays
// TODO: add support for nested objects
func NewTable(rows []Row) (*Table, error) {
	if len(rows) == 0 {
		return nil, ErrEmptyArray
	}

	columns, err := parseColumns(rows[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse columns for %v: %w", rows[0], err)
	}

	// TODO: add support for different row types
	for _, row := range rows[1:] {
		c, err := parseColumns(row)
		if err != nil {
			return nil, fmt.Errorf("failed to parse columns for %v: %w", row, err)
		}

		if !eqColumns(columns, c) {
			return nil, fmt.Errorf("rows have different schema %v != %v", columns, c)
		}
	}

	return &Table{
		columns: columns,
		rows:    rows,
	}, nil
}
