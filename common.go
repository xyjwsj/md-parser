package parser

// TokenType 表示标记类型
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenHeader
	TokenParagraph
	TokenList
	TokenListItem
	TokenCodeBlock
	TokenHorizontalRule
	TokenEmphasis
	TokenStrong
	TokenLink
	TokenImage
	TokenTable
	TokenTableRow
	TokenTableCell
	TokenText
)

// Token 表示一个Markdown标记
type Token struct {
	Type    TokenType
	Content string
	Level   int    // 用于标题和列表
	Indent  int    // 用于列表项缩进
	Link    string // 用于链接和图片
	Alt     string // 用于图片
}

// Lexer 词法分析器
type Lexer struct {
	input        string
	position     int  // 当前字符位置
	readPosition int  // 下一个字符位置
	ch           byte // 当前字符
}

// Node 表示AST中的一个节点
type Node struct {
	Type     TokenType
	Content  string
	Level    int
	Indent   int
	Link     string
	Alt      string
	Children []*Node
}

// Parser 语法分析器
type Parser struct {
	lexer   *Lexer
	current Token
	next    Token
}

// Renderer HTML渲染器
type Renderer struct{}
