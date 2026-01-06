package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Theme 颜色主题
type Theme struct {
	// 基础色
	Primary    lipgloss.Color // 主色调
	Secondary  lipgloss.Color // 次要色
	Background lipgloss.Color // 背景色
	Foreground lipgloss.Color // 前景色
	Muted      lipgloss.Color // 弱化色
	Border     lipgloss.Color // 边框色

	// 语义色
	Success lipgloss.Color // 成功/在线
	Warning lipgloss.Color // 警告
	Error   lipgloss.Color // 错误
	Info    lipgloss.Color // 信息

	// 特定元素
	NodeColor   lipgloss.Color // 节点标签
	NodeBg      lipgloss.Color // 节点背景
	AuthorColor lipgloss.Color // 作者名
	OPColor     lipgloss.Color // 楼主标记
	CountColor  lipgloss.Color // 回复数
	TimeColor   lipgloss.Color // 时间
	LinkColor   lipgloss.Color // 链接
	TitleColor  lipgloss.Color // 标题

	// 选中状态
	SelectedBg       lipgloss.Color // 选中行背景
	SelectedFg       lipgloss.Color // 选中行前景
	PrimaryFg        lipgloss.Color // Primary背景上的前景色（始终白色）
	ReplyCountBg     lipgloss.Color // 回复数背景
	ReplyCountBgSel  lipgloss.Color // 选中时回复数背景
	HeaderBg         lipgloss.Color // 顶部栏背景
	HeaderFg         lipgloss.Color // 顶部栏前景
}

// DarkTheme 深色主题
var DarkTheme = Theme{
	Primary:    lipgloss.Color("#7C3AED"), // 紫色
	Secondary:  lipgloss.Color("#10B981"), // 绿色
	Background: lipgloss.Color("#1F2937"),
	Foreground: lipgloss.Color("#F9FAFB"),
	Muted:      lipgloss.Color("#6B7280"),
	Border:     lipgloss.Color("#374151"),

	Success: lipgloss.Color("#10B981"),
	Warning: lipgloss.Color("#F59E0B"),
	Error:   lipgloss.Color("#EF4444"),
	Info:    lipgloss.Color("#3B82F6"),

	NodeColor:   lipgloss.Color("#8B5CF6"),
	NodeBg:      lipgloss.Color("#2D2D2D"),
	AuthorColor: lipgloss.Color("#60A5FA"),
	OPColor:     lipgloss.Color("#F472B6"),
	CountColor:  lipgloss.Color("#34D399"),
	TimeColor:   lipgloss.Color("#9CA3AF"),
	LinkColor:   lipgloss.Color("#38BDF8"),
	TitleColor:  lipgloss.Color("#FBBF24"),

	SelectedBg:      lipgloss.Color("#374151"),
	SelectedFg:      lipgloss.Color("#FFFFFF"),
	PrimaryFg:       lipgloss.Color("#FFFFFF"),
	ReplyCountBg:    lipgloss.Color("#065F46"),
	ReplyCountBgSel: lipgloss.Color("#047857"),
	HeaderBg:        lipgloss.Color("#374151"),
	HeaderFg:        lipgloss.Color("#D1D5DB"),
}

// LightTheme 浅色主题
var LightTheme = Theme{
	Primary:    lipgloss.Color("#7C3AED"),
	Secondary:  lipgloss.Color("#059669"),
	Background: lipgloss.Color("#FFFFFF"),
	Foreground: lipgloss.Color("#111827"),
	Muted:      lipgloss.Color("#6B7280"),
	Border:     lipgloss.Color("#E5E7EB"),

	Success: lipgloss.Color("#059669"),
	Warning: lipgloss.Color("#D97706"),
	Error:   lipgloss.Color("#DC2626"),
	Info:    lipgloss.Color("#2563EB"),

	NodeColor:   lipgloss.Color("#7C3AED"),
	NodeBg:      lipgloss.Color("#F3F4F6"),
	AuthorColor: lipgloss.Color("#2563EB"),
	OPColor:     lipgloss.Color("#DB2777"),
	CountColor:  lipgloss.Color("#059669"),
	TimeColor:   lipgloss.Color("#4B5563"),
	LinkColor:   lipgloss.Color("#0284C7"),
	TitleColor:  lipgloss.Color("#92400E"),

	SelectedBg:      lipgloss.Color("#E5E7EB"),
	SelectedFg:      lipgloss.Color("#000000"),
	PrimaryFg:       lipgloss.Color("#FFFFFF"),
	ReplyCountBg:    lipgloss.Color("#D1FAE5"),
	ReplyCountBgSel: lipgloss.Color("#A7F3D0"),
	HeaderBg:        lipgloss.Color("#E5E7EB"),
	HeaderFg:        lipgloss.Color("#1F2937"),
}

// CurrentTheme 当前使用的主题
var CurrentTheme Theme

// IsDarkTheme 是否是深色主题
var IsDarkTheme bool

func init() {
	// 自动检测终端背景色
	output := termenv.NewOutput(os.Stdout)
	if output.HasDarkBackground() {
		CurrentTheme = DarkTheme
		IsDarkTheme = true
	} else {
		CurrentTheme = LightTheme
		IsDarkTheme = false
	}
}

// ToggleTheme 切换主题
func ToggleTheme() {
	if IsDarkTheme {
		CurrentTheme = LightTheme
		IsDarkTheme = false
	} else {
		CurrentTheme = DarkTheme
		IsDarkTheme = true
	}
}
