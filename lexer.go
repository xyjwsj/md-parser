package parser

import "strings"

// NewLexer 创建一个新的词法分析器
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 读取下一个字符
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// 读取直到遇到换行符
func (l *Lexer) readUntilNewline() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 读取直到遇到指定字符
func (l *Lexer) readUntilChar(c byte) string {
	position := l.position
	for l.ch != c && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipEnter() {
	for l.ch == '\n' {
		l.readChar()
	}
}

// 跳过空白字符
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' {
		l.readChar()
	}
}

// 计算行首空格数
func (l *Lexer) countLeadingSpaces() int {
	position := l.position
	ch := l.ch
	count := 0
	for l.ch == ' ' || l.ch == '\t' {
		if l.ch == ' ' {
			count++
		} else {
			count += 4 // 假设一个制表符等于4个空格
		}
		l.readChar()
	}
	l.position = position
	l.ch = ch
	l.readPosition = position + 1
	return count
}

// 读取下一个标记
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipEnter()
	indent := l.countLeadingSpaces()
	l.skipWhitespace()

	if l.ch == 0 {
		return Token{Type: TokenEOF}
	}

	// 检查是否为代码块
	if l.ch == '`' && l.peekChar() == '`' && l.peekChar2() == '`' {
		l.readChar() // 消耗第一个 `
		l.readChar() // 消耗第二个 `
		l.readChar() // 消耗第三个 `
		tok.Type = TokenCodeBlock
		tok.Content = l.readCodeBlock()
		return tok
	}

	// 检查是否为标题
	if l.ch == '#' {
		tok.Type = TokenHeader
		level := 0
		for l.ch == '#' {
			level++
			l.readChar()
		}
		// 跳过标题后的空格
		for l.ch == ' ' {
			l.readChar()
		}
		tok.Level = level
		tok.Content = l.readUntilNewline()
		return tok
	}

	// 检查是否为列表
	if (l.ch == '*' || l.ch == '-' || l.ch == '+') && l.peekChar() == ' ' {
		//bullet := l.ch
		l.readChar() // 消耗列表标记
		l.readChar() // 消耗空格
		tType := TokenListItem
		tok.Type = tType
		tok.Indent = indent
		tok.Content = l.readUntilNewline()
		return tok
	}

	// 检查是否为数字列表
	if IsDigit(l.ch) && l.peekChar() == '.' && l.peekChar2() == ' ' {
		//indent := l.countLeadingSpaces()
		for IsDigit(l.ch) {
			l.readChar()
		}
		l.readChar() // 消耗 .
		l.readChar() // 消耗空格
		tType := TokenListItem
		tok.Type = tType
		tok.Indent = indent
		tok.Content = l.readUntilNewline()
		return tok
	}

	// 检查是否为水平线
	if l.ch == '-' && l.peekChar() == '-' && l.peekChar2() == '-' {
		l.readChar() // 消耗第一个 -
		l.readChar() // 消耗第二个 -
		l.readChar() // 消耗第三个 -
		for l.ch == '-' || l.ch == ' ' {
			l.readChar()
		}
		tok.Type = TokenHorizontalRule
		return tok
	}

	// 处理段落
	tok.Type = TokenParagraph
	tok.Content = l.readParagraph()
	tok.Indent = indent
	return tok
}

// 读取代码块
func (l *Lexer) readCodeBlock() string {
	var result strings.Builder
	for {
		line := l.readUntilNewline()
		if strings.HasPrefix(line, "```") {
			break
		}
		result.WriteString(line)
		result.WriteByte('\n')
		l.readChar() // 消耗换行符
	}
	return strings.TrimSuffix(result.String(), "\n")
}

// 读取段落
func (l *Lexer) readParagraph() string {
	var result strings.Builder
	result.WriteString(l.readUntilNewline())
	position := l.position

	// 检查下一行是否继续当前段落
	l.skipWhitespace()
	if l.ch != 0 && l.ch != '#' && l.ch != '*' && l.ch != '-' && l.ch != '+' && !IsDigit(l.ch) {
		l.position = position
		l.readChar() // 消耗换行符
		result.WriteByte('\n')
		result.WriteString(l.readParagraph())
	} else {
		//l.position = position
	}

	return result.String()
}

// 查看下一个字符
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// 查看下两个字符
func (l *Lexer) peekChar2() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}
