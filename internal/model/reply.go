package model

import "time"

// Reply 回复模型
type Reply struct {
	ID           string    // 回复 ID
	Floor        int       // 楼层号
	Author       string    // 作者用户名
	AuthorURL    string    // 作者主页 URL
	Content      string    // 回复内容 (纯文本或 Markdown)
	ContentHTML  string    // 原始 HTML 内容
	IsOP         bool      // 是否为楼主
	Time         time.Time // 回复时间
	RelativeTime string    // 相对时间 (如 "1小时前")
	Likes        int       // 感谢数
	Platform     string    // 来源平台
}
