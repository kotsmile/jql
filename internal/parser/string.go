package parser

type StringNode string

func (s StringNode) String() string {
	return string(s)
}

func (s StringNode) Value() string {
	return string(s)
}

func (s StringNode) Type() string {
	return "string"
}
