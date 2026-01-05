# V2EX TUI 开发规范

## 颜色主题规范

### 核心原则

1. **禁止硬编码颜色值**：所有 UI 组件中的颜色必须使用 `ui.CurrentTheme` 中定义的主题色，不得直接使用 `lipgloss.Color("#XXXXXX")` 硬编码颜色。

2. **只设置前景色，不设置背景色**：Terminal 应用的标准做法是只设置前景色（文字颜色），让终端本身的背景色生效。这样用户可以通过终端配置自己喜欢的背景色（如 iTerm2、Terminal.app 等都支持主题配置）。参考 lazygit、htop、tig 等主流终端应用的做法。

### 背景色使用例外

以下情况允许设置背景色（这些是"语义上的高亮"，而非填充背景）：

1. **选中/高亮项**：列表中的选中行需要用背景色突出显示
2. **状态栏/工具栏**：StatusBar、HelpBar 等固定的 UI 结构元素
3. **按钮/标签**：有焦点的导航项、回复数标签等

### 主题颜色分类

| 类别 | 颜色名 | 用途 |
|-----|--------|------|
| **基础色** | `Primary` | 主色调（紫色），用于强调元素 |
| | `Secondary` | 次要色（绿色），用于次级导航 |
| | `Background` | 背景色（仅用于语义高亮） |
| | `Foreground` | 前景色/正文文字 |
| | `Muted` | 弱化色，用于次要文字 |
| | `Border` | 边框色/分隔线 |
| **语义色** | `Success` | 成功/在线状态 |
| | `Warning` | 警告 |
| | `Error` | 错误 |
| | `Info` | 信息/快捷键提示 |
| **特定元素** | `NodeColor` | 节点标签 |
| | `AuthorColor` | 作者名 |
| | `OPColor` | 楼主标记 |
| | `CountColor` | 回复数 |
| | `TimeColor` | 时间 |
| | `TitleColor` | 标题 |
| **选中状态** | `SelectedBg` / `SelectedFg` | 列表选中行背景/前景 |
| | `PrimaryFg` | Primary 背景上的前景色（始终白色） |

### 代码示例

```go
// 正确 - 只设置前景色
titleStyle := lipgloss.NewStyle().
    Foreground(ui.CurrentTheme.Foreground)

contentStyle := lipgloss.NewStyle().
    Foreground(ui.CurrentTheme.TitleColor).
    Bold(true)

// 正确 - 选中项使用背景色（语义高亮）
selectedStyle := lipgloss.NewStyle().
    Foreground(ui.CurrentTheme.SelectedFg).
    Background(ui.CurrentTheme.SelectedBg).
    Bold(true)

// 正确 - 状态栏使用背景色（UI 结构）
statusBarStyle := lipgloss.NewStyle().
    Foreground(ui.CurrentTheme.PrimaryFg).
    Background(ui.CurrentTheme.Primary)

// 错误 - 普通文本不应设置背景色
titleStyle := lipgloss.NewStyle().
    Foreground(ui.CurrentTheme.Foreground).
    Background(ui.CurrentTheme.Background)  // 不要这样做

// 错误 - 硬编码颜色
titleStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#E5E5E5"))   // 不要这样做
```

### lipgloss.Place 使用

不要使用 `WithWhitespaceBackground` 填充背景：

```go
// 正确 - 不填充背景色
content := lipgloss.Place(
    width, height,
    lipgloss.Left, lipgloss.Top,
    view,
)

// 错误 - 不要使用 WithWhitespaceBackground
bgOpt := lipgloss.WithWhitespaceBackground(ui.CurrentTheme.Background)
content := lipgloss.Place(width, height, lipgloss.Left, lipgloss.Top, view, bgOpt)
```

### 添加新颜色

如需新颜色，请在 `internal/ui/theme.go` 的 `Theme` 结构体中添加，并同时更新 `DarkTheme` 和 `LightTheme` 两个主题的值。

## 快捷键

- `t` - 切换深色/亮色主题
