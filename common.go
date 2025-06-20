package parser

type NodeType int

var (
	NodeNone    NodeType = 0
	Heading     NodeType = 1 //标题
	OrderList   NodeType = 2 //有序列表
	UnOrderList NodeType = 3 //无序列表
	Text        NodeType = 4 //普通文本
	Paragraph   NodeType = 5 //段落
	References  NodeType = 6 //引用
	BlockCode   NodeType = 7 //块代码
	Table       NodeType = 8 //表格
	Divider     NodeType = 9 //分割线
)

type StyleType int

var (
	StyleNone StyleType = 0 //普通样式
	LineCode  StyleType = 1 //行内代码
	Link      StyleType = 2 //链接
	Image     StyleType = 3 //图片
	Font      StyleType = 4 //字体样式
)

type FontType int

var (
	FontNone      FontType = 0
	Bold          FontType = 1 // 粗体
	Italic        FontType = 2 // 斜体
	Strikethrough FontType = 3 // 删除线
)

type Item struct {
	Type  StyleType
	Text  string
	Items []Item
}

type Node struct {
	Type  NodeType // 行类型
	Level int      // 层级(结构层面)
	Rank  int      // 等级(样式层面)
	Info  []Item   // 当前行所有样式
}
