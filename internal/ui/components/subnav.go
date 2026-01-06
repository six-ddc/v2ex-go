package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-tui/internal/model"
	"github.com/six-ddc/v2ex-tui/internal/ui"
)

// Subnav 二级导航组件
type Subnav struct {
	nodes   []model.Node
	cursor  int  // 光标位置（浏览时移动）
	active  int  // 激活项（Enter 确认后的实际选中），-1 表示无激活项
	focused bool
	width   int
	offset  int // 用于横向滚动
}

// NewSubnav 创建二级导航
func NewSubnav() Subnav {
	return Subnav{
		cursor: 0,  // 光标默认第一个
		active: -1, // 默认无激活项
	}
}

// Init 初始化
func (s Subnav) Init() tea.Cmd {
	return nil
}

// Update 更新
func (s Subnav) Update(msg tea.Msg) (Subnav, tea.Cmd) {
	return s, nil
}

// SetNodes 设置节点列表（Tab 模式，无激活项）
func (s *Subnav) SetNodes(nodes []model.Node) {
	s.nodes = nodes
	s.cursor = 0  // 光标默认第一个
	s.active = -1 // Tab 模式下无激活项
	s.offset = 0
}

// SetActiveNode 设置激活的节点（节点模式）
func (s *Subnav) SetActiveNode(nodeCode string) {
	for i, node := range s.nodes {
		if node.Code == nodeCode {
			s.cursor = i
			s.active = i
			return
		}
	}
}

// Activate 激活当前光标位置的项
func (s *Subnav) Activate() {
	if s.cursor >= 0 && s.cursor < len(s.nodes) {
		s.active = s.cursor
	}
}

// SetFocused 设置焦点状态
func (s *Subnav) SetFocused(focused bool) {
	s.focused = focused
	if focused && s.active >= 0 {
		// 获得焦点时，光标恢复到激活项位置
		s.cursor = s.active
	}
}

// SetWidth 设置宽度
func (s *Subnav) SetWidth(width int) {
	s.width = width
}

// Selected 获取当前光标位置
func (s Subnav) Selected() int {
	return s.cursor
}

// SelectedNode 获取当前光标位置的节点
func (s Subnav) SelectedNode() model.Node {
	if s.cursor >= 0 && s.cursor < len(s.nodes) {
		return s.nodes[s.cursor]
	}
	return model.Node{}
}

// HasNodes 是否有节点
func (s Subnav) HasNodes() bool {
	return len(s.nodes) > 0
}

// MoveLeft 向左移动光标
func (s *Subnav) MoveLeft() {
	if s.cursor > 0 {
		s.cursor--
	}
}

// MoveRight 向右移动光标
func (s *Subnav) MoveRight() {
	if s.cursor < len(s.nodes)-1 {
		s.cursor++
	}
}

// View 渲染二级导航
func (s Subnav) View() string {
	// 没有节点时不显示
	if len(s.nodes) == 0 {
		return ""
	}

	var items []string
	for i, node := range s.nodes {
		var style lipgloss.Style
		isCursor := (i == s.cursor)
		isActive := (i == s.active)

		if s.focused && isCursor {
			// 有焦点时，光标位置高亮背景
			style = lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.PrimaryFg).
				Background(ui.CurrentTheme.Secondary).
				Bold(true).
				Padding(0, 1)
		} else if !s.focused && isActive && s.active >= 0 {
			// 无焦点时，激活项高亮（绿色文字）
			style = lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.Secondary).
				Bold(true).
				Padding(0, 1)
		} else {
			// 其他项 - 普通文字
			style = lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.Foreground).
				Padding(0, 1)
		}
		items = append(items, style.Render(node.Name))
	}

	content := strings.Join(items, "")

	// 如果内容太长，添加滚动指示
	contentWidth := lipgloss.Width(content)
	if contentWidth > s.width-4 {
		content = content + " ▸"
	}

	// 添加分隔线
	var separator string
	if s.focused {
		separator = lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Secondary).
			Render(strings.Repeat("─", s.width))
	} else {
		separator = lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Border).
			Render(strings.Repeat("─", s.width))
	}

	return content + "\n" + separator
}
