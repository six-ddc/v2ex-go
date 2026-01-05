package model

// Supplement 附言模型
type Supplement struct {
	Index        int    // 附言序号（第几条附言）
	Content      string // 附言内容（纯文本）
	ContentHTML  string // 附言内容（HTML）
	RelativeTime string // 相对时间（如 "11小时31分钟前"）
}
