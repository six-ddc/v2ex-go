package ui

import "github.com/charmbracelet/lipgloss"

// Styles UI 样式集合
type Styles struct {
	// 容器
	App       lipgloss.Style
	StatusBar lipgloss.Style
	HelpBar   lipgloss.Style

	// 导航栏
	Navbar        lipgloss.Style
	NavbarFocused lipgloss.Style
	NavItem       lipgloss.Style
	NavItemActive lipgloss.Style
	NavItemHover  lipgloss.Style

	// 二级导航
	Subnav        lipgloss.Style
	SubnavFocused lipgloss.Style
	SubnavItem    lipgloss.Style
	SubnavActive  lipgloss.Style

	// 列表
	TopicList        lipgloss.Style
	TopicListFocused lipgloss.Style
	ListItem         lipgloss.Style
	ListItemSelected lipgloss.Style
	ListItemTitle    lipgloss.Style
	ListItemMeta     lipgloss.Style

	// 帖子详情
	DetailContainer lipgloss.Style
	DetailHeader    lipgloss.Style
	Title           lipgloss.Style
	Author          lipgloss.Style
	Node            lipgloss.Style
	Time            lipgloss.Style
	Content         lipgloss.Style
	Separator       lipgloss.Style

	// 回复
	ReplyBox     lipgloss.Style
	ReplyHeader  lipgloss.Style
	ReplyContent lipgloss.Style
	Floor        lipgloss.Style
	OPBadge      lipgloss.Style

	// 搜索
	SearchContainer lipgloss.Style
	SearchInput     lipgloss.Style
	SearchPrompt    lipgloss.Style
	MatchHighlight  lipgloss.Style
	MatchCount      lipgloss.Style

	// 通用
	FocusedBorder lipgloss.Style
	NormalBorder  lipgloss.Style
}

// NewStyles 创建样式集合
func NewStyles(theme Theme) Styles {
	return Styles{
		// 容器
		App: lipgloss.NewStyle(),

		StatusBar: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(theme.Primary).
			Padding(0, 1),

		HelpBar: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Background(theme.Background).
			Padding(0, 1),

		// 导航栏
		Navbar: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Border).
			BorderBottom(true).
			Padding(0, 1),

		NavbarFocused: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Primary).
			BorderBottom(true).
			Padding(0, 1),

		NavItem: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(0, 1),

		NavItemActive: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			Padding(0, 1),

		NavItemHover: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(theme.Border).
			Padding(0, 1),

		// 二级导航
		Subnav: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Border).
			BorderBottom(true).
			Padding(0, 1),

		SubnavFocused: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Primary).
			BorderBottom(true).
			Padding(0, 1),

		SubnavItem: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(0, 1),

		SubnavActive: lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Bold(true).
			Padding(0, 1),

		// 列表
		TopicList: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Border),

		TopicListFocused: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Primary),

		ListItem: lipgloss.NewStyle().
			Padding(0, 1),

		ListItemSelected: lipgloss.NewStyle().
			Background(theme.Border).
			Foreground(theme.Foreground).
			Padding(0, 1),

		ListItemTitle: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		ListItemMeta: lipgloss.NewStyle().
			Foreground(theme.Muted),

		// 帖子详情
		DetailContainer: lipgloss.NewStyle().
			Padding(1, 2),

		DetailHeader: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Padding(0, 1),

		Title: lipgloss.NewStyle().
			Foreground(theme.TitleColor).
			Bold(true),

		Author: lipgloss.NewStyle().
			Foreground(theme.AuthorColor),

		Node: lipgloss.NewStyle().
			Foreground(theme.NodeColor).
			Background(theme.NodeBg).
			Padding(0, 1),

		Time: lipgloss.NewStyle().
			Foreground(theme.TimeColor),

		Content: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Separator: lipgloss.NewStyle().
			Foreground(theme.Border),

		// 回复
		ReplyBox: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(0, 1).
			MarginBottom(1),

		ReplyHeader: lipgloss.NewStyle().
			Foreground(theme.Muted),

		ReplyContent: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Floor: lipgloss.NewStyle().
			Foreground(theme.Muted),

		OPBadge: lipgloss.NewStyle().
			Foreground(theme.OPColor).
			Bold(true),

		// 搜索
		SearchContainer: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Info).
			BorderBottom(true).
			Padding(0, 1),

		SearchInput: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		SearchPrompt: lipgloss.NewStyle().
			Foreground(theme.Info),

		MatchHighlight: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),

		MatchCount: lipgloss.NewStyle().
			Foreground(theme.Muted),

		// 通用
		FocusedBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Primary),

		NormalBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.Border),
	}
}

// DefaultStyles 默认样式
var DefaultStyles = NewStyles(DarkTheme)
