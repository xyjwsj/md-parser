package parser

import (
	"log"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	//file, _ := os.ReadFile("/Users/wushaojie/Downloads/MDNote.md")
	file, _ := os.ReadFile("/Users/wushaojie/Downloads/MDNote-copy.md")
	lexer := NewLexer(string(file))
	parser := NewParser(lexer)
	ast := parser.Parse()
	log.Println(ast)
}
