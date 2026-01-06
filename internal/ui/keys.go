package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap 快捷键映射
type KeyMap struct {
	// 全局
	Quit        key.Binding
	Help        key.Binding
	Search      key.Binding
	Refresh     key.Binding
	Tab         key.Binding
	ShiftTab    key.Binding
	Enter       key.Binding
	Escape      key.Binding
	OpenBrowser key.Binding
	ToggleTheme key.Binding

	// 导航
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding

	// 列表导航
	Top          key.Binding
	Bottom       key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	PageUp       key.Binding
	PageDown     key.Binding

	// 详情视图
	PrevTopic  key.Binding
	NextTopic  key.Binding
	NextReply  key.Binding
	PrevReply  key.Binding
	CopyLink   key.Binding

	// 数字快捷键
	Num1 key.Binding
	Num2 key.Binding
	Num3 key.Binding
	Num4 key.Binding
	Num5 key.Binding
	Num6 key.Binding
	Num7 key.Binding
	Num8 key.Binding
	Num9 key.Binding
	Num0 key.Binding
}

// DefaultKeyMap 默认快捷键映射
var DefaultKeyMap = KeyMap{
	// 全局
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "退出"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "帮助"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "搜索"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "刷新"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "切换焦点"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("Shift+Tab", "反向切换"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "确认"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "返回"),
	),
	OpenBrowser: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "浏览器打开"),
	),
	ToggleTheme: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "切换主题"),
	),

	// 导航
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/←", "左"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/→", "右"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "上"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "下"),
	),

	// 列表导航
	Top: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "顶部"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "底部"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("Ctrl+U", "上翻半页"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("Ctrl+D", "下翻半页"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+b"),
		key.WithHelp("Ctrl+B", "上翻页"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+f", " "),
		key.WithHelp("Ctrl+F/Space", "下翻页"),
	),

	// 详情视图
	PrevTopic: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "上一篇"),
	),
	NextTopic: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "下一篇"),
	),
	NextReply: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "下条回复"),
	),
	PrevReply: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "上条回复"),
	),
	CopyLink: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "复制链接"),
	),

	// 数字快捷键
	Num1: key.NewBinding(key.WithKeys("1")),
	Num2: key.NewBinding(key.WithKeys("2")),
	Num3: key.NewBinding(key.WithKeys("3")),
	Num4: key.NewBinding(key.WithKeys("4")),
	Num5: key.NewBinding(key.WithKeys("5")),
	Num6: key.NewBinding(key.WithKeys("6")),
	Num7: key.NewBinding(key.WithKeys("7")),
	Num8: key.NewBinding(key.WithKeys("8")),
	Num9: key.NewBinding(key.WithKeys("9")),
	Num0: key.NewBinding(key.WithKeys("0")),
}
