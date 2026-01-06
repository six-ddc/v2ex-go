package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-tui/internal/model"
	"github.com/six-ddc/v2ex-tui/internal/ui"
)

// TopicList 主题列表组件
type TopicList struct {
	topics   []model.Topic
	selected int
	offset   int // 滚动偏移
	focused  bool
	width    int
	height   int
}

// NewTopicList 创建主题列表
func NewTopicList() TopicList {
	return TopicList{}
}

// Init 初始化
func (t TopicList) Init() tea.Cmd {
	return nil
}

// Update 更新
func (t TopicList) Update(msg tea.Msg) (TopicList, tea.Cmd) {
	return t, nil
}

// SetTopics 设置主题列表
func (t *TopicList) SetTopics(topics []model.Topic) {
	t.topics = topics
	t.selected = 0
	t.offset = 0
}

// SetFocused 设置焦点状态
func (t *TopicList) SetFocused(focused bool) {
	t.focused = focused
}

// SetSize 设置尺寸
func (t *TopicList) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// Selected 获取当前选中的索引
func (t TopicList) Selected() int {
	return t.selected
}

// SelectedTopic 获取当前选中的主题
func (t TopicList) SelectedTopic() model.Topic {
	if t.selected >= 0 && t.selected < len(t.topics) {
		return t.topics[t.selected]
	}
	return model.Topic{}
}

// Topics 获取主题列表
func (t TopicList) Topics() []model.Topic {
	return t.topics
}

// Len 获取列表长度
func (t TopicList) Len() int {
	return len(t.topics)
}

// MoveUp 向上移动
func (t *TopicList) MoveUp() {
	if t.selected > 0 {
		t.selected--
		t.ensureVisible()
	}
}

// MoveDown 向下移动
func (t *TopicList) MoveDown() {
	if t.selected < len(t.topics)-1 {
		t.selected++
		t.ensureVisible()
	}
}

// GoToTop 跳到顶部
func (t *TopicList) GoToTop() {
	t.selected = 0
	t.offset = 0
}

// GoToBottom 跳到底部
func (t *TopicList) GoToBottom() {
	if len(t.topics) > 0 {
		t.selected = len(t.topics) - 1
		t.ensureVisible()
	}
}

// PageUp 上翻页
func (t *TopicList) PageUp() {
	visibleItems := t.visibleItemCount()
	t.selected -= visibleItems
	if t.selected < 0 {
		t.selected = 0
	}
	t.ensureVisible()
}

// PageDown 下翻页
func (t *TopicList) PageDown() {
	visibleItems := t.visibleItemCount()
	t.selected += visibleItems
	if t.selected >= len(t.topics) {
		t.selected = len(t.topics) - 1
	}
	t.ensureVisible()
}

// HalfPageUp 上翻半页
func (t *TopicList) HalfPageUp() {
	visibleItems := t.visibleItemCount() / 2
	t.selected -= visibleItems
	if t.selected < 0 {
		t.selected = 0
	}
	t.ensureVisible()
}

// HalfPageDown 下翻半页
func (t *TopicList) HalfPageDown() {
	visibleItems := t.visibleItemCount() / 2
	t.selected += visibleItems
	if t.selected >= len(t.topics) {
		t.selected = len(t.topics) - 1
	}
	t.ensureVisible()
}

// visibleItemCount 计算可见项数量
func (t TopicList) visibleItemCount() int {
	itemHeight := 2
	if t.height <= 0 {
		return 10
	}
	return t.height / itemHeight
}

// ensureVisible 确保选中项可见
func (t *TopicList) ensureVisible() {
	visibleItems := t.visibleItemCount()
	if t.selected < t.offset {
		t.offset = t.selected
	} else if t.selected >= t.offset+visibleItems {
		t.offset = t.selected - visibleItems + 1
	}
}

// View 渲染主题列表
func (t TopicList) View() string {
	if len(t.topics) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Muted)
		return emptyStyle.Render("  暂无主题，请稍候...")
	}

	var items []string
	visibleItems := t.visibleItemCount()
	end := t.offset + visibleItems
	if end > len(t.topics) {
		end = len(t.topics)
	}

	contentWidth := t.width - 2

	for i := t.offset; i < end; i++ {
		topic := t.topics[i]
		isSelected := i == t.selected
		item := t.renderTopicItem(topic, i+1, isSelected, contentWidth)
		items = append(items, item)
	}

	return strings.Join(items, "\n")
}


// renderTopicItem 渲染单个主题项
func (t TopicList) renderTopicItem(topic model.Topic, index int, selected bool, width int) string {
	// 只有在有焦点时才显示选中效果
	showSelected := selected && t.focused

	indicator := "  "
	if showSelected {
		indicator = "> "
	}

	indexStr := fmt.Sprintf("%d. ", index)
	titleMaxWidth := width - len(indicator) - len(indexStr) - 4
	if titleMaxWidth < 10 {
		titleMaxWidth = 10
	}
	title := truncateString(topic.Title, titleMaxWidth)

	// 构建标题行
	var line1 string
	if showSelected {
		// 选中项使用反色背景（语义上的高亮）
		titleStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.SelectedFg).
			Background(ui.CurrentTheme.SelectedBg).
			Bold(true)
		prefixStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.SelectedFg).
			Background(ui.CurrentTheme.SelectedBg)
		line1 = prefixStyle.Render(indicator+indexStr) + titleStyle.Render(title)
		// 填充行尾保持选中背景色
		line1Width := lipgloss.Width(line1)
		if line1Width < width {
			line1 += prefixStyle.Render(strings.Repeat(" ", width-line1Width))
		}
	} else {
		// 非选中项不设置背景色
		titleStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Foreground)
		line1 = indicator + indexStr + titleStyle.Render(title)
	}

	// 构建元信息行
	var line2 string
	if showSelected {
		// 选中项使用反色背景
		baseStyle := lipgloss.NewStyle().Background(ui.CurrentTheme.SelectedBg)
		nodeStyle := baseStyle.Foreground(ui.CurrentTheme.NodeColor)
		authorStyle := baseStyle.Foreground(ui.CurrentTheme.AuthorColor)
		timeStyle := baseStyle.Foreground(ui.CurrentTheme.TimeColor)
		normalStyle := baseStyle.Foreground(ui.CurrentTheme.SelectedFg)

		line2 = normalStyle.Render("   [") +
			nodeStyle.Render(topic.Node.Name) +
			normalStyle.Render("] ") +
			authorStyle.Render(topic.Author) +
			normalStyle.Render(" · ") +
			timeStyle.Render(topic.RelativeTime)

		if topic.LastReplyBy != "" {
			line2 += normalStyle.Render(" · ") +
				authorStyle.Render(topic.LastReplyBy)
		}

		// 回复数
		if topic.ReplyCount > 0 {
			countStyle := lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.SelectedFg).
				Background(ui.CurrentTheme.ReplyCountBgSel).
				Padding(0, 1)
			line2 += normalStyle.Render(" ") + countStyle.Render(fmt.Sprintf("%d", topic.ReplyCount))
		}

		// 填充行尾保持选中背景色
		line2Width := lipgloss.Width(line2)
		if line2Width < width {
			line2 += normalStyle.Render(strings.Repeat(" ", width-line2Width))
		}
	} else {
		// 非选中项不设置背景色
		nodeStyle := lipgloss.NewStyle().Foreground(ui.CurrentTheme.NodeColor)
		authorStyle := lipgloss.NewStyle().Foreground(ui.CurrentTheme.AuthorColor)
		timeStyle := lipgloss.NewStyle().Foreground(ui.CurrentTheme.TimeColor)
		normalStyle := lipgloss.NewStyle().Foreground(ui.CurrentTheme.Muted)

		line2 = normalStyle.Render("   [") +
			nodeStyle.Render(topic.Node.Name) +
			normalStyle.Render("] ") +
			authorStyle.Render(topic.Author) +
			normalStyle.Render(" · ") +
			timeStyle.Render(topic.RelativeTime)

		if topic.LastReplyBy != "" {
			line2 += normalStyle.Render(" · ") +
				authorStyle.Render(topic.LastReplyBy)
		}

		// 回复数
		if topic.ReplyCount > 0 {
			countStyle := lipgloss.NewStyle().
				Foreground(ui.CurrentTheme.CountColor).
				Background(ui.CurrentTheme.ReplyCountBg).
				Padding(0, 1)
			line2 += normalStyle.Render(" ") + countStyle.Render(fmt.Sprintf("%d", topic.ReplyCount))
		}
	}

	return line1 + "\n" + line2
}

// sanitizeTitle 清理标题，移除换行符和特殊字符，确保单行显示
func sanitizeTitle(s string) string {
	// 替换各种换行符为空格
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	// 替换零宽字符和其他不可见字符
	s = strings.Map(func(r rune) rune {
		// 保留常规可打印字符
		if r < 32 && r != ' ' {
			return ' '
		}
		// 移除零宽字符
		if r == '\u200B' || r == '\u200C' || r == '\u200D' || r == '\uFEFF' {
			return -1
		}
		return r
	}, s)
	// 合并多个连续空格为单个空格
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	// 先清理标题
	s = sanitizeTitle(s)
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}
