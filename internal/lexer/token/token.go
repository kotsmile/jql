package token

import "fmt"

type TokenType string

const (
	Word      TokenType = "word"
	String              = "string"
	Semicolon           = "semicolon"
	// Number
	// Boolean
	// Null
	Comma = "comma"
	// Colon
	// Period
	// Equal
	// Plus
	// Minus
	// Slash
	// Percent
	// LessThan
	// GreaterThan
	// LessThanOrEqual
	// GreaterThanOrEqual
	// NotEqual
	// And
	// Or
	// Not
	// // Parentheses are smooth and curved (),
	// // brackets are square [], and braces are curly {}
	// LeftParenthesis
	// RightParenthesis
	// LeftBracket
	// RightBracket
	// LeftBrace
	// RightBrace
)

type Token struct {
	type_ TokenType
	value string
}

func New(type_ TokenType, value string) *Token {
	return &Token{
		type_: type_,
		value: value,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("{ type: %s, value: '%s' }", t.type_, t.value)
}

func (t Token) Is(type_ TokenType) bool {
	return t.type_ == type_
}

func (t Token) Value() string {
	return t.value
}

func (t Token) Type() TokenType {
	return t.type_
}
