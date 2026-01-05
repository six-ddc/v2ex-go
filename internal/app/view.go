package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/six-ddc/v2ex-go/internal/ui"
	"github.com/six-ddc/v2ex-go/internal/ui/components"
)

// 布局常量
const (
	statusBarHeight = 1
	navbarHeight    = 2 // 包含分隔线
	subnavHeight    = 2 // 包含分隔线
	helpBarHeight   = 1
)

// View 渲染视图
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	switch m.currentView {
	case ViewMain:
		return m.mainView()
	case ViewDetail:
		return m.detailView()
	default:
		return m.mainView()
	}
}

// mainView 渲染主视图
func (m Model) mainView() string {
	// 计算主题列表高度
	listHeight := m.height - statusBarHeight - navbarHeight - subnavHeight - helpBarHeight

	// 状态栏 - 固定 1 行
	statusBar := lipgloss.Place(
		m.width, statusBarHeight,
		lipgloss.Left, lipgloss.Top,
		m.statusBar.View(),
	)

	// 一级导航 - 固定 2 行
	navbar := lipgloss.Place(
		m.width, navbarHeight,
		lipgloss.Left, lipgloss.Top,
		m.navbar.View(),
	)

	// 二级导航 - 固定 2 行
	subnav := lipgloss.Place(
		m.width, subnavHeight,
		lipgloss.Left, lipgloss.Top,
		m.subnav.View(),
	)

	// 主题列表 - 自适应剩余空间
	var topicListContent string
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Info)
		topicListContent = loadingStyle.Render("  加载中...")
	} else if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Error)
		topicListContent = errorStyle.Render(fmt.Sprintf("  加载失败: %v", m.err))
	} else {
		topicListContent = m.topicList.View()
	}
	topicList := lipgloss.Place(
		m.width, listHeight,
		lipgloss.Left, lipgloss.Top,
		topicListContent,
	)

	// 帮助栏 - 固定 1 行
	m.helpBar.SetItems(components.MainViewHelp)
	helpBar := lipgloss.Place(
		m.width, helpBarHeight,
		lipgloss.Left, lipgloss.Top,
		m.helpBar.View(),
	)

	// 使用 JoinVertical 组合所有区域
	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		navbar,
		subnav,
		topicList,
		helpBar,
	)
}

// detailView 渲染详情视图
func (m Model) detailView() string {
	// 计算详情区域高度
	detailHeight := m.height - helpBarHeight

	// 详情内容
	var detailContent string
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(ui.CurrentTheme.Info)
		detailContent = loadingStyle.Render("加载中...")
	} else {
		detailContent = m.detail.View()
	}
	detail := lipgloss.Place(
		m.width, detailHeight,
		lipgloss.Left, lipgloss.Top,
		detailContent,
	)

	// 帮助栏
	m.helpBar.SetItems(components.DetailViewHelp)
	helpBar := lipgloss.Place(
		m.width, helpBarHeight,
		lipgloss.Left, lipgloss.Top,
		m.helpBar.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		detail,
		helpBar,
	)
}
