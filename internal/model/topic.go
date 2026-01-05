package model

import "time"

// Topic 主题模型
type Topic struct {
	ID            string       // 主题 ID
	Title         string       // 标题
	URL           string       // 完整 URL
	Author        string       // 作者用户名
	AuthorURL     string       // 作者主页 URL
	Node          Node         // 所属节点
	CreatedAt     time.Time    // 创建时间
	LastReply     time.Time    // 最后回复时间
	LastReplyBy   string       // 最后回复者
	ReplyCount    int          // 回复数量
	TotalPages    int          // 回复总页数
	CurrentPage   int          // 当前回复页码
	Content       string       // 正文内容 (Markdown 格式)
	ContentHTML   string       // 原始 HTML 内容
	Clicks        int          // 点击数
	RelativeTime  string       // 相对时间 (如 "2小时前")
	Supplements   []Supplement // 附言列表
}
