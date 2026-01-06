package api

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/six-ddc/v2ex-tui/internal/model"
)

// GetCurrentUser 获取当前登录用户信息
func (c *Client) GetCurrentUser() (*model.User, error) {
	doc, err := c.Get("/")
	if err != nil {
		return nil, err
	}

	user := &model.User{}

	// 检查是否已登录 - 查找用户名链接
	userSel := doc.Find("#Top .tools a[href^='/member/']").First()
	if userSel.Length() > 0 {
		user.LoggedIn = true
		user.Name = strings.TrimSpace(userSel.Text())

		// 解析通知数
		notifySel := doc.Find("#Top .tools a[href='/notifications']")
		notifyText := notifySel.Text()
		if matches := regexp.MustCompile(`(\d+)`).FindStringSubmatch(notifyText); len(matches) > 1 {
			if notify, err := strconv.Atoi(matches[1]); err == nil {
				user.Notify = notify
			}
		}

		// 解析余额
		balanceSel := doc.Find("#Top .balance_area, #money")
		balanceText := balanceSel.Text()

		// 解析银币和铜币
		// 格式可能是 "1234 银" 或 "1234/5678"
		silverRegex := regexp.MustCompile(`(\d+)\s*(?:银|S)`)
		if matches := silverRegex.FindStringSubmatch(balanceText); len(matches) > 1 {
			if silver, err := strconv.Atoi(matches[1]); err == nil {
				user.Silver = silver
			}
		}

		bronzeRegex := regexp.MustCompile(`(\d+)\s*(?:铜|B)`)
		if matches := bronzeRegex.FindStringSubmatch(balanceText); len(matches) > 1 {
			if bronze, err := strconv.Atoi(matches[1]); err == nil {
				user.Bronze = bronze
			}
		}

		// 解析头像
		avatarSel := doc.Find("#Top .avatar")
		if src, exists := avatarSel.Attr("src"); exists {
			user.Avatar = src
		}
	}

	return user, nil
}

// Login 用户登录 (暂时简化实现，后续需要处理验证码等)
func (c *Client) Login(username, password string) (*model.User, error) {
	// V2EX 登录比较复杂，需要:
	// 1. 获取登录页面，提取表单字段名 (动态生成)
	// 2. 处理验证码
	// 3. 提交表单
	// 暂时返回未登录用户
	return &model.User{
		LoggedIn: false,
	}, nil
}

// Logout 退出登录
func (c *Client) Logout() error {
	// 清除 Cookie
	c.cookies = nil
	return nil
}
