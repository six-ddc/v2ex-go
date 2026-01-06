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
	tabs     []model.Tab
	selected int
	focused  bool
	width    int
}

// NewNavbar 创建导航栏
func NewNavbar() Navbar {
	return Navbar{
		tabs:     model.DefaultTabs,
		selected: 9, // 默认选中 "全部"
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
}

// SetWidth 设置宽度
func (n *Navbar) SetWidth(width int) {
	n.width = width
}

// SetSelected 设置选中项
func (n *Navbar) SetSelected(index int) {
	if index >= 0 && index < len(n.tabs) {
		n.selected = index
	}
}

// Selected 获取当前选中的索引
func (n Navbar) Selected() int {
	return n.selected
}

// SelectedTab 获取当前选中的 Tab
func (n Navbar) SelectedTab() model.Tab {
	if n.selected >= 0 && n.selected < len(n.tabs) {
		return n.tabs[n.selected]
	}
	return model.Tab{}
}

// MoveLeft 向左移动
func (n *Navbar) MoveLeft() {
	if n.selected > 0 {
		n.selected--
	}
}

// MoveRight 向右移动
func (n *Navbar) MoveRight() {
	if n.selected < len(n.tabs)-1 {
		n.selected++
	}
}

// JumpTo 跳转到指定索引
func (n *Navbar) JumpTo(index int) {
	if index >= 0 && index < len(n.tabs) {
		n.selected = index
	}
}

// View 渲染导航栏
func (n Navbar) View() string {
	var items []string

	for i, tab := range n.tabs {
		var rendered string
		if i == n.selected {
			if n.focused {
				// 选中且有焦点 - 高亮背景（语义上的高亮）
				style := lipgloss.NewStyle().
					Foreground(ui.CurrentTheme.PrimaryFg).
					Background(ui.CurrentTheme.Primary).
					Bold(true).
					Padding(0, 1)
				rendered = style.Render(tab.Name)
			} else {
				// 选中但无焦点 - 用 [xxx] 包裹标识当前选中
				style := lipgloss.NewStyle().
					Foreground(ui.CurrentTheme.Primary).
					Bold(true)
				rendered = style.Render("[" + tab.Name + "]")
			}
		} else {
			// 未选中 - 灰色文字，不设置背景色
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
