package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kotsmile/jql/internal/parser"
	"github.com/kotsmile/jql/internal/tableui"
)

var (
	ErrLoadMissingFilename   = errors.New("'load' command: missing filename")
	ErrLoadWrongTypeFilename = errors.New("'load' command: wrong type for filename")
	ErrLoadAsExpected        = errors.New("'load' command: expected 'as' keyword")
	ErrLoadExpectedTableName = errors.New("'load' command: expected table name")
)

type Engine struct {
	loadedTables map[string]*Table
	writer       io.Writer
}

func New(w io.Writer) *Engine {
	return &Engine{
		loadedTables: make(map[string]*Table),
		writer:       w,
	}
}

func (e *Engine) Process(query *parser.AstNode) error {
	switch value := query.Value().(type) {
	case *parser.KeywordNode:
		switch value.Value() {
		case parser.LoadKeyword.String():
			if len(query.Children()) < 1 {
				return ErrLoadMissingFilename
			}

			filename, ok := query.Children()[0].Value().(parser.StringNode)
			if !ok {
				return ErrLoadWrongTypeFilename
			}

			index := strings.Index(filename.Value(), ".json")
			var tablename parser.StringNode = parser.StringNode(filepath.Base(filename.Value()[:index]))

			if len(query.Children()) > 1 {
				asCommand, ok := query.Children()[1].Value().(*parser.KeywordNode)
				if !ok {
					return ErrLoadAsExpected
				}
				if asCommand.Value() != parser.AsKeyword.String() {
					return ErrLoadAsExpected
				}

				if len(query.Children()[1].Children()) < 1 {
					return ErrLoadExpectedTableName
				}

				tablename, ok = query.Children()[1].Children()[0].Value().(parser.StringNode)
				if !ok {
					return ErrLoadExpectedTableName
				}
			}

			if err := e.loadCommand(filename.Value(), tablename.Value()); err != nil {
				return fmt.Errorf("failed to load table: %w", err)
			}

			e.writer.Write([]byte(fmt.Sprintf("Loaded table '%s'\n", tablename.Value())))
		case parser.TablesKeyword.String():
			for name := range e.loadedTables {
				e.writer.Write([]byte(fmt.Sprintf("%s\n", name)))
			}
		case parser.SelectKeyword.String():
			if len(query.Children()) < 1 {
				return fmt.Errorf("'select' command: missing columns")
			}

			fromNode := query.Children()[len(query.Children())-1]
			if _, ok := fromNode.Value().(*parser.KeywordNode); !ok {
				return fmt.Errorf("'select' command: missing 'from' keyword")
			}

			columns := make([]string, 0)
			for _, c := range query.Children()[:len(query.Children())-1] {
				c, ok := c.Value().(parser.StringNode)
				if !ok {
					return fmt.Errorf("'select' command: wrong type for column name")
				}
				columns = append(columns, c.Value())
			}

			if len(fromNode.Children()) != 1 {
				return fmt.Errorf("'select' command: wrong number of arguments for 'from' keyword")
			}

			tablename, ok := fromNode.Children()[0].Value().(parser.StringNode)
			if !ok {
				return fmt.Errorf("'select' command: wrong type for table name")
			}
			if err := e.selectCommand(tablename.Value(), columns); err != nil {
				return fmt.Errorf("failed to select columns: %w", err)
			}

		}

	default:
		return errors.New("unknown command")
	}

	return nil
}

func (c *Engine) LoadTable(filename string, tablename string) error {
	return c.loadCommand(filename, tablename)
}

func (c *Engine) loadCommand(filename string, tablename string) error {
	_, ok := c.loadedTables[tablename]
	if ok {
		return fmt.Errorf("table '%s' already loaded", tablename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var rows []Row

	if err := json.NewDecoder(file).Decode(&rows); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	table, err := NewTable(rows)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	c.loadedTables[tablename] = table

	return nil
}

func (c *Engine) selectCommand(tablename string, columns []string) error {
	table, ok := c.loadedTables[tablename]
	if !ok {
		return fmt.Errorf("table '%s' not found", tablename)
	}

	allColumns := make([]string, 0)
	for column := range table.columns {
		allColumns = append(allColumns, column)
	}

	unpackedColumns := make([]string, 0)
	for _, column := range columns {
		if column == "*" {
			unpackedColumns = append(unpackedColumns, allColumns...)
		} else {
			unpackedColumns = append(unpackedColumns, column)
		}
	}

	for _, column := range unpackedColumns {
		if _, ok := table.columns[column]; !ok {
			return fmt.Errorf("column '%s' not found in table '%s'", column, tablename)
		}
	}

	var cs []string
	var rs []tableui.Row

	for _, column := range unpackedColumns {
		cs = append(cs, column)
	}

	for _, row := range table.rows {
		r := make(tableui.Row, 0)
		for _, column := range unpackedColumns {
			r = append(r, fmt.Sprintf("%v", row[column]))
		}
		rs = append(rs, r)
	}

	tui := tableui.New(cs, rs)
	if err := tui.Render(c.writer); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
}
