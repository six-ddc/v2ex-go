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
	nodes    []model.Node
	selected int
	focused  bool
	width    int
	offset   int // 用于横向滚动
}

// NewSubnav 创建二级导航
func NewSubnav() Subnav {
	return Subnav{}
}

// Init 初始化
func (s Subnav) Init() tea.Cmd {
	return nil
}

// Update 更新
func (s Subnav) Update(msg tea.Msg) (Subnav, tea.Cmd) {
	return s, nil
}

// SetNodes 设置节点列表
func (s *Subnav) SetNodes(nodes []model.Node) {
	s.nodes = nodes
	s.selected = 0
	s.offset = 0
}

// SetFocused 设置焦点状态
func (s *Subnav) SetFocused(focused bool) {
	s.focused = focused
}

// SetWidth 设置宽度
func (s *Subnav) SetWidth(width int) {
	s.width = width
}

// Selected 获取当前选中的索引
func (s Subnav) Selected() int {
	return s.selected
}

// SelectedNode 获取当前选中的节点
func (s Subnav) SelectedNode() model.Node {
	if s.selected >= 0 && s.selected < len(s.nodes) {
		return s.nodes[s.selected]
	}
	return model.Node{}
}

// HasNodes 是否有节点
func (s Subnav) HasNodes() bool {
	return len(s.nodes) > 0
}

// MoveLeft 向左移动
func (s *Subnav) MoveLeft() {
	if s.selected > 0 {
		s.selected--
	}
}

// MoveRight 向右移动
func (s *Subnav) MoveRight() {
	if s.selected < len(s.nodes)-1 {
		s.selected++
	}
}

// View 渲染二级导航
func (s Subnav) View() string {
	var content string

	if len(s.nodes) == 0 {
		// 没有子节点时显示提示
		emptyStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Muted)
		content = emptyStyle.Render("  (无子节点)")
	} else {
		var items []string

		for i, node := range s.nodes {
			var style lipgloss.Style
			if i == s.selected {
				if s.focused {
					// 选中且有焦点 - 高亮背景（语义上的高亮）
					style = lipgloss.NewStyle().
						Foreground(ui.CurrentTheme.PrimaryFg).
						Background(ui.CurrentTheme.Secondary).
						Bold(true).
						Padding(0, 1)
				} else {
					// 选中但无焦点 - 绿色文字，不设置背景色
					style = lipgloss.NewStyle().
						Foreground(ui.CurrentTheme.Secondary).
						Bold(true).
						Padding(0, 1)
				}
			} else {
				// 未选中 - 前景色文字，不设置背景色
				style = lipgloss.NewStyle().
					Foreground(ui.CurrentTheme.Foreground).
					Padding(0, 1)
			}
			items = append(items, style.Render(node.Name))
		}

		content = strings.Join(items, "")

		// 如果内容太长，添加滚动指示
		contentWidth := lipgloss.Width(content)
		if contentWidth > s.width-4 {
			content = content + " ▸"
		}
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
