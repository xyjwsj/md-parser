package parser

type NodeType int

var (
	Heading     NodeType = 1  //标题
	OrderList   NodeType = 2  //有序列表
	UnOrderList NodeType = 3  //无序列表
	Text        NodeType = 5  //普通问题
	Paragraph   NodeType = 6  //段落
	References  NodeType = 7  //引用
	LineCode    NodeType = 8  //行内代码
	BlockCode   NodeType = 9  //块代码
	Link        NodeType = 10 //链接
	Image       NodeType = 11 //图片
	Table       NodeType = 12 //表格
	Divider     NodeType = 13 //分割线
)

type Content struct {
	Type  NodeType
	Level int
	Text  string
}

type Node struct {
	Info  Content
	Nodes []Node
}
