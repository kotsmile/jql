package parser

type KeywordType string

const (
	LoadKeyword   KeywordType = "load"
	TablesKeyword KeywordType = "tables"
	AsKeyword     KeywordType = "as"
	SelectKeyword KeywordType = "select"
	FromKeyword   KeywordType = "from"
)

func NewKeyword(type_ KeywordType) *KeywordNode {
	return &KeywordNode{
		type_: type_,
	}
}

func (k KeywordType) String() string {
	return string(k)
}

type KeywordNode struct {
	type_ KeywordType
}

func (k *KeywordNode) String() string {
	return k.type_.String()
}

func (k *KeywordNode) Value() string {
	return k.type_.String()
}

func (k *KeywordNode) Type() string {
	return "keyword"
}
