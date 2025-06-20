package parser

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	file, _ := os.ReadFile("/Users/wushaojie/Downloads/MDNote.md")
	MdDoc().Parse(string(file))
}
