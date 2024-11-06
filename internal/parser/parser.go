package parser

import (
	"errors"

	"github.com/kotsmile/jql/internal/lexer/token"
	"github.com/kotsmile/jql/util"
)

var (
	ErrUnknownKeyword                = errors.New("unknown keyword")
	ErrUnexpectedToken               = errors.New("unexpected token")
	ErrMissingFileNameLoadCommand    = errors.New("missing file name for 'load' command")
	ErrMissingTableNameLoadCommand   = errors.New("missing table name for 'load' command")
	ErrMissingFromKeyword            = errors.New("'select' command: missing 'from' keyword")
	ErrMissingTableNameSelectCommand = errors.New("'select' command: missing table name")
)

type tokenInterator interface {
	Next() (*token.Token, error)
}

type parser struct {
	tokens tokenInterator
	logger util.Logger
}

func New(i tokenInterator, logger util.Logger) *parser {
	return &parser{
		tokens: i,
		logger: logger,
	}
}

func (p *parser) Parse() ([]*AstNode, error) {
	var queries []*AstNode

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

func (p *parser) parseNode(tokens []token.Token) (*AstNode, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	var root *AstNode = &AstNode{}

	cmdToken := tokens[0]
	if cmdToken.Is(token.Word) {
		switch cmdToken.Value() {
		case LoadKeyword.String():
			if len(tokens) < 2 {
				return nil, ErrMissingFileNameLoadCommand
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
					return nil, ErrMissingTableNameLoadCommand
				}

				tablenameToken := tokens[1]
				if tablenameToken.Is(token.Word) || tablenameToken.Is(token.String) {
					node := NewAstNode(NewKeyword(AsKeyword))
					node.AppendChild(NewAstNode(StringNode(tablenameToken.Value())))
					root.AppendChild(node)
					tokens = tokens[2:]
				}
			}
		case TablesKeyword.String():
			root.value = NewKeyword(TablesKeyword)
			return root, nil
		case SelectKeyword.String():
			root.value = NewKeyword(SelectKeyword)
			if len(tokens) < 2 {
				return nil, ErrMissingFromKeyword
			}

			foundFromKeyword := false
			index := 1
			for _, t := range tokens[1:] {
				if t.Is(token.Word) && t.Value() == FromKeyword.String() {
					foundFromKeyword = true
					break
				}

				if t.Value() != "," {
					root.AppendChild(NewAstNode(StringNode(t.Value())))
				}
				index++
			}
			if !foundFromKeyword {
				return nil, ErrMissingFromKeyword
			}

			fromKeyword := NewAstNode(NewKeyword(FromKeyword))
			root.AppendChild(fromKeyword)

			if len(tokens) < index+1 {
				return nil, ErrMissingTableNameSelectCommand
			}
			for _, t := range tokens[index+1:] {
				fromKeyword.AppendChild(NewAstNode(StringNode(t.Value())))
			}

			return root, nil

		default:
			return nil, ErrUnknownKeyword
		}
	} else {
		return nil, ErrUnexpectedToken
	}

skip:
	return root, nil
}
