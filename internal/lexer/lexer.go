package lexer

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/kotsmile/jql/internal/lexer/token"
	"github.com/kotsmile/jql/util"
)

var (
	separators = []rune{' ', '\t', '\n', '\r'}
	symbols    = []rune{
		/* parentheses */ '(', ')', '[', ']', '{', '}',
		/* arithmetic */ '+', '-', '/', '%', '*',
		/* punctuation */ ',', ':', ';', '\'', '"', '.',
		/* logical */ '|', '&', '?', '!',
		/* comparison */ '<', '>', '=',
	}
)

type lexer struct {
	index   int
	cmd     string
	cmdLeft string
	logger  util.Logger
}

func New(logger util.Logger) *lexer {
	return &lexer{
		index:  0,
		logger: logger,
	}
}

func (p *lexer) Lex(cmd string) {
	p.cmd = cmd
	p.cmdLeft = cmd
}

func (p *lexer) Next() (t *token.Token, err error) {
	word, rest := nextWord(p.cmdLeft, separators, symbols)
	p.cmdLeft = rest

	if word == "" {
		t = nil
		err = nil
	} else if word == ";" {
		t = token.New(token.Semicolon, word)
	} else if word[0] == '"' {
		index := strings.Index(rest, "\"")
		if index == -1 {
			return nil, errors.New("unterminated string")
		}

		content := word[1:] + rest[:index]
		p.cmdLeft = rest[index+1:]

		t = token.New(token.String, content)
		err = nil
	} else {
		t = token.New(token.Word, word)
		err = nil
	}

	if t != nil {
		p.logger.WithFields(util.LoggerFields{
			"type":  t.Type(),
			"value": t.Value(),
		}).Debug("token")
	}

	return
}

func (p *lexer) Peek() (*token.Token, error) {
	cmdLeft := p.cmdLeft
	t, err := p.Next()
	p.cmdLeft = cmdLeft

	return t, err
}

func (p *lexer) Collect() ([]token.Token, error) {
	var tokens []token.Token

	for {
		token, err := p.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}

		tokens = append(tokens, *token)
	}

	return tokens, nil
}

func nextWord(input string, whitespaces []rune, special []rune) (string, string) {
	whitespacesMap := make(map[rune]struct{})
	for _, w := range whitespaces {
		whitespacesMap[w] = struct{}{}
	}
	specialMap := make(map[rune]struct{})
	for _, s := range special {
		specialMap[s] = struct{}{}
	}

	start := 0
	for start < len(input) {
		r, size := utf8.DecodeRuneInString(input[start:])
		if _, ok := whitespacesMap[r]; !ok {
			break
		} else {
			start += size
		}
	}

	if start >= len(input) {
		return "", ""
	}

	end := start
	r, size := utf8.DecodeRuneInString(input[end:])
	if _, ok := specialMap[r]; ok {
		end += size
		return input[start:end], input[end:]
	}

	for end < len(input) {
		r, size := utf8.DecodeRuneInString(input[end:])
		if _, ok := specialMap[r]; ok {
			return input[start:end], input[end:]
		}
		if _, ok := whitespacesMap[r]; ok {
			return input[start:end], input[end:]
		}

		end += size
	}

	return input[start:end], input[end:]
}
