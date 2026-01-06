# V2EX TUI

一个现代化的 V2EX 终端客户端，基于 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 框架构建。

## 特性

- Vim 风格的键盘导航 (hjkl)
- 顶部双层导航设计，参考 V2EX 官网布局
- 支持深色/浅色主题切换
- 响应式布局，自适应终端窗口
- 终端超链接支持 (OSC 8)
- 无限滚动加载回复

## 安装

```bash
go install github.com/six-ddc/v2ex-tui@latest
```

或从源码构建：

```bash
git clone https://github.com/six-ddc/v2ex-tui.git
cd v2ex-tui
go build -o v2ex-tui
```

## 使用

```bash
./v2ex-tui
```

## 快捷键

### 全局

| 快捷键 | 功能 |
|--------|------|
| `q` / `Esc` | 返回/退出 |
| `Tab` | 切换焦点区域 |
| `Shift+Tab` | 反向切换焦点 |
| `r` | 刷新当前视图 |
| `t` | 切换深色/浅色主题 |
| `?` | 显示帮助 |

### 导航栏

| 快捷键 | 功能 |
|--------|------|
| `h` / `←` | 向左切换 |
| `l` / `→` | 向右切换 |
| `Enter` | 选中当前项 |
| `1-9, 0` | 快速跳转到对应 Tab |

### 主题列表

| 快捷键 | 功能 |
|--------|------|
| `j` / `↓` | 向下移动 |
| `k` / `↑` | 向上移动 |
| `Enter` / `l` | 打开帖子 |
| `g` | 跳到顶部 |
| `G` | 跳到底部 |
| `Ctrl+D` | 向下翻半页 |
| `Ctrl+U` | 向上翻半页 |
| `Ctrl+F` / `Space` | 向下翻页 |
| `Ctrl+B` | 向上翻页 |

### 帖子详情

| 快捷键 | 功能 |
|--------|------|
| `j/k` | 滚动内容 |
| `g/G` | 跳到顶部/底部 |
| `n` | 加载更多回复 |
| `[/]` | 上一篇/下一篇帖子 |
| `o` | 在浏览器中打开 |
| `q` | 返回列表 |

## 依赖

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI 框架
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - 样式库
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML 解析
- [resty](https://github.com/go-resty/resty) - HTTP 客户端

## 许可证

MIT License
