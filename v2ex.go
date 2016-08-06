package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	requests "github.com/levigross/grequests"
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
)

var Session *requests.Session

func init() {
	Session = requests.NewSession(nil)
}

func Login(username, password string) error {

	resp, err := Session.Get("https://www.v2ex.com/signin", &requests.RequestOptions{
		UserAgent: UserAgent,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("resp.StatusCode=%d", resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return err
	}

	formData := make(map[string]string)

	doc.Find("input.sl").Each(func(i int, s *goquery.Selection) {
		if val, _ := s.Attr("type"); val == "text" {
			val, _ = s.Attr("name")
			formData[val] = username
		} else if val, _ := s.Attr("type"); val == "password" {
			val, _ = s.Attr("name")
			formData[val] = password
		}
	})

	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, h1 := s.Attr("name")
		value, h2 := s.Attr("value")
		if h1 && h2 && len(value) > 0 {
			formData[name] = value
		}
	})

	Headers := make(map[string]string)
	Headers["Referer"] = "https://www.v2ex.com/signin"
	Headers["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err = Session.Post("https://www.v2ex.com/signin", &requests.RequestOptions{
		Data:      formData,
		UserAgent: UserAgent,
		Headers:   Headers,
	})

	return err
}

func ParseTopicByTab(tab string, uInfo *UserInfo, tabList [][]string) (ret []TopicInfo) {
	url := fmt.Sprintf("https://www.v2ex.com/?tab=%s", tab)
	resp, err := Session.Get(url, &requests.RequestOptions{
		UserAgent: UserAgent,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer log.Println(url, "status_code", resp.StatusCode)
	if resp.StatusCode != 200 {
		return
	}
	doc, err := goquery.NewDocumentFromResponse(resp.RawResponse)
	if err != nil {
		log.Println(err)
		return
	}
	uInfo.Name = strings.TrimSpace(doc.Find("span.bigger a").Text())
	doc.Find("a.fade").Each(func(i int, s *goquery.Selection) {
		if v, has := s.Attr("href"); has && v == "/notifications" {
			uInfo.Notify = s.Text()
		}
	})
	sliverStr := doc.Find("a.balance_area").Text()
	sliverLst := strings.Split(sliverStr, " ")
	setSli := false
	for _, sli := range sliverLst {
		if len(sli) > 0 {
			if !setSli {
				uInfo.Silver, _ = strconv.Atoi(sli)
				setSli = true
			} else {
				uInfo.Bronze, _ = strconv.Atoi(sli)
				break
			}
		}
	}
	log.Println("UserInfo", uInfo)
	doc.Find("div.box div.cell").Each(func(i int, s *goquery.Selection) {
		if s.Next().HasClass("cell item") && s.Prev().HasClass("inner") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				hrefSplit := strings.Split(href, "/")
				if hrefSplit[len(hrefSplit)-2] == "go" {
					href = hrefSplit[len(hrefSplit)-1]
					tabList[1] = append(tabList[1], href)
					tabList[0] = append(tabList[0], s.Text())
				}
			})
		}
	})
	doc.Find("div.cell.item").Each(func(i int, s *goquery.Selection) {
		topic := TopicInfo{}
		title := s.Find(".item_title a")
		topic.Title = title.Text()
		topic.Title = strings.Replace(topic.Title, "[", "<", -1)
		topic.Title = strings.Replace(topic.Title, "]", ">", -1)
		topic.Url, _ = title.Attr("href")
		topic.Url = "https://www.v2ex.com" + topic.Url
		info := s.Find(".small.fade").Text()
		info = strings.Replace(info, " ", "", -1)
		info = strings.Replace(info, string([]rune{0xA0}), "", -1) // 替换&nbsp
		infoList := strings.Split(info, "•")
		topic.Node = s.Find("a.node").Text()
		topic.Author = infoList[1]
		if len(infoList) > 2 {
			topic.Time = infoList[2]
			if len(infoList) > 3 {
				topic.LastReply = infoList[3]
			}
		}
		replyCount := s.Find("a.count_livid").Text()
		if replyCount != "" {
			topic.ReplyCount, _ = strconv.Atoi(replyCount)
		}
		ret = append(ret, topic)
	})
	return
}

func ParseTopicByNode(node string, page int) (ret []TopicInfo) {
	var url string
	if page > 1 {
		url = fmt.Sprintf("https://www.v2ex.com/go/%s?p=%d", node, page)
	} else {
		url = fmt.Sprintf("https://www.v2ex.com/go/%s", node)
	}
	resp, err := Session.Get(url, &requests.RequestOptions{
		UserAgent: UserAgent,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer log.Println(url, "status_code", resp.StatusCode)
	if resp.StatusCode != 200 {
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		log.Println(err)
		return
	}
	sel := doc.Find("div#TopicsNode")
	// log.Println(sel.Size())
	sel.Find("div.cell").Each(func(i int, s *goquery.Selection) {
		info := s.Find(".small.fade").Text()
		info = strings.Replace(info, " ", "", -1)
		info = strings.Replace(info, string([]rune{0xA0}), "", -1) // 替换&nbsp
		if len(info) == 0 {
			return
		}
		topic := TopicInfo{}
		title := s.Find(".item_title a")
		topic.Title = title.Text()
		topic.Title = strings.Replace(topic.Title, "[", "<", -1)
		topic.Title = strings.Replace(topic.Title, "]", ">", -1)
		topic.Url, _ = title.Attr("href")
		topic.Url = "https://www.v2ex.com" + topic.Url
		infoList := strings.Split(info, "•")
		topic.Node = node
		topic.Author = infoList[0]
		if len(infoList) > 2 {
			topic.Time = infoList[1]
			if len(infoList) > 3 {
				topic.LastReply = infoList[2]
			}
		}
		replyCount := s.Find("a.count_livid").Text()
		if replyCount != "" {
			topic.ReplyCount, _ = strconv.Atoi(replyCount)
		}
		ret = append(ret, topic)
		// log.Println(s.Find(".small.fade a.node").Text())
		// log.Println(s.Find(".small.fade strong a").Text())
	})
	return
}

func ParseReply(url string, reply *ReplyList) error {
	*reply = ReplyList{}
	resp, err := Session.Get(url, &requests.RequestOptions{
		UserAgent: UserAgent,
	})
	if err != nil {
		return err
	}
	defer log.Println(url, "status_code", resp.StatusCode)
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("resp.StatusCode=%d", resp.StatusCode))
	}
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return err
	}

	head := doc.Find("div.header small.gray").Text()
	head = strings.Replace(head, string([]rune{0xA0}), "", -1)
	head = strings.Replace(head, " ", "", -1)
	headList := strings.Split(head, "·")
	if len(headList) != 3 {
		// 这里好像是部分帖子需要登录才能浏览
		// https://www.v2ex.com/t/297344#reply12
		return errors.New("maybe need to login...")
	}
	reply.Lz = headList[0]
	reply.PostTime = headList[1]
	reply.ClickNum = headList[2]

	doc.Find("div.topic_content").Each(func(i int, sel *goquery.Selection) {
		selMD := sel.Find("div.markdown_body")
		contentList := []string{}
		if selMD.Size() > 0 {
			sel = selMD.Children()
			for _, node := range sel.Nodes {
				if node.Type == html.ElementNode {
					if node.Data == "p" {
						cnode := node.FirstChild
						pList := []string{}
						for cnode != nil {
							if cnode.Type == html.TextNode {
								pList = append(pList, cnode.Data)
							} else if cnode.Type == html.ElementNode {
								if cnode.Data == "a" {
									var href, text string
									for _, attr := range cnode.Attr {
										if attr.Key == "href" {
											href = attr.Val
										}
									}
									if len(href) > 0 && cnode.FirstChild != nil {
										text = cnode.FirstChild.Data
									}
									if href == text {
										pList = append(pList, href)
									} else {
										pList = append(pList, fmt.Sprintf("<%s>(%s)", text, href))
									}
								} else if cnode.Data == "img" {
									for _, attr := range cnode.Attr {
										if attr.Key == "src" {
											pList = append(pList, attr.Val)
										}
									}
								} else if cnode.Data == "strong" {
									pList = append(pList, cnode.FirstChild.Data)
								}
							}
							cnode = cnode.NextSibling
						}
						contentList = append(contentList, strings.Join(pList, ""))
					} else if node.Data == "ol" || node.Data == "ul" {
						cnode := node.FirstChild
						idx := 1
						olList := []string{}
						for cnode != nil {
							if cnode.Data == "li" {
								if node.Data == "ol" {
									olList = append(olList, fmt.Sprintf("%d. %s", idx, cnode.FirstChild.Data))
								} else {
									olList = append(olList, fmt.Sprintf("* %s", cnode.FirstChild.Data))
								}
								idx++
							}
							cnode = cnode.NextSibling
						}
						contentList = append(contentList, strings.Join(olList, "\n"))
					} else if node.Data == "img" {
						for _, attr := range node.Attr {
							if attr.Key == "src" {
								contentList = append(contentList, attr.Val)
							}
						}
					} else if node.Data[0] == 'h' && (node.Data[1] >= '1' && node.Data[1] <= '6') {
						var hstr string
						if node.FirstChild.Type == html.TextNode {
							hstr = node.FirstChild.Data
						} else if node.FirstChild.Type == html.ElementNode {
							if node.FirstChild.Data == "a" {
								var href, text string
								for _, attr := range node.FirstChild.Attr {
									if attr.Key == "href" {
										href = attr.Val
									}
								}
								if len(href) > 0 && node.FirstChild.FirstChild != nil {
									text = node.FirstChild.FirstChild.Data
								}
								if href == text {
									hstr = href
								} else {
									hstr = fmt.Sprintf("<%s>(%s)", text, href)
								}
							}
						}
						contentList = append(contentList, fmt.Sprintf("%s %s", strings.Repeat("#", int(node.Data[1]-'0')), hstr))
					}
				}
			}
			content := strings.Join(contentList, "\n\n")
			content = strings.Replace(content, "[", "<", -1)
			content = strings.Replace(content, "]", ">", -1)
			reply.Content = append(reply.Content, content)
		} else {
			cnode := sel.Nodes[0].FirstChild
			for cnode != nil {
				if cnode.Type == html.TextNode {
					contentList = append(contentList, cnode.Data)
				} else if cnode.Data == "img" {
					for _, attr := range cnode.Attr {
						if attr.Key == "src" {
							contentList = append(contentList, attr.Val)
						}
					}
				} else if cnode.Data == "a" {
					for _, attr := range cnode.Attr {
						if attr.Key == "href" {
							contentList = append(contentList, attr.Val)
							break
						}
					}
				}
				cnode = cnode.NextSibling
			}
			content := strings.Join(contentList, "")
			content = strings.Replace(content, "[", "<", -1)
			content = strings.Replace(content, "]", ">", -1)
			reply.Content = append(reply.Content, content)
		}
	})

	doc.Find("div#Main div.box").Find("div").Each(func(i int, sel *goquery.Selection) {
		idAttr, has := sel.Attr("id")
		if has && string(idAttr[:2]) == "r_" {
			replySel := sel.Find("div.reply_content")
			if len(replySel.Nodes) == 0 {
				return
			}
			info := ReplyInfo{}
			contentList := []string{}
			cnode := replySel.Nodes[0].FirstChild
			for cnode != nil {
				if cnode.Type == html.TextNode {
					contentList = append(contentList, cnode.Data)
				} else if cnode.Data == "img" {
					for _, attr := range cnode.Attr {
						if attr.Key == "src" {
							contentList = append(contentList, attr.Val)
						}
					}
				} else if cnode.Data == "a" {
					for _, attr := range cnode.Attr {
						if attr.Key == "href" {
							if strings.HasPrefix(attr.Val, "/member") {
								contentList = append(contentList, cnode.FirstChild.Data)
							} else {
								contentList = append(contentList, attr.Val)
							}
							break
						}
					}
				}
				cnode = cnode.NextSibling
			}
			info.Reply = strings.Join(contentList, "")
			info.Reply = strings.Replace(info.Reply, "[", "<", -1)
			info.Reply = strings.Replace(info.Reply, "]", ">", -1)
			info.Floor, _ = strconv.Atoi(sel.Find("span.no").Text())
			info.Member = sel.Find("a.dark").Text()
			info.Source = sel.Find("span.fade.small").Text()
			reply.List = append(reply.List, info)
			// log.Println(info.Floor, info.Member, info.Time)
		}
	})

	return nil
}
