package ui

import (
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

var (
	htmlConverter *md.Converter
)

// Hyperlink 创建终端可点击超链接 (OSC 8)
// 支持 iTerm2, GNOME Terminal, Windows Terminal 等
func Hyperlink(url, text string) string {
	// 补全相对路径
	if strings.HasPrefix(url, "/") {
		url = "https://www.v2ex.com" + url
	}
	// OSC 8 格式: \x1b]8;;URL\x07TEXT\x1b]8;;\x07
	return fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", url, text)
}

func init() {
	// 初始化 HTML to Markdown 转换器，禁用转义以避免 \. \- 等
	htmlConverter = md.NewConverter("", true, &md.Options{
		EscapeMode: "disabled",
	})

	// 自定义链接处理规则
	htmlConverter.AddRules(
		md.Rule{
			Filter: []string{"a"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				href, _ := selec.Attr("href")
				if href == "" {
					return &content
				}

				// 创建可点击超链接
				result := Hyperlink(href, content)
				return &result
			},
		},
		md.Rule{
			Filter: []string{"img"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				src, _ := selec.Attr("src")
				alt, _ := selec.Attr("alt")
				if src == "" {
					return nil
				}

				// 图片显示为可点击的 [图片] 链接
				text := "🖼️  图片"
				if alt != "" {
					text = "🖼️  " + alt
				}
				result := Hyperlink(src, text)
				return &result
			},
		},
	)
}

// RenderHTML 将 HTML 内容渲染为终端格式
// 保留 OSC 8 超链接
func RenderHTML(html string) string {
	if html == "" {
		return ""
	}

	result, err := htmlConverter.ConvertString(html)
	if err != nil {
		return html
	}

	return strings.TrimSpace(result)
}
