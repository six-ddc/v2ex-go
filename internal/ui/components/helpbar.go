package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-go/internal/ui"
)

// HelpItem 帮助项
type HelpItem struct {
	Key  string
	Desc string
}

// HelpBar 底部快捷键提示栏组件
type HelpBar struct {
	items []HelpItem
	width int
}

// NewHelpBar 创建帮助栏
func NewHelpBar() HelpBar {
	return HelpBar{}
}

// SetItems 设置帮助项
func (h *HelpBar) SetItems(items []HelpItem) {
	h.items = items
}

// SetWidth 设置宽度
func (h *HelpBar) SetWidth(width int) {
	h.width = width
}

// View 渲染帮助栏
func (h HelpBar) View() string {
	keyStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Info).
		Background(ui.CurrentTheme.HeaderBg).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Muted).
		Background(ui.CurrentTheme.HeaderBg)

	sepStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Border).
		Background(ui.CurrentTheme.HeaderBg)

	var items []string
	for _, item := range h.items {
		items = append(items, keyStyle.Render(item.Key)+descStyle.Render(" "+item.Desc))
	}

	content := strings.Join(items, sepStyle.Render("  │  "))

	// 填充到整行宽度
	contentWidth := lipgloss.Width(content)
	if contentWidth < h.width {
		padding := lipgloss.NewStyle().Background(ui.CurrentTheme.HeaderBg).Render(strings.Repeat(" ", h.width-contentWidth))
		content = content + padding
	}

	return content
}

// MainViewHelp 主视图帮助项
var MainViewHelp = []HelpItem{
	{Key: "Tab", Desc: "切换焦点"},
	{Key: "h/l", Desc: "切换Tab"},
	{Key: "j/k", Desc: "上下"},
	{Key: "Enter", Desc: "打开"},
	{Key: "t", Desc: "主题"},
	{Key: "q", Desc: "退出"},
}

// DetailViewHelp 详情视图帮助项
var DetailViewHelp = []HelpItem{
	{Key: "j/k", Desc: "滚动"},
	{Key: "g/G", Desc: "顶部/底部"},
	{Key: "n", Desc: "加载更多"},
	{Key: "[/]", Desc: "上/下篇"},
	{Key: "o", Desc: "浏览器"},
	{Key: "q", Desc: "返回"},
}

// SearchViewHelp 搜索视图帮助项
var SearchViewHelp = []HelpItem{
	{Key: "Enter", Desc: "选择"},
	{Key: "Esc", Desc: "取消"},
	{Key: "↑/↓", Desc: "移动"},
	{Key: "Ctrl+N/P", Desc: "下/上一个"},
}
