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

type columnDefinition struct {
	ColumnType columnType
}

type (
	columns map[string]columnDefinition
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

func parseColumns(data []Row) (columns, error) {
	columnsCollection := make(map[string][]columnDefinition)
	for _, row := range data {
		for key, value := range row {
			switch value.(type) {
			case string:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: StringType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == StringType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: StringType,
						})
					}
				}
			case float64:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: NumberType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == NumberType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: NumberType,
						})
					}
				}
			case int:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: NumberType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == NumberType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: NumberType,
						})
					}
				}
			case bool:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: BooleanType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == BooleanType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: BooleanType,
						})
					}
				}
			case []any:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: ArrayType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == ArrayType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: ArrayType,
						})
					}
				}
			case map[string]any:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: ObjectType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == ObjectType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: ObjectType,
						})
					}
				}
			case nil:
				if cv, ok := columnsCollection[key]; !ok {
					columnsCollection[key] = append(columnsCollection[key], columnDefinition{
						ColumnType: NullType,
					})
				} else {
					found := false
					for _, c := range cv {
						if c.ColumnType == NullType {
							found = true
							break
						}
					}

					if !found {
						columnsCollection[key] = append(columnsCollection[key], columnDefinition{
							ColumnType: NullType,
						})
					}
				}
			default:
				return nil, fmt.Errorf("unknown type: %T", value)
			}
		}
	}

	columns := make(columns)
	for key, columnDefinitions := range columnsCollection {
		if len(columnDefinitions) == 1 {
			columns[key] = columnDefinitions[0]
		} else {
			var ts []columnType

			for _, columnDefinition := range columnDefinitions {
				if columnDefinition.ColumnType != NullType {
					found := false
					for _, t := range ts {
						if t == columnDefinition.ColumnType {
							found = true
							break
						}
					}
					if !found {
						ts = append(ts, columnDefinition.ColumnType)
					}
				}
			}
			if len(ts) == 1 {
				columns[key] = columnDefinition{
					ColumnType: ts[0],
				}
			} else {
				columns[key] = columnDefinition{
					ColumnType: StringType,
				}
			}
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

	columns, err := parseColumns(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to parse columns for %v: %w", rows[0], err)
	}

	// // TODO: add support for different row types
	// for _, row := range rows[1:] {
	// 	c, err := parseColumns(row)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to parse columns for %v: %w", row, err)
	// 	}
	//
	// 	if !eqColumns(columns, c) {
	// 		return nil, fmt.Errorf("rows have different schema %v != %v", columns, c)
	// 	}
	// }

	return &Table{
		columns: columns,
		rows:    rows,
	}, nil
}

func (t *Table) Columns() columns {
	return t.columns
}

func (t *Table) Rows() rows {
	return t.rows
}

type SqliteColumn struct {
	Name       string
	SqliteType string
}

func (t *Table) ToSqliteTypes() []SqliteColumn {
	var cs []SqliteColumn
	for name, column := range t.columns {
		dbType := ""
		switch column.ColumnType {
		case StringType:
			dbType = "TEXT"
		case NumberType:
			dbType = "REAL"
		case BooleanType:
			dbType = "BOOLEAN"
		case ArrayType:
			dbType = "TEXT"
		case ObjectType:
			dbType = "TEXT"
		case NullType:
			dbType = "TEXT"
		}
		cs = append(cs, SqliteColumn{
			Name:       name,
			SqliteType: dbType,
		})
	}
	return cs
}
