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

func main() {
	debug := false

	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	logger := util.NewLogger(debug)
	e := engine.New()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		l := lexer.New(logger.WithField("module", "lexer"))
		l.Lex(cmd)

		p := parser.New(l, logger.WithField("module", "parser"))
		queries, err := p.Parse()
		if err != nil {
			logger.Error(err)
			continue
		}

		for _, q := range queries {
			if err := e.Process(q); err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}
