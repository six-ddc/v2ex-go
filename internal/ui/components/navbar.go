package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-tui/internal/model"
	"github.com/six-ddc/v2ex-tui/internal/ui"
)

// Navbar 一级导航栏组件
type Navbar struct {
	tabs    []model.Tab
	cursor  int  // 光标位置（浏览时移动）
	active  int  // 激活项（Enter 确认后的实际选中）
	focused bool
	width   int
}

// NewNavbar 创建导航栏
func NewNavbar() Navbar {
	return Navbar{
		tabs:   model.DefaultTabs,
		cursor: 9, // 默认 "全部"
		active: 9,
	}
}

// Init 初始化
func (n Navbar) Init() tea.Cmd {
	return nil
}

// Update 更新
func (n Navbar) Update(msg tea.Msg) (Navbar, tea.Cmd) {
	return n, nil
}

// SetFocused 设置焦点状态
func (n *Navbar) SetFocused(focused bool) {
	n.focused = focused
	if focused {
		// 获得焦点时，光标恢复到激活项位置
		n.cursor = n.active
	}
}

// SetWidth 设置宽度
func (n *Navbar) SetWidth(width int) {
	n.width = width
}

// SetSelected 设置选中项（同时设置光标和激活项）
func (n *Navbar) SetSelected(index int) {
	if index >= 0 && index < len(n.tabs) {
		n.cursor = index
		n.active = index
	}
}

// Activate 激活当前光标位置的项
func (n *Navbar) Activate() {
	n.active = n.cursor
}

// Selected 获取当前光标位置
func (n Navbar) Selected() int {
	return n.cursor
}

// SelectedTab 获取当前光标位置的 Tab
func (n Navbar) SelectedTab() model.Tab {
	if n.cursor >= 0 && n.cursor < len(n.tabs) {
		return n.tabs[n.cursor]
	}
	return model.Tab{}
}

// MoveLeft 向左移动光标
func (n *Navbar) MoveLeft() {
	if n.cursor > 0 {
		n.cursor--
	}
}

// MoveRight 向右移动光标
func (n *Navbar) MoveRight() {
	if n.cursor < len(n.tabs)-1 {
		n.cursor++
	}
}

// JumpTo 跳转到指定索引（同时设置光标和激活项）
func (n *Navbar) JumpTo(index int) {
	if index >= 0 && index < len(n.tabs) {
		n.cursor = index
		n.active = index
	}
}

// View 渲染导航栏
func (n Navbar) View() string {
	var items []string

	for i, tab := range n.tabs {
		var rendered string
		isCursor := (i == n.cursor)
		isActive := (i == n.active)

		if n.focused && isCursor {
			// 有焦点时，光标位置高亮背景
			style := lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.PrimaryFg).
				Background(ui.CurrentTheme.Primary).
				Bold(true).
				Padding(0, 1)
			rendered = style.Render(tab.Name)
		} else if !n.focused && isActive {
			// 无焦点时，激活项用颜色标识
			style := lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.Primary).
				Bold(true).
				Padding(0, 1)
			rendered = style.Render(tab.Name)
		} else {
			// 其他项 - 灰色文字
			style := lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.Muted).
				Padding(0, 1)
			rendered = style.Render(tab.Name)
		}
		items = append(items, rendered)
	}

	content := strings.Join(items, "")

	// 添加分隔线
	var separator string
	if n.focused {
		separator = lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Primary).
			Render(strings.Repeat("─", n.width))
	} else {
		separator = lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Border).
			Render(strings.Repeat("─", n.width))
	}

	return content + "\n" + separator
}
