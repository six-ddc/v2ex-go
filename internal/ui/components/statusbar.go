package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-tui/internal/model"
	"github.com/six-ddc/v2ex-tui/internal/ui"
)

// StatusBar 顶部状态栏组件
type StatusBar struct {
	user  *model.User
	width int
}

// NewStatusBar 创建状态栏
func NewStatusBar() StatusBar {
	return StatusBar{}
}

// SetUser 设置用户信息
func (s *StatusBar) SetUser(user *model.User) {
	s.user = user
}

// SetWidth 设置宽度
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// View 渲染状态栏
func (s StatusBar) View() string {
	title := lipgloss.NewStyle().
		Foreground(ui.CurrentTheme.PrimaryFg).
		Background(ui.CurrentTheme.Primary).
		Bold(true).
		Padding(0, 1).
		Render("V2EX Terminal")

	var userInfo string
	if s.user != nil && s.user.LoggedIn {
		// 已登录用户信息
		notify := ""
		if s.user.Notify > 0 {
			notify = fmt.Sprintf(" 🔔%d", s.user.Notify)
		}
		balance := fmt.Sprintf("💰%d/%d", s.user.Silver, s.user.Bronze)
		userInfo = lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.PrimaryFg).
			Background(ui.CurrentTheme.Primary).
			Padding(0, 1).
			Render(fmt.Sprintf("[%s]%s %s", s.user.Name, notify, balance))
	} else {
		userInfo = lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.PrimaryFg).
			Background(ui.CurrentTheme.Primary).
			Padding(0, 1).
			Render("[未登录]")
	}

	// 计算中间填充
	titleWidth := lipgloss.Width(title)
	userWidth := lipgloss.Width(userInfo)
	padding := s.width - titleWidth - userWidth
	if padding < 0 {
		padding = 0
	}

	spacer := lipgloss.NewStyle().
		Background(ui.CurrentTheme.Primary).
		Width(padding).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, title, spacer, userInfo)
}
