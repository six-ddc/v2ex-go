package model

// Node 节点模型
type Node struct {
	Name     string // 显示名称
	Code     string // URL 编码
	Icon     string // 图标 (emoji)
	Children []Node // 子节点
}

// Tab 一级导航项
type Tab struct {
	Name   string // 显示名称
	Code   string // URL 编码
	HasSub bool   // 是否有二级节点
}

// DefaultTabs 预定义的一级导航
var DefaultTabs = []Tab{
	{Name: "技术", Code: "tech", HasSub: true},
	{Name: "创意", Code: "creative", HasSub: true},
	{Name: "好玩", Code: "play", HasSub: true},
	{Name: "Apple", Code: "apple", HasSub: true},
	{Name: "酷工作", Code: "jobs", HasSub: true},
	{Name: "交易", Code: "deals", HasSub: true},
	{Name: "城市", Code: "city", HasSub: true},
	{Name: "问与答", Code: "qna", HasSub: true},
	{Name: "最热", Code: "hot", HasSub: false},
	{Name: "全部", Code: "all", HasSub: true},
	{Name: "R2", Code: "r2", HasSub: false},
	{Name: "节点", Code: "nodes", HasSub: false},
	{Name: "关注", Code: "members", HasSub: false},
}
