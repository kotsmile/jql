package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kotsmile/jql/internal/engine"
	"github.com/kotsmile/jql/internal/lexer"
	"github.com/kotsmile/jql/internal/parser"
	"github.com/kotsmile/jql/util"
)

func processCmd(cmd string, e *engine.Engine, logger util.Logger) error {
	l := lexer.New(logger.WithField("module", "lexer"))
	l.Lex(cmd)

	p := parser.New(l, logger.WithField("module", "parser"))
	queries, err := p.Parse()
	if err != nil {
		return err
	}

	for _, q := range queries {
		if err := e.Process(q); err != nil {
			continue
		}
	}

	return nil
}

func main() {
	debug := false

	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	logger := util.NewLogger(debug)
	e := engine.New(os.Stdout)

	e.LoadTable("./examples/simple.json", "simple")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if err := processCmd(cmd, e, logger); err != nil {
			logger.Errorf("failed to execute command: %s", err)
		}
	}
}
