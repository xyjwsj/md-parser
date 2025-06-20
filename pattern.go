package parser

import (
	"log"
	"regexp"
	"strings"
)

type callFun func()

type LinePattern struct {
	Pattern string
	Type    NodeType
}

var lineFlags = []LinePattern{
	{
		Pattern: `^#{0,6} (.*)`,
		Type:    Heading,
	},
	{
		Pattern: `^\> (.*)`,
		Type:    References,
	},
	{
		Pattern: `^[-+*] (.*)`,
		Type:    UnOrderList,
	},
	{
		Pattern: `^\d. (.*)`,
		Type:    OrderList,
	},
	{
		Pattern: `^\|(.*)\|$`,
		Type:    Table,
	},
	{
		Pattern: "^\\x60{3}",
		Type:    BlockCode,
	},
	{
		Pattern: `^[\*{3}\-{3}_{3}] (.*)`,
		Type:    Divider,
	},
}

var headingStylePattern = `(#{0,6}) `

func ParseNode(lineStr string) Node {
	level := 0
	length := getLeadingWhitespaceLength(lineStr)
	if length > 0 {
		lineStr = lineStr[length:]
		level = length
	}
	for _, item := range lineFlags {
		groups := pattern(item.Pattern, lineStr)
		if len(groups) > 1 {
			rank := 0
			if item.Type == Heading {
				strArr := pattern(headingStylePattern, lineStr)
				s := strArr[len(strArr)-1]
				rank = strings.Count(s, "#")
			}
			log.Println(level, rank, item.Type, groups[len(groups)-1])
			matchItem(groups[len(groups)-1])
			break
		} else {

		}
	}
	return Node{}
}

func pattern(p, str string) []string {
	re := regexp.MustCompile(p)
	return re.FindStringSubmatch(str)
}

var levelPattern = `^(\s*)`

func getLeadingWhitespaceLength(s string) int {
	re := regexp.MustCompile(levelPattern)
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		return len(match[1])
	}
	return 0
}

var itemFlags = map[string]StyleType{
	"`([^`]+)":                              LineCode,
	`$$([^]]+)$$\(([^)]+)(\s+"([^"]+)")?\)`: Link,
	`^!\[([^]]+)\]\(([^)]+)(\s+"[^"]+")?\)`: Image,
	`\*{2}([^*]+)\*{2}`:                     Bold,
	`([*_])((?:[^_*]|\n)+)\1`:               Italic,
	`\~{2}([^~]+)\~{2}`:                     Strikethrough,
}

func matchItem(str string) []Item {
	data := make(map[StyleType][]string, 30)
	for key, itm := range itemFlags {
		strArr := pattern(key, str)
		if len(strArr) > 1 {
			data[itm] = strArr[1:]
		}
	}

	log.Println(data)

	return nil
}
