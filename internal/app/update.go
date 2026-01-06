package app

import (
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/six-ddc/v2ex-tui/internal/model"
	"github.com/six-ddc/v2ex-tui/internal/ui"
	"github.com/six-ddc/v2ex-tui/internal/ui/components"
)

// Init 初始化
func (m Model) Init() tea.Cmd {
	// 启动时加载默认 Tab 的数据
	return m.loadTopicsByTab("all")
}

// Update 更新逻辑
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateComponentSizes()
		m.ready = true
		return m, nil

	case topicsLoadedMsg:
		m.loading = false
		m.initialized = true
		m.err = nil
		m.topics = msg.topics
		// Tab 模式下更新二级节点列表（可能为空）
		// 节点模式下保持当前节点列表不变
		if !m.nodeMode {
			m.subNodes = msg.nodes
			m.subnav.SetNodes(m.subNodes)
			// 节点列表变化后需要重新计算布局高度
			m.updateComponentSizes()
		} else {
			// 节点模式下更新总页数
			m.nodeTotalPages = msg.totalPages
		}
		if msg.user != nil {
			m.user = msg.user
			m.statusBar.SetUser(m.user)
		}
		m.topicList.SetTopics(m.topics)
		return m, nil

	case topicDetailLoadedMsg:
		m.loading = false
		m.err = nil
		m.topicDetail = msg.topic
		m.replies = msg.replies
		m.replyPage = msg.topic.CurrentPage
		m.detail.SetTopic(m.topicDetail, m.replies)
		m.detail.SetNavInfo(m.topicList.Selected(), m.topicList.Len())
		m.currentView = ViewDetail
		return m, nil

	case replyPageLoadedMsg:
		m.err = nil
		// 追加新回复到现有列表（无限滚动模式）
		m.topicDetail.TotalPages = msg.totalPages
		m.replies = append(m.replies, msg.replies...)
		m.replyPage = msg.page
		m.detail.SetLoading(false) // 先清除加载状态
		m.detail.AppendReplies(msg.replies, msg.page)
		return m, nil

	case nodeTopicsAppendedMsg:
		m.err = nil
		// 追加新主题到列表（无限滚动模式）
		m.nodePage = msg.page
		m.nodeTotalPages = msg.totalPages
		// 只更新 topicList，m.topics 不再单独维护
		m.topicList.AppendTopics(msg.topics)
		m.topicList.SetLoading(false) // 清除加载状态
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg.err
		// 清除组件级加载状态
		m.topicList.SetLoading(false)
		m.detail.SetLoading(false)
		return m, nil
	}

	return m, nil
}

// handleKeyMsg 处理键盘输入
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 全局快捷键
	switch {
	case key.Matches(msg, m.keys.Quit):
		if m.currentView == ViewDetail {
			// 从详情返回列表
			m.currentView = ViewMain
			return m, nil
		}
		return m, tea.Quit

	case key.Matches(msg, m.keys.Escape):
		if m.currentView == ViewDetail {
			m.currentView = ViewMain
			return m, nil
		}
		return m, nil

	case key.Matches(msg, m.keys.Help):
		return m, nil

	case key.Matches(msg, m.keys.Refresh):
		return m, m.refresh()

	case key.Matches(msg, m.keys.ToggleTheme):
		ui.ToggleTheme()
		return m, nil
	}

	// 根据当前视图处理
	switch m.currentView {
	case ViewMain:
		return m.handleMainViewKey(msg)
	case ViewDetail:
		return m.handleDetailViewKey(msg)
	}

	return m, nil
}

// handleMainViewKey 主视图键盘处理
func (m Model) handleMainViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Tab 切换焦点
	switch {
	case key.Matches(msg, m.keys.Tab):
		m.cycleFocus(true)
		return m, nil
	case key.Matches(msg, m.keys.ShiftTab):
		m.cycleFocus(false)
		return m, nil
	}

	// 数字快捷键 - 切换 Tab
	switch {
	case key.Matches(msg, m.keys.Num1):
		return m.jumpToTab(0)
	case key.Matches(msg, m.keys.Num2):
		return m.jumpToTab(1)
	case key.Matches(msg, m.keys.Num3):
		return m.jumpToTab(2)
	case key.Matches(msg, m.keys.Num4):
		return m.jumpToTab(3)
	case key.Matches(msg, m.keys.Num5):
		return m.jumpToTab(4)
	case key.Matches(msg, m.keys.Num6):
		return m.jumpToTab(5)
	case key.Matches(msg, m.keys.Num7):
		return m.jumpToTab(6)
	case key.Matches(msg, m.keys.Num8):
		return m.jumpToTab(7)
	case key.Matches(msg, m.keys.Num9):
		return m.jumpToTab(8)
	case key.Matches(msg, m.keys.Num0):
		return m.jumpToTab(9)
	}

	// 根据焦点区域处理
	switch m.focusedPane {
	case PaneNavbar:
		return m.handleNavbarKey(msg)
	case PaneSubnav:
		return m.handleSubnavKey(msg)
	case PaneTopicList:
		return m.handleTopicListKey(msg)
	}

	return m, nil
}

// handleNavbarKey 处理导航栏按键
func (m Model) handleNavbarKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left):
		m.navbar.MoveLeft()
		return m, nil
	case key.Matches(msg, m.keys.Right):
		m.navbar.MoveRight()
		return m, nil
	case key.Matches(msg, m.keys.Down):
		// 下键切换到下一组（subnav 或 topiclist）
		m.updateFocusState(false)
		if m.subnav.HasNodes() {
			m.focusedPane = PaneSubnav
		} else {
			m.focusedPane = PaneTopicList
		}
		m.updateFocusState(true)
		return m, nil
	case key.Matches(msg, m.keys.Enter):
		// 激活当前光标位置的 Tab
		m.navbar.Activate()
		tab := m.navbar.SelectedTab()
		m.loading = true
		m.nodeMode = false
		// 切换后焦点移到帖子列表
		m.switchFocusToTopicList()
		return m, m.loadTopicsByTab(tab.Code)
	}
	return m, nil
}

// handleSubnavKey 处理二级导航按键
func (m Model) handleSubnavKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left):
		m.subnav.MoveLeft()
		return m, nil
	case key.Matches(msg, m.keys.Right):
		m.subnav.MoveRight()
		return m, nil
	case key.Matches(msg, m.keys.Up):
		// 上键切换到 navbar
		m.updateFocusState(false)
		m.focusedPane = PaneNavbar
		m.updateFocusState(true)
		return m, nil
	case key.Matches(msg, m.keys.Down):
		// 下键切换到 topiclist
		m.updateFocusState(false)
		m.focusedPane = PaneTopicList
		m.updateFocusState(true)
		return m, nil
	case key.Matches(msg, m.keys.Enter):
		// 激活当前光标位置的节点
		m.subnav.Activate()
		node := m.subnav.SelectedNode()
		if node.Code != "" {
			m.loading = true
			m.nodeMode = true
			m.currentNode = node.Code
			m.nodePage = 1
			// 切换后焦点移到帖子列表
			m.switchFocusToTopicList()
			return m, m.loadTopicsByNode(node.Code, 1)
		}
		return m, nil
	}
	return m, nil
}

// handleTopicListKey 处理主题列表按键
func (m Model) handleTopicListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		// 在第一个帖子时，上键切换到上一组
		if m.topicList.Selected() == 0 {
			m.updateFocusState(false)
			if m.subnav.HasNodes() {
				m.focusedPane = PaneSubnav
			} else {
				m.focusedPane = PaneNavbar
			}
			m.updateFocusState(true)
			return m, nil
		}
		m.topicList.MoveUp()
		return m, nil
	case key.Matches(msg, m.keys.Down):
		m.topicList.MoveDown()
		return m.checkAndLoadMoreTopics()
	case key.Matches(msg, m.keys.Enter), key.Matches(msg, m.keys.Right):
		// 打开帖子详情
		topic := m.topicList.SelectedTopic()
		if topic.URL != "" {
			m.loading = true
			return m, m.loadTopicDetail(topic.URL)
		}
		return m, nil
	case key.Matches(msg, m.keys.Top):
		m.topicList.GoToTop()
		return m, nil
	case key.Matches(msg, m.keys.Bottom):
		m.topicList.GoToBottom()
		return m.checkAndLoadMoreTopics()
	case key.Matches(msg, m.keys.PageUp):
		m.topicList.PageUp()
		return m, nil
	case key.Matches(msg, m.keys.PageDown):
		m.topicList.PageDown()
		return m.checkAndLoadMoreTopics()
	case key.Matches(msg, m.keys.HalfPageUp):
		m.topicList.HalfPageUp()
		return m, nil
	case key.Matches(msg, m.keys.HalfPageDown):
		m.topicList.HalfPageDown()
		return m.checkAndLoadMoreTopics()
	}
	return m, nil
}

// checkAndLoadMoreTopics 检查是否到底部并加载更多主题（节点模式）
func (m Model) checkAndLoadMoreTopics() (tea.Model, tea.Cmd) {
	// 仅在节点模式下生效
	if !m.nodeMode {
		return m, nil
	}

	// 如果已经在加载中，不重复触发
	if m.topicList.IsLoading() {
		return m, nil
	}

	// 检测是否到底部且还有下一页
	if m.topicList.AtBottom() && m.nodePage < m.nodeTotalPages {
		return m.loadMoreTopics()
	}

	return m, nil
}

// loadMoreTopics 加载更多主题（无限滚动）
func (m Model) loadMoreTopics() (tea.Model, tea.Cmd) {
	if !m.nodeMode || m.topicList.IsLoading() {
		return m, nil
	}

	nextPage := m.nodePage + 1
	if nextPage > m.nodeTotalPages {
		return m, nil
	}

	m.topicList.SetLoading(true) // 设置组件级加载状态
	return m, m.loadMoreNodeTopics(m.currentNode, nextPage)
}

// handleDetailViewKey 详情视图键盘处理
func (m Model) handleDetailViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 使用 tea.KeyType 精确匹配按键类型
	switch msg.Type {
	case tea.KeyDown:
		m.detail.ScrollDown()
		return m.checkAndLoadMoreReplies()
	case tea.KeyUp:
		m.detail.ScrollUp()
		return m, nil
	case tea.KeyPgDown:
		m.detail.PageDown()
		return m.checkAndLoadMoreReplies()
	case tea.KeyPgUp:
		m.detail.PageUp()
		return m, nil
	case tea.KeySpace:
		m.detail.PageDown()
		return m.checkAndLoadMoreReplies()
	case tea.KeyRunes:
		// 处理字符按键
		switch string(msg.Runes) {
		case "j":
			m.detail.ScrollDown()
			return m.checkAndLoadMoreReplies()
		case "k":
			m.detail.ScrollUp()
			return m, nil
		case "g":
			m.detail.GoToTop()
			return m, nil
		case "G":
			m.detail.GoToBottom()
			return m.checkAndLoadMoreReplies()
		case "[":
			return m.navigateTopic(-1)
		case "]":
			return m.navigateTopic(1)
		case "o":
			// macOS: 使用 open 命令在浏览器中打开
			if runtime.GOOS == "darwin" && m.topicDetail != nil && m.topicDetail.URL != "" {
				exec.Command("open", m.topicDetail.URL).Start()
			}
			return m, nil
		}
	case tea.KeyCtrlF:
		m.detail.PageDown()
		return m.checkAndLoadMoreReplies()
	case tea.KeyCtrlB:
		m.detail.PageUp()
		return m, nil
	case tea.KeyCtrlD:
		m.detail.HalfPageDown()
		return m.checkAndLoadMoreReplies()
	case tea.KeyCtrlU:
		m.detail.HalfPageUp()
		return m, nil
	}

	// 不处理其他按键，保持当前状态
	return m, nil
}

// checkAndLoadMoreReplies 检查是否到底部并加载更多回复
func (m Model) checkAndLoadMoreReplies() (tea.Model, tea.Cmd) {
	// 如果已经在加载中，不重复触发
	if m.detail.IsLoading() {
		return m, nil
	}

	// 检测是否到底部且还有下一页
	if m.detail.AtBottom() && m.detail.HasNextPage() {
		return m.loadMoreReplies()
	}

	return m, nil
}

// loadMoreReplies 加载更多回复
func (m Model) loadMoreReplies() (tea.Model, tea.Cmd) {
	if m.topicDetail == nil || m.detail.IsLoading() {
		return m, nil
	}

	nextPage := m.replyPage + 1
	if nextPage > m.topicDetail.TotalPages {
		return m, nil
	}

	m.detail.SetLoading(true) // 设置组件级加载状态
	return m, m.loadReplyPage(m.topicDetail.URL, nextPage)
}

// cycleFocus 循环切换焦点
func (m *Model) cycleFocus(forward bool) {
	m.updateFocusState(false)

	if forward {
		switch m.focusedPane {
		case PaneNavbar:
			if m.subnav.HasNodes() {
				m.focusedPane = PaneSubnav
			} else {
				m.focusedPane = PaneTopicList
			}
		case PaneSubnav:
			m.focusedPane = PaneTopicList
		case PaneTopicList:
			m.focusedPane = PaneNavbar
		}
	} else {
		switch m.focusedPane {
		case PaneNavbar:
			m.focusedPane = PaneTopicList
		case PaneSubnav:
			m.focusedPane = PaneNavbar
		case PaneTopicList:
			if m.subnav.HasNodes() {
				m.focusedPane = PaneSubnav
			} else {
				m.focusedPane = PaneNavbar
			}
		}
	}

	m.updateFocusState(true)
}

// updateFocusState 更新焦点状态
func (m *Model) updateFocusState(focused bool) {
	switch m.focusedPane {
	case PaneNavbar:
		m.navbar.SetFocused(focused)
	case PaneSubnav:
		m.subnav.SetFocused(focused)
	case PaneTopicList:
		m.topicList.SetFocused(focused)
	}
}

// jumpToTab 跳转到指定 Tab
func (m Model) jumpToTab(index int) (tea.Model, tea.Cmd) {
	if index < len(m.tabs) {
		m.navbar.SetSelected(index)
		tab := m.navbar.SelectedTab()
		m.loading = true
		m.nodeMode = false
		// 切换后焦点移到帖子列表
		m.switchFocusToTopicList()
		return m, m.loadTopicsByTab(tab.Code)
	}
	return m, nil
}

// switchFocusToTopicList 切换焦点到帖子列表并选中第一条
func (m *Model) switchFocusToTopicList() {
	// 先取消当前焦点
	m.updateFocusState(false)
	// 切换到帖子列表
	m.focusedPane = PaneTopicList
	// 设置新焦点
	m.updateFocusState(true)
}

// navigateTopic 导航到上/下一篇帖子
func (m Model) navigateTopic(direction int) (tea.Model, tea.Cmd) {
	currentIdx := m.topicList.Selected()
	newIdx := currentIdx + direction

	if newIdx >= 0 && newIdx < m.topicList.Len() {
		if direction > 0 {
			m.topicList.MoveDown()
		} else {
			m.topicList.MoveUp()
		}

		topic := m.topicList.SelectedTopic()
		if topic.URL != "" {
			m.loading = true
			return m, m.loadTopicDetail(topic.URL)
		}
	}
	return m, nil
}

// refresh 刷新当前视图
func (m *Model) refresh() tea.Cmd {
	m.loading = true
	if m.nodeMode {
		return m.loadTopicsByNode(m.currentNode, m.nodePage)
	}
	tab := m.navbar.SelectedTab()
	return m.loadTopicsByTab(tab.Code)
}

// updateComponentSizes 更新组件尺寸
func (m *Model) updateComponentSizes() {
	// 使用 view.go 中定义的布局常量
	// statusBarHeight = 1, navbarHeight = 2, subnavHeight = 2, helpBarHeight = 1

	m.statusBar.SetWidth(m.width)
	m.navbar.SetWidth(m.width)
	m.subnav.SetWidth(m.width)
	m.helpBar.SetWidth(m.width)

	// 计算实际的 subnav 高度（有节点时显示，无节点时不显示）
	actualSubnavHeight := 0
	if m.subnav.HasNodes() {
		actualSubnavHeight = subnavHeight
	}

	// 主题列表高度 = 总高度 - 固定区域
	listHeight := m.height - statusBarHeight - navbarHeight - actualSubnavHeight - helpBarHeight
	if listHeight < 5 {
		listHeight = 5
	}
	m.topicList.SetSize(m.width, listHeight)

	// 详情视图高度 = 总高度 - 帮助栏
	detailHeight := m.height - helpBarHeight
	m.detail.SetSize(m.width, detailHeight)

	// 设置帮助栏内容
	m.helpBar.SetItems(components.MainViewHelp)
}

// Messages

type topicsLoadedMsg struct {
	topics     []model.Topic
	nodes      []model.Node
	user       *model.User
	totalPages int // 节点模式下的总页数
}

type topicDetailLoadedMsg struct {
	topic   *model.Topic
	replies []model.Reply
}

type replyPageLoadedMsg struct {
	replies    []model.Reply
	page       int
	totalPages int
}

// nodeTopicsAppendedMsg 节点主题追加消息（用于无限滚动）
type nodeTopicsAppendedMsg struct {
	topics     []model.Topic
	page       int
	totalPages int
}

type errMsg struct {
	err error
}

// Commands

func (m *Model) loadTopicsByTab(tab string) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		topics, nodes, user, err := client.GetTopicsByTab(tab)
		if err != nil {
			return errMsg{err: err}
		}
		return topicsLoadedMsg{
			topics: topics,
			nodes:  nodes,
			user:   user,
		}
	}
}

func (m *Model) loadTopicsByNode(nodeCode string, page int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		result, err := client.GetTopicsByNode(nodeCode, page)
		if err != nil {
			return errMsg{err: err}
		}
		return topicsLoadedMsg{
			topics:     result.Topics,
			nodes:      nil,
			user:       nil,
			totalPages: result.TotalPages,
		}
	}
}

// loadMoreNodeTopics 加载更多节点主题（追加模式）
func (m *Model) loadMoreNodeTopics(nodeCode string, page int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		result, err := client.GetTopicsByNode(nodeCode, page)
		if err != nil {
			return errMsg{err: err}
		}
		return nodeTopicsAppendedMsg{
			topics:     result.Topics,
			page:       page,
			totalPages: result.TotalPages,
		}
	}
}

func (m *Model) loadTopicDetail(url string) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		topic, replies, err := client.GetTopicDetail(url)
		if err != nil {
			return errMsg{err: err}
		}
		return topicDetailLoadedMsg{
			topic:   topic,
			replies: replies,
		}
	}
}

func (m *Model) loadReplyPage(topicURL string, page int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		topic, replies, err := client.GetTopicDetailPage(topicURL, page)
		if err != nil {
			return errMsg{err: err}
		}
		return replyPageLoadedMsg{
			replies:    replies,
			page:       page,
			totalPages: topic.TotalPages,
		}
	}
}
