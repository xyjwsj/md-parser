package parser

import "strings"

// 移动到下一个标记
func (p *Parser) nextToken() {
	p.current = p.next
	p.next = p.lexer.NextToken()
}

func (p *Parser) resetPreviousToken() {

}

// Parse 解析Markdown文档
func (p *Parser) Parse() *Node {
	root := &Node{Type: TokenText}

	nextToken := true

	for p.current.Type != TokenEOF {
		switch p.current.Type {
		case TokenHeader:
			root.Children = append(root.Children, p.parseHeader())
		case TokenParagraph:
			root.Children = append(root.Children, p.parseParagraph())
		case TokenListItem:
			root.Children = append(root.Children, p.parseList())
			nextToken = false
		case TokenCodeBlock:
			root.Children = append(root.Children, p.parseCodeBlock())
		case TokenHorizontalRule:
			root.Children = append(root.Children, p.parseHorizontalRule())
		case TokenTable:
			root.Children = append(root.Children, p.parseTable())
		}

		if nextToken {
			p.nextToken()
		} else {
			nextToken = true
		}
	}

	return root
}

// 解析标题
func (p *Parser) parseHeader() *Node {
	return &Node{
		Type:     TokenHeader,
		Level:    p.current.Level,
		Content:  p.current.Content,
		Children: p.parseInline(p.current.Content),
	}
}

// 解析段落
func (p *Parser) parseParagraph() *Node {
	paragraph := &Node{
		Type:    TokenParagraph,
		Content: p.current.Content,
	}

	// 解析段落内的内联元素
	paragraph.Children = p.parseInline(p.current.Content)

	return paragraph
}

// 解析列表
func (p *Parser) parseList() *Node {
	// 创建主列表根节点
	list := &Node{
		Type:   TokenList,
		Indent: p.current.Indent,
	}

	// 使用栈保存当前可用的 TokenList 节点
	stack := []*Node{list}

	for p.current.Type == TokenListItem {
		item := &Node{
			Type:    TokenListItem,
			Indent:  p.current.Indent,
			Content: p.current.Content,
		}
		item.Children = p.parseInline(p.current.Content)

		// 找到合适的父 TokenList（栈中最后一个缩进小于当前项的）
		var parentList *Node
		for i := len(stack) - 1; i >= 0; i-- {
			if stack[i].Indent < item.Indent {
				parentList = stack[i]
				break
			}
		}

		// 如果没找到合适的父节点，默认使用根列表
		if parentList == nil {
			parentList = list
		}

		// 添加到父 TokenList 下
		parentList.Children = append(parentList.Children, item)

		// 如果下一个项缩进更深，说明属于当前项的子列表，需创建新 TokenList
		if p.next.Type == TokenListItem && p.next.Indent > item.Indent {
			subList := &Node{
				Type:   TokenList,
				Indent: item.Indent + 1,
			}
			item.Children = append(item.Children, subList) // 子列表作为当前项的子节点
			stack = append(stack, subList)
		}

		p.nextToken()
	}

	return list
}

// 解析代码块
func (p *Parser) parseCodeBlock() *Node {
	return &Node{
		Type:    TokenCodeBlock,
		Content: p.current.Content,
	}
}

// 解析水平线
func (p *Parser) parseHorizontalRule() *Node {
	return &Node{
		Type: TokenHorizontalRule,
	}
}

// 解析内联元素
func (p *Parser) parseInline(content string) []*Node {
	var nodes []*Node
	var currentText strings.Builder

	for i := 0; i < len(content); i++ {
		// 检查强调
		if i+1 < len(content) && content[i] == '*' && content[i+1] == '*' {
			// 保存当前文本
			if currentText.Len() > 0 {
				nodes = append(nodes, &Node{
					Type:    TokenText,
					Content: currentText.String(),
				})
				currentText.Reset()
			}

			// 查找结束标记
			end := i + 2
			for end < len(content)-1 && !(content[end] == '*' && content[end+1] == '*') {
				end++
			}

			if end < len(content)-1 {
				// 找到结束标记
				strongContent := content[i+2 : end]
				strongNode := &Node{
					Type:    TokenStrong,
					Content: strongContent,
				}
				strongNode.Children = p.parseInline(strongContent)
				nodes = append(nodes, strongNode)
				i = end + 1 // 跳过结束标记
			} else {
				// 没有找到结束标记，当作普通文本处理
				currentText.WriteByte('*')
				currentText.WriteByte('*')
			}
			continue
		}

		// 检查强调（单星号）
		if content[i] == '*' {
			// 保存当前文本
			if currentText.Len() > 0 {
				nodes = append(nodes, &Node{
					Type:    TokenText,
					Content: currentText.String(),
				})
				currentText.Reset()
			}

			// 查找结束标记
			end := i + 1
			for end < len(content) && content[end] != '*' {
				end++
			}

			if end < len(content) {
				// 找到结束标记
				emphasisContent := content[i+1 : end]
				emphasisNode := &Node{
					Type:    TokenEmphasis,
					Content: emphasisContent,
				}
				emphasisNode.Children = p.parseInline(emphasisContent)
				nodes = append(nodes, emphasisNode)
				i = end // 跳过结束标记
			} else {
				// 没有找到结束标记，当作普通文本处理
				currentText.WriteByte('*')
			}
			continue
		}

		// 检查链接
		if content[i] == '[' {
			// 保存当前文本
			if currentText.Len() > 0 {
				nodes = append(nodes, &Node{
					Type:    TokenText,
					Content: currentText.String(),
				})
				currentText.Reset()
			}

			// 查找链接文本结束标记
			textEnd := i + 1
			for textEnd < len(content) && content[textEnd] != ']' {
				textEnd++
			}

			if textEnd < len(content) && textEnd+2 < len(content) && content[textEnd+1] == '(' {
				// 找到链接文本结束标记，继续查找URL结束标记
				urlEnd := textEnd + 2
				for urlEnd < len(content) && content[urlEnd] != ')' {
					urlEnd++
				}

				if urlEnd < len(content) {
					// 找到完整的链接
					linkText := content[i+1 : textEnd]
					linkURL := content[textEnd+2 : urlEnd]
					linkNode := &Node{
						Type:    TokenLink,
						Content: linkText,
						Link:    linkURL,
					}
					linkNode.Children = p.parseInline(linkText)
					nodes = append(nodes, linkNode)
					i = urlEnd // 跳过整个链接
					continue
				}
			}
		}

		// 检查图片
		if i+1 < len(content) && content[i] == '!' && content[i+1] == '[' {
			// 保存当前文本
			if currentText.Len() > 0 {
				nodes = append(nodes, &Node{
					Type:    TokenText,
					Content: currentText.String(),
				})
				currentText.Reset()
			}

			// 查找图片描述结束标记
			altEnd := i + 2
			for altEnd < len(content) && content[altEnd] != ']' {
				altEnd++
			}

			if altEnd < len(content) && altEnd+2 < len(content) && content[altEnd+1] == '(' {
				// 找到图片描述结束标记，继续查找URL结束标记
				urlEnd := altEnd + 2
				for urlEnd < len(content) && content[urlEnd] != ')' {
					urlEnd++
				}

				if urlEnd < len(content) {
					// 找到完整的图片
					altText := content[i+2 : altEnd]
					imgURL := content[altEnd+2 : urlEnd]
					imgNode := &Node{
						Type:    TokenImage,
						Content: altText,
						Link:    imgURL,
					}
					nodes = append(nodes, imgNode)
					i = urlEnd // 跳过整个图片
					continue
				}
			}
		}

		// 普通字符
		currentText.WriteByte(content[i])
	}

	// 添加剩余的文本
	if currentText.Len() > 0 {
		nodes = append(nodes, &Node{
			Type:    TokenText,
			Content: currentText.String(),
		})
	}

	return nodes
}

func (p *Parser) parseTable() *Node {
	tableContent := p.current.Content
	lines := strings.Split(tableContent, "\n")

	if len(lines) < 2 {
		return nil
	}

	headers := strings.Split(strings.Trim(lines[0], "|"), "|")
	rows := lines[2:]

	// 验证分隔行是否是有效的表格分隔符
	validSeparator := true
	for _, c := range lines[1] {
		if c != '-' && c != ' ' && c != ':' && c != '|' {
			validSeparator = false
			break
		}
	}

	if !validSeparator || len(rows) == 0 {
		return nil
	}

	table := &Node{
		Type:    TokenTable,
		Content: tableContent,
	}

	// 添加表头
	headerRow := &Node{Type: TokenTableRow}
	for _, header := range headers {
		cell := &Node{Type: TokenTableCell, Content: strings.TrimSpace(header)}
		headerRow.Children = append(headerRow.Children, cell)
	}
	table.Children = append(table.Children, headerRow)

	// 添加数据行
	for _, row := range rows {
		cells := strings.Split(strings.Trim(row, "|"), "|")
		dataRow := &Node{Type: TokenTableRow}
		for _, cell := range cells {
			cellNode := &Node{Type: TokenTableCell, Content: strings.TrimSpace(cell)}
			dataRow.Children = append(dataRow.Children, cellNode)
		}
		table.Children = append(table.Children, dataRow)
	}

	return table
}

// NewParser 创建一个新的语法分析器
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	// 读取两个标记，初始化current和next
	p.current = p.lexer.NextToken()
	p.next = p.lexer.NextToken()
	return p
}
