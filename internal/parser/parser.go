package parser

import (
	"errors"

	"github.com/kotsmile/jql/internal/lexer/token"
	"github.com/kotsmile/jql/util"
)

var (
	ErrUnknownKeyword                 = errors.New("unknown keyword")
	ErrUnexpectedToken                = errors.New("unexpected token")
	ErrMissingFileNameLoadCommand     = errors.New("missing file name for 'load' command")
	ErrMissingTableNameLoadCommand    = errors.New("missing table name for 'load' command")
	ErrMissingFromKeyword             = errors.New("'select' command: missing 'from' keyword")
	ErrMissingTableNameSelectCommand  = errors.New("'select' command: missing table name")
	ErrMissingColumnNameSelectCommand = errors.New("'select' command: missing column name")
	ErrEmptyCommand                   = errors.New("empty command")
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
	var root *AstNode = &AstNode{}

	cmdToken, ok := util.Next(&tokens)
	if !ok {
		return nil, ErrEmptyCommand
	}

	if cmdToken.Is(token.Word) {
		switch cmdToken.Value() {
		case LoadKeyword.String():
			filenameToken, ok := util.Next(&tokens)
			if !ok {
				return nil, ErrMissingFileNameLoadCommand
			}

			if filenameToken.Is(token.String) {
				root.value = NewKeyword(LoadKeyword)
				root.AppendChild(NewAstNode(StringNode(filenameToken.Value())))
			}

			asToken, ok := util.Next(&tokens)
			if !ok {
				goto skip
			}

			if asToken.Is(token.Word) && asToken.Value() == AsKeyword.String() {
				tablenameToken, ok := util.Next(&tokens)
				if !ok {
					return nil, ErrMissingTableNameLoadCommand
				}

				if tablenameToken.Is(token.Word) || tablenameToken.Is(token.String) {
					node := NewAstNode(NewKeyword(AsKeyword))
					node.AppendChild(NewAstNode(StringNode(tablenameToken.Value())))
					root.AppendChild(node)
				}
			}
		case TablesKeyword.String():
			root.value = NewKeyword(TablesKeyword)
			return root, nil
		case SelectKeyword.String():
			root.value = NewKeyword(SelectKeyword)
			for len(tokens) > 0 {
				t, ok := util.Next(&tokens)
				if !ok {
					return nil, ErrMissingFromKeyword
				}
				if t.Is(token.Word) && t.Value() == FromKeyword.String() {
					fromNode := NewAstNode(NewKeyword(FromKeyword))
					root.AppendChild(fromNode)

					tableName, ok := util.Next(&tokens)
					if !ok {
						return nil, ErrMissingTableNameSelectCommand
					}
					fromNode.AppendChild(NewAstNode(StringNode(tableName.Value())))
					break
				}

				if t.Is(token.Comma) {
					continue
				}
				if t.Is(token.Word) {
					child := NewAstNode(StringNode(t.Value()))
					root.AppendChild(child)

					possibleAsToken, ok := util.Peek(tokens)
					if !ok {
						return nil, ErrMissingFromKeyword
					}
					if !possibleAsToken.Is(token.Word) {
						continue
					}
					if possibleAsToken.Value() != AsKeyword.String() {
						continue
					}

					util.Next(&tokens)

					asNode := NewAstNode(NewKeyword(AsKeyword))
					child.AppendChild(asNode)

					columnNameToken, ok := util.Next(&tokens)
					if !ok {
						return nil, ErrMissingColumnNameSelectCommand
					}
					if !columnNameToken.Is(token.Word) {
						continue
					}

					tableNameNode := NewAstNode(StringNode(columnNameToken.Value()))
					asNode.AppendChild(tableNameNode)
				}
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
