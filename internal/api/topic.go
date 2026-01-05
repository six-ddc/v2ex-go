package api

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/six-ddc/v2ex-go/internal/model"
)

// 预编译正则表达式，避免重复编译
var (
	topicIDRegex    = regexp.MustCompile(`/t/(\d+)`)
	numberRegex     = regexp.MustCompile(`(\d+)`)
	replyCountRegex = regexp.MustCompile(`(\d+)\s*条回复`)
)

// GetTopicsByTab 按 Tab 获取主题列表，同时返回二级节点列表和用户信息
func (c *Client) GetTopicsByTab(tab string) ([]model.Topic, []model.Node, *model.User, error) {
	path := "/?tab=" + tab
	doc, err := c.Get(path)
	if err != nil {
		return nil, nil, nil, err
	}

	topics := parseTopicListByTab(doc)
	nodes := parseSubNodes(doc)
	user := parseUserInfo(doc)

	return topics, nodes, user, nil
}

// GetTopicsByNode 按节点获取主题列表
func (c *Client) GetTopicsByNode(nodeCode string, page int) ([]model.Topic, error) {
	path := "/go/" + nodeCode
	if page > 1 {
		path += "?p=" + strconv.Itoa(page)
	}

	doc, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	return parseTopicListByNode(doc, nodeCode), nil
}

// GetTopicDetail 获取帖子详情和回复
func (c *Client) GetTopicDetail(topicURL string) (*model.Topic, []model.Reply, error) {
	return c.GetTopicDetailPage(topicURL, 1)
}

// GetTopicDetailPage 获取帖子详情和指定页码的回复
func (c *Client) GetTopicDetailPage(topicURL string, page int) (*model.Topic, []model.Reply, error) {
	// 从完整 URL 提取路径部分（如 /t/1183189）
	path := topicURL
	if strings.HasPrefix(path, "https://www.v2ex.com") {
		path = strings.TrimPrefix(path, "https://www.v2ex.com")
	} else if strings.HasPrefix(path, "http://www.v2ex.com") {
		path = strings.TrimPrefix(path, "http://www.v2ex.com")
	}

	// 去掉锚点部分（如 #reply123）
	if idx := strings.Index(path, "#"); idx != -1 {
		path = path[:idx]
	}

	// 构造带分页的路径
	if page > 1 {
		if strings.Contains(path, "?") {
			path += "&p=" + strconv.Itoa(page)
		} else {
			path += "?p=" + strconv.Itoa(page)
		}
	}

	doc, err := c.Get(path)
	if err != nil {
		return nil, nil, err
	}

	topic := parseTopicDetail(doc, topicURL)
	topic.CurrentPage = page
	replies := parseReplies(doc)

	return topic, replies, nil
}

// parseUserInfo 解析用户信息
func parseUserInfo(doc *goquery.Document) *model.User {
	user := &model.User{}

	// 解析用户名
	user.Name = strings.TrimSpace(doc.Find("span.bigger a").Text())
	if user.Name != "" {
		user.LoggedIn = true
	}

	// 解析通知数
	doc.Find("a.fade").Each(func(i int, s *goquery.Selection) {
		if href, has := s.Attr("href"); has && href == "/notifications" {
			user.Notify, _ = strconv.Atoi(strings.TrimSpace(s.Text()))
		}
	})

	// 解析银币和铜币
	balanceText := doc.Find("a.balance_area").Text()
	balanceList := strings.Split(balanceText, " ")
	setSilver := false
	for _, b := range balanceList {
		b = strings.TrimSpace(b)
		if len(b) > 0 {
			if !setSilver {
				user.Silver, _ = strconv.Atoi(b)
				setSilver = true
			} else {
				user.Bronze, _ = strconv.Atoi(b)
				break
			}
		}
	}

	return user
}

// parseTopicListByTab 解析 Tab 页面的主题列表
func parseTopicListByTab(doc *goquery.Document) []model.Topic {
	var topics []model.Topic

	doc.Find("div.cell.item").Each(func(i int, s *goquery.Selection) {
		topic := model.Topic{}

		// 解析标题和链接
		titleSel := s.Find(".item_title a")
		topic.Title = strings.TrimSpace(titleSel.Text())
		topic.Title = strings.ReplaceAll(topic.Title, "[", "<")
		topic.Title = strings.ReplaceAll(topic.Title, "]", ">")

		if href, exists := titleSel.Attr("href"); exists {
			topic.URL = "https://www.v2ex.com" + href
			// 从 URL 提取 ID
			if matches := topicIDRegex.FindStringSubmatch(href); len(matches) > 1 {
				topic.ID = matches[1]
			}
		}

		// 解析节点
		topic.Node = model.Node{
			Name: strings.TrimSpace(s.Find("a.node").Text()),
		}
		if href, exists := s.Find("a.node").Attr("href"); exists {
			topic.Node.Code = strings.TrimPrefix(href, "/go/")
		}

		// 解析 topic_info
		infoText := s.Find(".topic_info").Text()
		infoText = strings.ReplaceAll(infoText, " ", "")
		infoText = strings.ReplaceAll(infoText, string(rune(0xA0)), "") // 替换 &nbsp;
		infoList := strings.Split(infoText, "•")

		if len(infoList) > 1 {
			topic.Author = infoList[1]
		}
		if len(infoList) > 2 {
			topic.RelativeTime = infoList[2]
		}
		if len(infoList) > 3 {
			topic.LastReplyBy = infoList[3]
		}

		// 解析回复数
		replyCount := strings.TrimSpace(s.Find("a.count_livid").Text())
		if replyCount != "" {
			topic.ReplyCount, _ = strconv.Atoi(replyCount)
		}

		if topic.Title != "" {
			topics = append(topics, topic)
		}
	})

	return topics
}

// parseTopicListByNode 解析节点页面的主题列表
func parseTopicListByNode(doc *goquery.Document, nodeCode string) []model.Topic {
	var topics []model.Topic

	doc.Find("div#TopicsNode div.cell").Each(func(i int, s *goquery.Selection) {
		infoText := s.Find(".small.fade").Text()
		infoText = strings.ReplaceAll(infoText, " ", "")
		infoText = strings.ReplaceAll(infoText, string(rune(0xA0)), "")
		if len(infoText) == 0 {
			return
		}

		topic := model.Topic{}

		// 解析标题和链接
		titleSel := s.Find(".item_title a")
		topic.Title = strings.TrimSpace(titleSel.Text())
		topic.Title = strings.ReplaceAll(topic.Title, "[", "<")
		topic.Title = strings.ReplaceAll(topic.Title, "]", ">")

		if href, exists := titleSel.Attr("href"); exists {
			topic.URL = "https://www.v2ex.com" + href
			if matches := topicIDRegex.FindStringSubmatch(href); len(matches) > 1 {
				topic.ID = matches[1]
			}
		}

		topic.Node = model.Node{
			Name: nodeCode,
			Code: nodeCode,
		}

		infoList := strings.Split(infoText, "•")
		if len(infoList) > 0 {
			topic.Author = infoList[0]
		}
		if len(infoList) > 1 {
			topic.RelativeTime = infoList[1]
		}
		if len(infoList) > 2 {
			topic.LastReplyBy = infoList[2]
		}

		replyCount := strings.TrimSpace(s.Find("a.count_livid").Text())
		if replyCount != "" {
			topic.ReplyCount, _ = strconv.Atoi(replyCount)
		}

		if topic.Title != "" {
			topics = append(topics, topic)
		}
	})

	return topics
}

// parseSubNodes 解析二级节点列表
func parseSubNodes(doc *goquery.Document) []model.Node {
	var nodes []model.Node

	// 查找 Tab 下方的节点列表
	doc.Find("div.box div.cell").Each(func(i int, s *goquery.Selection) {
		// 查找紧邻 cell item 之前且紧邻 inner 之后的 cell
		if s.Next().HasClass("cell") && s.Next().HasClass("item") && s.Prev().HasClass("inner") {
			s.Find("a").Each(func(j int, a *goquery.Selection) {
				href, _ := a.Attr("href")
				hrefSplit := strings.Split(href, "/")
				if len(hrefSplit) >= 2 && hrefSplit[len(hrefSplit)-2] == "go" {
					node := model.Node{
						Name: strings.TrimSpace(a.Text()),
						Code: hrefSplit[len(hrefSplit)-1],
					}
					nodes = append(nodes, node)
				}
			})
		}
	})

	return nodes
}

// parseTopicDetail 解析帖子详情
func parseTopicDetail(doc *goquery.Document, topicURL string) *model.Topic {
	topic := &model.Topic{
		URL:        topicURL,
		TotalPages: 1, // 默认1页
	}

	// 从 URL 提取 ID
	if matches := topicIDRegex.FindStringSubmatch(topicURL); len(matches) > 1 {
		topic.ID = matches[1]
	}

	// 解析标题
	topic.Title = strings.TrimSpace(doc.Find("h1").First().Text())

	// 解析 header 信息
	headerText := doc.Find("div.header small.gray").Text()
	headerText = strings.ReplaceAll(headerText, string(rune(0xA0)), "")
	headerText = strings.ReplaceAll(headerText, " ", "")
	headerList := strings.Split(headerText, "·")

	if len(headerList) >= 1 {
		topic.Author = headerList[0]
	}
	if len(headerList) >= 2 {
		topic.RelativeTime = headerList[1]
	}
	if len(headerList) >= 3 {
		// 解析点击数
		clickText := headerList[2]
		if matches := numberRegex.FindStringSubmatch(clickText); len(matches) > 1 {
			topic.Clicks, _ = strconv.Atoi(matches[1])
		}
	}

	// 解析节点
	nodeSel := doc.Find(".header a[href^='/go/']").First()
	topic.Node = model.Node{
		Name: strings.TrimSpace(nodeSel.Text()),
	}
	if href, exists := nodeSel.Attr("href"); exists {
		topic.Node.Code = strings.TrimPrefix(href, "/go/")
	}

	// 解析正文内容 - 在主内容区的 box 中找第一个 cell 下的 topic_content
	// 注意：附言(subtle)中也有 topic_content，所以要排除 subtle
	doc.Find("div#Main div.box div.cell > div.topic_content").First().Each(func(i int, sel *goquery.Selection) {
		topic.ContentHTML, _ = sel.Html()

		// 检查是否有 markdown_body
		selMD := sel.Find("div.markdown_body")
		if selMD.Length() > 0 {
			topic.Content = parseMarkdownContent(selMD)
		} else {
			topic.Content = parseSimpleContent(sel)
		}
	})

	// 解析附言（div.subtle）
	doc.Find("div.subtle").Each(func(i int, sel *goquery.Selection) {
		supplement := model.Supplement{
			Index: i + 1,
		}

		// 解析附言时间信息
		fadeText := sel.Find("span.fade").Text()
		fadeText = strings.ReplaceAll(fadeText, string(rune(0xA0)), " ")
		// 格式："第 1 条附言  ·  11 小时 31 分钟前"
		if parts := strings.Split(fadeText, "·"); len(parts) >= 2 {
			supplement.RelativeTime = strings.TrimSpace(parts[1])
		}

		// 解析附言内容
		contentSel := sel.Find("div.topic_content")
		supplement.ContentHTML, _ = contentSel.Html()
		supplement.Content = strings.TrimSpace(contentSel.Text())

		topic.Supplements = append(topic.Supplements, supplement)
	})

	// 解析回复数
	doc.Find(".cell .gray").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if matches := replyCountRegex.FindStringSubmatch(text); len(matches) > 1 {
			topic.ReplyCount, _ = strconv.Atoi(matches[1])
		}
	})

	// 解析分页信息
	doc.Find("a.page_normal, a.page_current").Each(func(i int, s *goquery.Selection) {
		pageNum, _ := strconv.Atoi(strings.TrimSpace(s.Text()))
		if pageNum > topic.TotalPages {
			topic.TotalPages = pageNum
		}
	})

	return topic
}

// parseMarkdownContent 解析 Markdown 格式的内容
func parseMarkdownContent(sel *goquery.Selection) string {
	var contentList []string

	sel.Children().Each(func(i int, child *goquery.Selection) {
		node := child.Get(0)
		if node.Type == html.ElementNode {
			switch node.Data {
			case "p":
				contentList = append(contentList, parseParagraph(node))
			case "ol", "ul":
				contentList = append(contentList, parseList(node))
			case "img":
				if src := getAttr(node, "src"); src != "" {
					contentList = append(contentList, src)
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				level := int(node.Data[1] - '0')
				text := parseHeading(node)
				contentList = append(contentList, strings.Repeat("#", level)+" "+text)
			case "pre":
				// 代码块
				code := child.Find("code").Text()
				contentList = append(contentList, "```\n"+code+"\n```")
			case "blockquote":
				text := strings.TrimSpace(child.Text())
				contentList = append(contentList, "> "+text)
			}
		}
	})

	content := strings.Join(contentList, "\n\n")
	content = strings.ReplaceAll(content, "[", "<")
	content = strings.ReplaceAll(content, "]", ">")
	return content
}

// parseSimpleContent 解析简单文本内容
func parseSimpleContent(sel *goquery.Selection) string {
	var contentList []string

	if len(sel.Nodes) == 0 {
		return ""
	}

	cnode := sel.Nodes[0].FirstChild
	for cnode != nil {
		switch {
		case cnode.Type == html.TextNode:
			contentList = append(contentList, cnode.Data)
		case cnode.Data == "img":
			if src := getAttr(cnode, "src"); src != "" {
				contentList = append(contentList, src)
			}
		case cnode.Data == "a":
			if href := getAttr(cnode, "href"); href != "" {
				contentList = append(contentList, href)
			}
		case cnode.Data == "br":
			contentList = append(contentList, "\n")
		}
		cnode = cnode.NextSibling
	}

	content := strings.Join(contentList, "")
	content = strings.ReplaceAll(content, "[", "<")
	content = strings.ReplaceAll(content, "]", ">")
	return content
}

// parseParagraph 解析段落
func parseParagraph(node *html.Node) string {
	var parts []string
	cnode := node.FirstChild
	for cnode != nil {
		switch {
		case cnode.Type == html.TextNode:
			parts = append(parts, cnode.Data)
		case cnode.Type == html.ElementNode:
			switch cnode.Data {
			case "a":
				href := getAttr(cnode, "href")
				text := getInnerText(cnode)
				if href == text {
					parts = append(parts, href)
				} else {
					parts = append(parts, "<"+text+">("+href+")")
				}
			case "img":
				if src := getAttr(cnode, "src"); src != "" {
					parts = append(parts, src)
				}
			case "strong", "b":
				parts = append(parts, getInnerText(cnode))
			case "em", "i":
				parts = append(parts, getInnerText(cnode))
			case "code":
				parts = append(parts, "`"+getInnerText(cnode)+"`")
			case "br":
				parts = append(parts, "\n")
			}
		}
		cnode = cnode.NextSibling
	}
	return strings.Join(parts, "")
}

// parseList 解析列表
func parseList(node *html.Node) string {
	var items []string
	idx := 1
	cnode := node.FirstChild
	for cnode != nil {
		if cnode.Data == "li" {
			text := getInnerText(cnode)
			if node.Data == "ol" {
				items = append(items, strconv.Itoa(idx)+". "+text)
			} else {
				items = append(items, "* "+text)
			}
			idx++
		}
		cnode = cnode.NextSibling
	}
	return strings.Join(items, "\n")
}

// parseHeading 解析标题
func parseHeading(node *html.Node) string {
	if node.FirstChild == nil {
		return ""
	}
	if node.FirstChild.Type == html.TextNode {
		return node.FirstChild.Data
	}
	if node.FirstChild.Type == html.ElementNode && node.FirstChild.Data == "a" {
		href := getAttr(node.FirstChild, "href")
		text := getInnerText(node.FirstChild)
		if href == text {
			return href
		}
		return "<" + text + ">(" + href + ")"
	}
	return ""
}

// parseReplies 解析回复列表
func parseReplies(doc *goquery.Document) []model.Reply {
	var replies []model.Reply

	// 获取楼主用户名
	opName := strings.TrimSpace(doc.Find("div.header small.gray").Text())
	if parts := strings.Split(opName, "·"); len(parts) > 0 {
		opName = strings.ReplaceAll(parts[0], string(rune(0xA0)), "")
		opName = strings.TrimSpace(opName)
	}

	doc.Find("div#Main div.box div[id^='r_']").Each(func(i int, sel *goquery.Selection) {
		reply := model.Reply{}

		// 解析回复 ID
		if id, exists := sel.Attr("id"); exists {
			reply.ID = strings.TrimPrefix(id, "r_")
		}

		// 解析回复内容
		replySel := sel.Find("div.reply_content")
		if replySel.Length() == 0 {
			return
		}

		reply.ContentHTML, _ = replySel.Html()

		var contentList []string
		if len(replySel.Nodes) > 0 {
			cnode := replySel.Nodes[0].FirstChild
			for cnode != nil {
				switch {
				case cnode.Type == html.TextNode:
					contentList = append(contentList, cnode.Data)
				case cnode.Data == "img":
					if src := getAttr(cnode, "src"); src != "" {
						contentList = append(contentList, src)
					}
				case cnode.Data == "a":
					href := getAttr(cnode, "href")
					if strings.HasPrefix(href, "/member") {
						contentList = append(contentList, getInnerText(cnode))
					} else {
						contentList = append(contentList, href)
					}
				case cnode.Data == "br":
					contentList = append(contentList, "\n")
				}
				cnode = cnode.NextSibling
			}
		}
		reply.Content = strings.Join(contentList, "")
		reply.Content = strings.ReplaceAll(reply.Content, "[", "<")
		reply.Content = strings.ReplaceAll(reply.Content, "]", ">")

		// 解析楼层
		floorText := strings.TrimSpace(sel.Find("span.no").Text())
		reply.Floor, _ = strconv.Atoi(floorText)

		// 解析作者
		reply.Author = strings.TrimSpace(sel.Find("a.dark").Text())
		reply.IsOP = (reply.Author == opName)

		// 解析时间
		reply.RelativeTime = strings.TrimSpace(sel.Find("span.ago").Text())

		// 解析感谢数
		likesText := strings.TrimSpace(sel.Find("span.small.fade").Text())
		if likesText != "" {
			reply.Likes, _ = strconv.Atoi(likesText)
		}

		if reply.Content != "" {
			replies = append(replies, reply)
		}
	})

	return replies
}

// getAttr 获取节点属性
func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getInnerText 获取节点内部文本
func getInnerText(node *html.Node) string {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return node.FirstChild.Data
	}
	return ""
}
