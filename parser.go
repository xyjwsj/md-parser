package parser

import "strings"

type LineCall func(line int, info Content)

type MDDocument struct {
	lineItem []Node
	call     LineCall
}

func MdDoc() *MDDocument {
	return &MDDocument{}
}

func (doc *MDDocument) Parse(content string) {
	split := strings.Split(content, "\n")
	for _, line := range split {
		ParseLine(line)
	}
}

func (doc *MDDocument) Walk() {
	if doc.call != nil {
		for line, item := range doc.lineItem {
			doc.callChildNode(line, &item)
		}
	}
}

func (doc *MDDocument) callChildNode(line int, node *Node) {
	if doc.call != nil {
		if node.Nodes == nil {
			doc.call(line, node.Info)
			return
		}
		for _, itm := range node.Nodes {
			doc.callChildNode(line, &itm)
		}
	}
}
