package parser

import "strings"

// 移动到下一个标记
func (p *Parser) nextToken() {
	p.current = p.next
	p.next = p.lexer.NextToken()
}

// Parse 解析Markdown文档
func (p *Parser) Parse() *Node {
	root := &Node{Type: TokenText}

	for p.current.Type != TokenEOF {
		switch p.current.Type {
		case TokenHeader:
			root.Children = append(root.Children, p.parseHeader())
		case TokenParagraph:
			root.Children = append(root.Children, p.parseParagraph())
		case TokenListItem:
			root.Children = append(root.Children, p.parseList())
		case TokenCodeBlock:
			root.Children = append(root.Children, p.parseCodeBlock())
		case TokenHorizontalRule:
			root.Children = append(root.Children, p.parseHorizontalRule())
		}
		p.nextToken()
	}

	return root
}

// 解析标题
func (p *Parser) parseHeader() *Node {
	return &Node{
		Type:    TokenHeader,
		Level:   p.current.Level,
		Content: p.current.Content,
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
	list := &Node{Type: TokenList}

	// 使用栈保存当前层级的父节点（每个节点代表一个 List 或 ListItem）
	stack := []*Node{list}

	for p.current.Type == TokenListItem {
		currentItem := &Node{
			Type:    TokenListItem,
			Indent:  p.current.Indent,
			Content: p.current.Content,
		}
		currentItem.Children = p.parseInline(p.current.Content)

		// 找到合适的父节点（即缩进小于当前项的最后一个节点）
		var parentNode *Node
		for i := len(stack) - 1; i >= 0; i-- {
			if stack[i].Indent < currentItem.Indent {
				parentNode = stack[i]
				break
			}
		}

		// 如果没找到，说明当前项是顶级列表项
		if parentNode == nil {
			parentNode = list
		}

		// 添加当前项到父节点下
		parentNode.Children = append(parentNode.Children, currentItem)

		// 如果下一个项缩进更深，则将当前项压入栈中作为新层级起点
		if p.next.Type == TokenListItem && p.next.Indent > currentItem.Indent {
			stack = append(stack, currentItem)
		} else {
			// 否则，清理栈中不再适用的层级
			for len(stack) > 0 && stack[len(stack)-1].Indent >= currentItem.Indent {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, currentItem)
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

// NewParser 创建一个新的语法分析器
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	// 读取两个标记，初始化current和next
	p.current = p.lexer.NextToken()
	p.next = p.lexer.NextToken()
	return p
}
