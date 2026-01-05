package model

// User 用户模型
type User struct {
	Name     string // 用户名
	Avatar   string // 头像 URL
	Notify   int    // 未读通知数
	Silver   int    // 银币
	Bronze   int    // 铜币
	LoggedIn bool   // 是否已登录
}
