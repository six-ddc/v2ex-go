package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-go/internal/model"
	"github.com/six-ddc/v2ex-go/internal/ui"
)

// Detail 帖子详情视图组件
type Detail struct {
	topic    *model.Topic
	replies  []model.Reply
	viewport viewport.Model
	width    int
	height   int

	// 导航状态
	currentIndex int // 当前帖子在列表中的索引
	totalCount   int // 列表总数
}

// NewDetail 创建详情视图
func NewDetail() Detail {
	vp := viewport.New(80, 20)
	return Detail{
		viewport: vp,
	}
}

// Init 初始化
func (d Detail) Init() tea.Cmd {
	return nil
}

// Update 更新
func (d Detail) Update(msg tea.Msg) (Detail, tea.Cmd) {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return d, cmd
}

// SetTopic 设置帖子详情
func (d *Detail) SetTopic(topic *model.Topic, replies []model.Reply) {
	d.topic = topic
	d.replies = replies
	d.updateContent()
}

// SetNavInfo 设置导航信息
func (d *Detail) SetNavInfo(currentIndex, totalCount int) {
	d.currentIndex = currentIndex
	d.totalCount = totalCount
}

// SetSize 设置尺寸
func (d *Detail) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.viewport.Width = width - 2
	d.viewport.Height = height - detailHeaderHeight // 减去顶部标题栏
	d.updateContent()
}

// ScrollUp 向上滚动
func (d *Detail) ScrollUp() {
	d.viewport.LineUp(1)
}

// ScrollDown 向下滚动
func (d *Detail) ScrollDown() {
	d.viewport.LineDown(1)
}

// GoToTop 跳到顶部
func (d *Detail) GoToTop() {
	d.viewport.GotoTop()
}

// GoToBottom 跳到底部
func (d *Detail) GoToBottom() {
	d.viewport.GotoBottom()
}

// PageUp 上翻页
func (d *Detail) PageUp() {
	d.viewport.ViewUp()
}

// PageDown 下翻页
func (d *Detail) PageDown() {
	d.viewport.ViewDown()
}

// HalfPageUp 上翻半页
func (d *Detail) HalfPageUp() {
	d.viewport.HalfViewUp()
}

// HalfPageDown 下翻半页
func (d *Detail) HalfPageDown() {
	d.viewport.HalfViewDown()
}

// AtBottom 检测是否滚动到底部
func (d *Detail) AtBottom() bool {
	return d.viewport.AtBottom()
}

// HasNextPage 检测是否还有下一页回复
func (d *Detail) HasNextPage() bool {
	if d.topic == nil {
		return false
	}
	return d.topic.CurrentPage < d.topic.TotalPages
}

// CurrentPage 获取当前回复页码
func (d *Detail) CurrentPage() int {
	if d.topic == nil {
		return 1
	}
	return d.topic.CurrentPage
}

// AppendReplies 追加回复到列表（用于无限滚动加载）
func (d *Detail) AppendReplies(replies []model.Reply, newPage int) {
	// 保存当前滚动位置
	currentYOffset := d.viewport.YOffset

	d.replies = append(d.replies, replies...)
	if d.topic != nil {
		d.topic.CurrentPage = newPage
	}
	d.updateContent()

	// 恢复滚动位置
	d.viewport.SetYOffset(currentYOffset)
}

// updateContent 更新视口内容
func (d *Detail) updateContent() {
	if d.topic == nil {
		d.viewport.SetContent("")
		return
	}

	contentWidth := d.viewport.Width - 2

	var content strings.Builder

	// 标题（可点击跳转到原帖）
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.TitleColor).
		Bold(true)
	titleText := titleStyle.Render(d.topic.Title)
	if d.topic.URL != "" {
		titleText = ui.Hyperlink(d.topic.URL, titleText)
	}
	content.WriteString(titleText)
	content.WriteString("\n")

	// 分隔线
	sepStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Border)
	content.WriteString(sepStyle.Render(strings.Repeat("─", contentWidth)))
	content.WriteString("\n")

	// 作者信息
	authorStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.AuthorColor)
	timeStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.TimeColor)
	normalStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Foreground)
	// 作者名可点击
	authorText := normalStyle.Render("@") + authorStyle.Render(d.topic.Author)
	authorText = ui.Hyperlink("/member/"+d.topic.Author, authorText)
	metaLine := authorText +
		normalStyle.Render(" · ") + timeStyle.Render(d.topic.RelativeTime) +
		normalStyle.Render(fmt.Sprintf(" · 点击 %d · 回复 %d", d.topic.Clicks, d.topic.ReplyCount))
	content.WriteString(metaLine)
	content.WriteString("\n\n")

	// 正文内容
	if d.topic.ContentHTML != "" {
		rendered := ui.RenderHTML(d.topic.ContentHTML)
		content.WriteString(rendered)
		content.WriteString("\n")
	} else if d.topic.Content != "" {
		content.WriteString(d.topic.Content)
		content.WriteString("\n")
	}

	// 渲染附言
	if len(d.topic.Supplements) > 0 {
		content.WriteString("\n")
		for _, supplement := range d.topic.Supplements {
			content.WriteString(d.renderSupplement(supplement, contentWidth))
			content.WriteString("\n")
		}
	}

	// 回复区分隔线和标题
	content.WriteString("\n")
	content.WriteString(sepStyle.Render(strings.Repeat("─", contentWidth)))
	content.WriteString("\n")

	// 回复统计和分页信息
	replyHeaderStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Muted)
	if d.topic.TotalPages > 1 {
		if d.topic.CurrentPage < d.topic.TotalPages {
			content.WriteString(replyHeaderStyle.Render(
				fmt.Sprintf("共 %d 条回复 · 已加载 %d/%d 页 · 滚动到底部自动加载更多",
					d.topic.ReplyCount, d.topic.CurrentPage, d.topic.TotalPages)))
		} else {
			content.WriteString(replyHeaderStyle.Render(
				fmt.Sprintf("共 %d 条回复 · 已全部加载",
					d.topic.ReplyCount)))
		}
	} else if d.topic.ReplyCount > 0 {
		content.WriteString(replyHeaderStyle.Render(
			fmt.Sprintf("共 %d 条回复", d.topic.ReplyCount)))
	}
	content.WriteString("\n\n")

	// 回复列表
	for _, reply := range d.replies {
		content.WriteString(d.renderReply(reply, contentWidth))
		content.WriteString("\n")
	}

	d.viewport.SetContent(content.String())
}

// renderSupplement 渲染附言
func (d Detail) renderSupplement(supplement model.Supplement, width int) string {
	var content strings.Builder

	// 附言框样式 - 使用虚线边框区分正文
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.CurrentTheme.Warning).
		Padding(0, 1).
		Width(width)

	// 附言头部
	headerStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Warning).
		Bold(true)
	timeStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.TimeColor)
	normalStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Foreground)

	header := headerStyle.Render(fmt.Sprintf("第 %d 条附言", supplement.Index))
	if supplement.RelativeTime != "" {
		header += normalStyle.Render(" · ") + timeStyle.Render(supplement.RelativeTime)
	}
	content.WriteString(header)
	content.WriteString("\n")

	// 附言内容
	if supplement.ContentHTML != "" {
		rendered := ui.RenderHTML(supplement.ContentHTML)
		content.WriteString(rendered)
	} else if supplement.Content != "" {
		content.WriteString(supplement.Content)
	}

	return boxStyle.Render(content.String())
}

// renderReply 渲染单条回复
func (d Detail) renderReply(reply model.Reply, width int) string {
	var content strings.Builder

	// 回复框样式 - 使用完整宽度
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.CurrentTheme.Border).
		Padding(0, 1).
		Width(width)

	// 头部: #楼层 作者 · 时间
	floorStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Muted)
	authorStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.AuthorColor)
	timeStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.TimeColor)
	normalStyle := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.Foreground)

	// 作者名可点击
	authorText := normalStyle.Render(" @") + authorStyle.Render(reply.Author)
	authorText = ui.Hyperlink("/member/"+reply.Author, authorText)

	header := floorStyle.Render(fmt.Sprintf("#%d", reply.Floor)) + authorText

	// 如果是楼主
	if reply.IsOP {
		opBadge := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.OPColor).
			Bold(true).
			Render(" [OP]")
		header += opBadge
	}

	if reply.RelativeTime != "" {
		header += normalStyle.Render(" · ") + timeStyle.Render(reply.RelativeTime)
	}

	// 感谢数
	if reply.Likes > 0 {
		likesStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.OPColor)
		header += normalStyle.Render(" · ") + likesStyle.Render(fmt.Sprintf("❤️ %d", reply.Likes))
	}

	content.WriteString(header)
	content.WriteString("\n")

	// 回复内容
	if reply.ContentHTML != "" {
		rendered := ui.RenderHTML(reply.ContentHTML)
		content.WriteString(rendered)
	} else if reply.Content != "" {
		content.WriteString(reply.Content)
	}

	return boxStyle.Render(content.String())
}

// 布局常量
const (
	detailHeaderHeight = 1
)

// renderHeader 渲染顶部导航栏
func (d Detail) renderHeader() string {
	// 使用反转颜色作为header背景（这是语义上需要的高亮）
	style := lipgloss.NewStyle().
		Background(ui.CurrentTheme.Primary).
		Foreground(ui.CurrentTheme.PrimaryFg).
		Width(d.width)

	// 简单拼接：返回 | 节点 | 页码
	content := fmt.Sprintf("  ← [q] 返回    [ %s ]    %d / %d",
		d.topic.Node.Name,
		d.currentIndex+1,
		d.totalCount,
	)

	return style.Render(content)
}

// View 渲染详情视图
func (d Detail) View() string {
	if d.topic == nil {
		return ""
	}

	// 顶部栏 - 强制固定 1 行高度
	header := lipgloss.Place(
		d.width, detailHeaderHeight,
		lipgloss.Left, lipgloss.Top,
		d.renderHeader(),
	)

	// 内容区域 - 强制固定高度，确保不会溢出
	viewportHeight := d.height - detailHeaderHeight
	content := lipgloss.Place(
		d.width, viewportHeight,
		lipgloss.Left, lipgloss.Top,
		d.viewport.View(),
	)

	// 组合
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
	)
}
