package parser

import (
	"errors"

	"github.com/kotsmile/jql/internal/lexer/token"
	"github.com/kotsmile/jql/internal/logger"
)

type tokenInterator interface {
	Next() (*token.Token, error)
}

type parser struct {
	tokens tokenInterator
	logger logger.Logger
}

func New(i tokenInterator, logger logger.Logger) *parser {
	return &parser{
		tokens: i,
		logger: logger,
	}
}

func (p *parser) Parse() ([]*astNode, error) {
	var queries []*astNode

	var tokens []token.Token
	for {
		t, err := p.tokens.Next()
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}

		if t.Is(token.Semicolon) {
			node, err := p.parseNode(tokens)
			if err != nil {
				return nil, err
			}

			queries = append(queries, node)
			tokens = make([]token.Token, 0)
			continue
		}
		tokens = append(tokens, *t)
	}

	for i, q := range queries {
		p.logger.Debugf("Query: %d\n"+q.String()+"\n", i+1)
	}

	return queries, nil
}

func (p *parser) parseNode(tokens []token.Token) (*astNode, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	var root *astNode = &astNode{}

	cmdToken := tokens[0]
	if cmdToken.Is(token.Word) {
		switch cmdToken.Value() {
		case LoadKeyword.String():
			if len(tokens) < 2 {
				return nil, errors.New("missing file name for load command")
			}

			filenameToken := tokens[1]
			if filenameToken.Is(token.String) {
				root.value = NewKeyword(LoadKeyword)
				root.AppendChild(NewAstNode(StringNode(filenameToken.Value())))

				tokens = tokens[2:]
			}

			if len(tokens) == 0 {
				goto skip
			}

			asToken := tokens[0]
			if asToken.Is(token.Word) && asToken.Value() == AsKeyword.String() {
				if len(tokens) < 2 {
					return nil, errors.New("missing table name for load command")
				}

				tablenameToken := tokens[1]
				if tablenameToken.Is(token.Word) || tablenameToken.Is(token.String) {
					node := NewAstNode(NewKeyword(AsKeyword))
					node.AppendChild(NewAstNode(StringNode(tablenameToken.Value())))
					root.AppendChild(node)
					tokens = tokens[2:]
				}
			}

		default:
			return nil, errors.New("unknown keyword")
		}
	} else {
		return nil, errors.New("unexpected token")
	}

skip:
	return root, nil
}
