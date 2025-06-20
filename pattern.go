package parser

import (
	"log"
	"regexp"
)

type callFun func()

var flags = map[string]NodeType{
	"(#{0,}) (.*)":   Heading,
	`^\> (.*)`:       References,
	`[-+*] (.*)`:     UnOrderList,
	`\d. (.*)`:       OrderList,
	"`([^`\\r\\n]+)": LineCode,
	"^(?:\\`{3})([a-zA-Z0-9]*)[\\r\\n]+([\\s\\S]*?)[\\r\\n]+\\`{3}": BlockCode,
	`$$([^]]+)$$\(([^)]+)(\s+"([^"]+)")?\)`:                         Link,
	`^!\[([^]]+)\]\(([^)]+)(\s+"[^"]+")?\)`:                         Image,
	"^|(.*)":                                                        Table,
	`[\*{3}\-{3}_{3}] (.*)`:                                         Divider,
	"":                                                              Text,
}

func ParseLine(lineStr string) Node {
	for p, t := range flags {
		log.Println(t)
		strings := pattern(p, lineStr)
		if len(strings) > 1 {

		}
	}
	return Node{}
}

func pattern(p, str string) []string {
	re := regexp.MustCompile(p)
	return re.FindStringSubmatch(str)
}

func font() {
	pattern := `\*{2}[^\*{2}]+\*{2}`
	log.Println(pattern)
}
