package engine

import (
	"errors"

	"github.com/kotsmile/jql/internal/parser"
)

type Engine struct {
	loadedTables map[string]Table
}

func New() *Engine {
	return &Engine{
		loadedTables: make(map[string]Table),
	}
}

func (c *Engine) AddTable(name string, table Table) {
	c.loadedTables[name] = table
}

func (c *Engine) Process(query parser.Query) error {
	return errors.New("not implemented")
}
