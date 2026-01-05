package app

import (
	"github.com/six-ddc/v2ex-go/internal/api"
	"github.com/six-ddc/v2ex-go/internal/model"
	"github.com/six-ddc/v2ex-go/internal/ui"
	"github.com/six-ddc/v2ex-go/internal/ui/components"
)

// View 当前视图类型
type View int

const (
	ViewMain   View = iota // 主视图 (导航 + 主题列表)
	ViewDetail             // 帖子详情视图
	ViewSearch             // 搜索模式
	ViewHelp               // 帮助面板
)

// Pane 焦点区域
type Pane int

const (
	PaneNavbar    Pane = iota // 一级导航栏
	PaneSubnav                // 二级节点导航
	PaneTopicList             // 主题列表
)

// Model 应用主 Model
type Model struct {
	// 状态
	currentView View
	focusedPane Pane
	loading     bool
	err         error
	ready       bool
	initialized bool // 是否已经初始化加载过数据

	// 数据
	client      *api.Client
	user        *model.User
	tabs        []model.Tab
	subNodes    []model.Node
	topics      []model.Topic
	topicDetail *model.Topic
	replies     []model.Reply

	// 节点模式
	nodeMode    bool   // 是否处于节点浏览模式
	currentNode string // 当前节点代码
	nodePage    int    // 节点页码

	// 详情页回复分页
	replyPage int // 当前回复页码

	// UI 组件
	statusBar components.StatusBar
	navbar    components.Navbar
	subnav    components.Subnav
	topicList components.TopicList
	detail    components.Detail
	helpBar   components.HelpBar

	// 配置
	keys   ui.KeyMap
	styles ui.Styles
	width  int
	height int
}

// NewModel 创建新的 Model
func NewModel() Model {
	m := Model{
		currentView: ViewMain,
		focusedPane: PaneTopicList, // 默认焦点在主题列表
		loading:     true,          // 初始加载状态

		client:   api.NewClient(),
		tabs:     model.DefaultTabs,
		nodePage: 1,

		statusBar: components.NewStatusBar(),
		navbar:    components.NewNavbar(),
		subnav:    components.NewSubnav(),
		topicList: components.NewTopicList(),
		detail:    components.NewDetail(),
		helpBar:   components.NewHelpBar(),

		keys:   ui.DefaultKeyMap,
		styles: ui.DefaultStyles,
	}

	// 设置初始焦点
	m.topicList.SetFocused(true)

	return m
}
