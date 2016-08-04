package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	requests "github.com/levigross/grequests"
	rw "github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	ui "github.com/six-ddc/termui"
	// "io/ioutil"
	"errors"
	"log"
	"math/rand"
	// "net/http"
	// "net/http/cookiejar"
	// "net/url"
	"golang.org/x/net/html"
	"os"
	// "golang.org/x/net/publicsuffix"
	// "encoding/json"
	"strconv"
	"strings"
	"time"
)

var (
	uiTab   *UITab
	uiTopic *UITopicList
	uiReply *UIReplyList
	uiLog   *UILog
	uiUser  *UIUser

	TopicRows *ui.Row
	ReplyRows *ui.Row

	Session *requests.Session

	LastCtrlW     int64
	CurrState     State
	CurrBodyState State

	ShortKeys  []byte
	MatchList  []int
	MatchIndex int
)

type State int

const (
	StateDefault State = iota
	StateTab
	StateBody
	StateMax
)

const (
	BodyStateTopic State = StateMax + 1
	BodyStateReply State = StateMax + 2
)

type UserInfo struct {
	Name   string
	Notify string
	Silver int
	Bronze int
}

type ReplyInfo struct {
	Floor  int
	Member string
	Reply  string
	Time   string
	Up     string
}

type ReplyList struct {
	Title    string
	Url      string
	Content  []string
	Lz       string
	PostTime string
	ClickNum string
	List     []ReplyInfo
}

type TopicInfo struct {
	Title      string
	Url        string
	Author     string
	AuthorImg  string
	Node       string
	Time       string
	LastReply  string
	ReplyCount int
}

type TopicType uint16

const (
	TopicTab TopicType = iota
	TopicNode
)

type ScrollList struct {
	ui.List
	AllItems []string
	Index    int
}

func NewScrollList() *ScrollList {
	s := &ScrollList{List: *ui.NewList()}
	return s
}

func (l *ScrollList) SetItem(i int, item string) {
	l.Items[i] = item
}

func (l *ScrollList) SetAllItems(items []string, updateLastList bool) {
	if updateLastList {
		l.AllItems = items
	}
	sz := l.InnerHeight()
	if len(items) < sz {
		sz = len(items)
	}
	l.Items = make([]string, sz)
	copy(l.Items, items) // 复制长度以较小的slice为准
	l.ResetBgColor()
}

func (l *ScrollList) SetBgColor(i int, color ui.Attribute) {
	if len(l.ItemBg) != len(l.Items) {
		l.ResetBgColor()
	}
	l.ItemBg[i] = color
}

func (l *ScrollList) Highlight(b bool) {
	if b {
		l.BorderFg = ui.ColorRed
	} else {
		l.BorderFg = ui.ColorDefault
	}
}

func (l *ScrollList) ResetBgColor() {
	l.ItemBg = make([]ui.Attribute, len(l.Items))
	for i, _ := range l.ItemBg {
		l.ItemBg[i] = ui.ThemeAttr("list.item.bg")
	}
}

func (l *ScrollList) ScrollDown() {
	sz := len(l.AllItems)
	screen_heigth := l.InnerHeight()
	if sz > screen_heigth+l.Index {
		l.ResetBgColor()
		l.Index++
		l.SetAllItems(l.AllItems[l.Index:], false)
		ui.Render(l)
	}
}

func (l *ScrollList) PageDown() {
	sz := len(l.AllItems)
	screen_heigth := l.InnerHeight()
	if sz < screen_heigth {
		return
	}
	index := l.Index + screen_heigth
	if index > sz-screen_heigth {
		index = sz - screen_heigth
		if index == l.Index {
			return
		}
	}
	l.Index = index
	l.SetAllItems(l.AllItems[l.Index:], false)
	ui.Render(l)
}

func (l *ScrollList) PageUp() {
	screen_heigth := l.InnerHeight()
	index := l.Index - screen_heigth
	if index < 0 {
		index = 0
		if index == l.Index {
			return
		}
	}
	l.Index = index
	l.SetAllItems(l.AllItems[l.Index:], false)
	ui.Render(l)
}

func (l *ScrollList) ScrollUp() {
	if l.Index > 0 {
		l.ResetBgColor()
		l.Index--
		l.SetAllItems(l.AllItems[l.Index:], false)
		ui.Render(l)
	}
}

var userAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.82 Safari/537.36"

type UILog struct {
	ui.List
	Label string
	Index int
}

func NewLog() *UILog {
	l := &UILog{Index: 0, Label: "Log [C-l]"}
	l.List = *ui.NewList()
	l.Height = 5
	l.BorderLabel = l.Label
	l.Items = make([]string, l.Height-2)
	return l
}

func (l *UILog) Write(p []byte) (n int, err error) {
	str := fmt.Sprintf("[%d] %s", uiLog.Index+1, p)
	if uiLog.Items[len(uiLog.Items)-1] != "" {
		i := 0
		for ; i < len(uiLog.Items)-1; i++ {
			uiLog.Items[i] = uiLog.Items[i+1]
		}
		uiLog.Items[i] = str
	} else {
		for i, item := range uiLog.Items {
			if item == "" {
				uiLog.Items[i] = str
				break
			}
		}
	}
	uiLog.Index++
	ui.Render(l)
	return len(p), nil
}

type UIUser struct {
	ui.List
	User *UserInfo
}

func NewUser() *UIUser {
	u := &UIUser{User: &UserInfo{}}
	u.List = *ui.NewList()
	u.Height = 5
	u.Items = make([]string, u.Height-2)
	return u
}

func (u *UIUser) Fresh() {
	if len(u.User.Name) > 0 {
		u.Items[0] = fmt.Sprintf("[%s](fg-green)", u.User.Name)
		balance := fmt.Sprintf("%d/%d", u.User.Silver, u.User.Bronze)
		space_width := u.InnerWidth() - 1 - rw.StringWidth(u.User.Notify) - rw.StringWidth(balance)
		if space_width > 0 {
			u.Items[2] = fmt.Sprintf("%s%s%s", u.User.Notify, strings.Repeat(" ", space_width), balance)
		} else {
			u.Items[2] = fmt.Sprintf("%d %s", u.User.Notify, balance)
		}
	}
}

type UITab struct {
	ui.List
	Label        string
	NameList     [][]string
	ChildList    [][][]string
	CurrTab      int
	CurrChildTab int
}

func NewTab() *UITab {
	t := &UITab{Label: "Tab [C-t]", CurrTab: 9 /*(全部)*/, CurrChildTab: -1}
	t.NameList = [][]string{
		{"技术", "创意", "好玩", "Apple", "酷工作", "交易", "城市", "问与答", "最热", "全部", "R2", "节点", "关注"},
		{"tech", "creative", "play", "apple", "jobs", "deals", "city", "qna", "hot", "all", "r2", "nodes", "members"}}
	t.ChildList = make([][][]string, len(t.NameList[0]))
	t.List = *ui.NewList()
	t.BorderLabel = t.Label
	t.Height = 2 + 2
	t.Items = make([]string, 2)
	t.ResetTabList()
	return t
}

func (t *UITab) Highlight(b bool) {
	if b {
		t.BorderFg = ui.ColorRed
	} else {
		t.BorderFg = ui.ColorDefault
	}
}

func (t *UITab) ResetTabList() {
	strList := []string{}
	for i, names := range (t.NameList)[0] {
		if i == t.CurrTab {
			strList = append(strList, fmt.Sprintf("[%s](bg-red)", names))
		} else {
			strList = append(strList, names)
		}
	}
	t.Items[0] = strings.Join(strList, " ")
	childList := (t.ChildList)[t.CurrTab]
	if len(childList) == 0 {
		t.Items[1] = ""
		return
	}
	strList = strList[:0]
	for i, names := range childList[0] {
		if i == t.CurrChildTab {
			strList = append(strList, fmt.Sprintf("[%s](bg-red)", names))
		} else {
			strList = append(strList, names)
		}
	}
	t.Items[1] = strings.Join(strList, " ")
}

func (t *UITab) UpdateLabel() {
	str := t.Label
	if len(ShortKeys) > 0 {
		str = fmt.Sprintf("%s (%s)", str, ShortKeys)
	}
	t.BorderLabel = str
}

func (t *UITab) GetTabNode() (cate string, node string) {
	cate = t.NameList[1][t.CurrTab]
	if t.CurrChildTab >= 0 {
		childList := (t.ChildList)[t.CurrTab]
		if len(childList) > 0 {
			node = childList[1][t.CurrChildTab]
		}
	}
	return
}

func (t *UITab) MatchTab() {
	strList := []string{}
	count := 0
	MatchList = MatchList[:0]
	for i, names := range (t.NameList)[1] {
		str := matchKey([]byte(names), ShortKeys)
		if str != names {
			if count == MatchIndex {
				strList = append(strList, fmt.Sprintf("[%s](bg-red)<%s>", t.NameList[0][i], str))
			} else {
				strList = append(strList, fmt.Sprintf("[%s](bg-blue)<%s>", t.NameList[0][i], str))
			}
			count++
			MatchList = append(MatchList, i)
		} else {
			strList = append(strList, fmt.Sprintf("%s<%s>", t.NameList[0][i], names))
		}
	}
	t.Items[0] = strings.Join(strList, " ")
	childList := (t.ChildList)[t.CurrTab]
	if len(childList) == 0 {
		return
	}
	strList = strList[:0]
	for i, names := range childList[1] {
		str := matchKey([]byte(names), ShortKeys)
		if str != names {
			if count == MatchIndex {
				strList = append(strList, fmt.Sprintf("[%s](bg-red)<%s>", childList[0][i], str))
			} else {
				strList = append(strList, fmt.Sprintf("[%s](bg-blue)<%s>", childList[0][i], str))
			}
			count++
			MatchList = append(MatchList, i+len(t.NameList[1]))
		} else {
			strList = append(strList, fmt.Sprintf("%s<%s>", childList[0][i], names))
		}
	}
	t.Items[1] = strings.Join(strList, " ")
}

type UITopicList struct {
	ScrollList
	Label        string
	Name         string
	Type         TopicType
	AllTopicInfo []TopicInfo
}

func NewTopicList() *UITopicList {
	l := &UITopicList{Label: "Title [C-p]", Type: TopicTab}
	l.ScrollList = *NewScrollList()
	l.BorderLabel = l.Label
	return l
}

func (l *UITopicList) UpdateLabel() {
	str := l.Label
	if len(l.Name) > 0 {
		str = fmt.Sprintf("%s (%s)", str, l.Name)
	}
	if len(ShortKeys) > 0 {
		str = fmt.Sprintf("%s (%s)", str, ShortKeys)
	}
	l.BorderLabel = str
}

func (l *UITopicList) MatchTopic() {
	count := 0
	MatchList = MatchList[:0]
	l.ResetBgColor()
	log.Println("+", len(l.Items), len(l.AllItems), l.Index)
	for i := 0; i < len(l.Items); i++ {
		item := l.AllItems[i+l.Index]
		match_str := []byte(item)[:10]
		str := matchKey(match_str, ShortKeys)
		if str != string(match_str) {
			if count == MatchIndex {
				l.SetBgColor(i, ui.ColorRed)
			} else {
				l.SetBgColor(i, ui.ColorBlue)
			}
			count++
			MatchList = append(MatchList, i)
			l.SetItem(i, fmt.Sprintf("%s%s", str, []byte(item)[10:]))
		} else {
			l.SetItem(i, item)
		}
	}
}

func (l *UITopicList) Fresh(cate, node string) {
	l.Index = 0
	log.Println(cate, node)
	resetMatch()
	l.Name = "..."
	l.UpdateLabel()
	ui.Render(l)
	if len(node) > 0 {
		l.Name = node
		l.Type = TopicNode
		l.AllTopicInfo = parseTopicByNode(l.Name)
	} else {
		l.Name = cate
		l.Type = TopicTab
		l.AllTopicInfo = parseTopicByTab(l.Name)
	}
	l.DrawTopic()
	l.UpdateLabel()
}

func (l *UITopicList) DrawTopic() {
	lst := make([]string, len(l.AllTopicInfo))
	for i, info := range l.AllTopicInfo {
		prefix := fmt.Sprintf("<%02d> <%s> ", i, randID())
		prefix_width := rw.StringWidth(prefix)
		title := info.Title
		title_witth := rw.StringWidth(title)
		var suffix string
		// if len(info.Time) > 0 {
		if len(info.Time) < 0 {
			if l.Type == TopicTab {
				suffix = fmt.Sprintf("[<%d>](fg-bold,fg-blue) %s %s [%s](fg-green)", info.ReplyCount, info.Node, info.Time, info.Author)
			} else {
				suffix = fmt.Sprintf("[<%d>](fg-bold,fg-blue) %s [%s](fg-green)", info.ReplyCount, info.Time, info.Author)
			}
		} else {
			if l.Type == TopicTab {
				suffix = fmt.Sprintf("[<%d>](fg-bold,fg-blue) %s [%s](fg-green)", info.ReplyCount, info.Node, info.Author)
			} else {
				suffix = fmt.Sprintf("[<%d>](fg-bold,fg-blue) [%s](fg-green)", info.ReplyCount, info.Author)
			}
		}
		suffix_width := rw.StringWidth(suffix) - rw.StringWidth("[](fg-bold,fg-blue)[](fg-green)")
		space_width := l.InnerWidth() - 1 - (prefix_width + suffix_width + title_witth)
		if space_width < 0 {
			trim_width := l.InnerWidth() - 1 - prefix_width - suffix_width
			title_rune := []rune(title)
			w := 0
			ellip_widh := rw.StringWidth("…")
			for i, ch := range title_rune {
				w += rw.RuneWidth(ch)
				if w > trim_width-ellip_widh {
					if i > 0 {
						title = string(title_rune[:i]) + "…"
					} else {
						title = ""
					}
					break
				}
			}
			space_width = 0
		}
		lst[i] = fmt.Sprintf("%s%s%s%s", prefix, title, strings.Repeat(" ", space_width), suffix)
	}
	l.SetAllItems(lst, true)
}

type UIReplyList struct {
	ScrollList
	Topic TopicInfo
	Reply ReplyList
	Label string
}

func NewReplyList() *UIReplyList {
	r := &UIReplyList{ScrollList: *NewScrollList(), Label: "Reply [C-p]"}
	return r
}

func (l *UIReplyList) UpdateLabel() {
	str := l.Label
	if len(l.Reply.Url) > 0 {
		str = fmt.Sprintf("%s (%s)", str, l.Reply.Url)
	}
	/*
		if len(ShortKeys) > 0 {
			str = fmt.Sprintf("%s (%s)", str, ShortKeys)
		}
	*/
	l.BorderLabel = str
}

func resetMatch() {
	ShortKeys = ShortKeys[:0]
	MatchList = MatchList[:0]
	MatchIndex = 0
}

func switchState(st State) {
	log.Println(st, CurrState, CurrBodyState)
	resetMatch()
	if st == BodyStateReply && CurrBodyState != BodyStateReply {
		ui.Body.Rows[1] = ReplyRows
		ui.Body.Align()
		ui.Render(ui.Body)
		st = StateBody
		CurrBodyState = BodyStateReply
	} else if st == BodyStateTopic && CurrBodyState != BodyStateTopic {
		ui.Body.Rows[1] = TopicRows
		ui.Body.Align()
		ui.Render(ui.Body)
		st = StateBody
		CurrBodyState = BodyStateTopic
	}
	switch st {
	case StateDefault:
		uiTab.ResetTabList()
		uiTab.Highlight(false)
		uiTab.UpdateLabel()

		if CurrBodyState == BodyStateTopic {
			uiTopic.Highlight(false)
			uiTopic.UpdateLabel()
		} else {
			uiReply.Highlight(false)
			uiReply.UpdateLabel()
		}

		CurrState = StateDefault
	case StateTab:
		uiTab.MatchTab()
		uiTab.Highlight(true)
		uiTab.UpdateLabel()

		if CurrBodyState == BodyStateTopic {
			uiTopic.Highlight(false)
			uiTopic.UpdateLabel()
		} else {
			uiReply.Highlight(false)
			uiReply.UpdateLabel()
		}

		CurrState = StateTab
	case StateBody:
		uiTab.ResetTabList()
		uiTab.Highlight(false)
		uiTab.UpdateLabel()

		if CurrBodyState == BodyStateTopic {
			uiTopic.Highlight(true)
			uiTopic.UpdateLabel()
		} else {
			uiReply.Highlight(true)
			uiReply.UpdateLabel()
		}

		CurrState = StateBody
	}
}

func matchKey(str, key []byte) string {
	color_rc := "[c](fg-green)"
	color_bytes := []byte(color_rc)
	key_map := make(map[byte]uint16)
	for _, c := range key {
		key_map[c]++
	}
	name_map := make(map[byte]uint16)
	for _, c := range str {
		name_map[c]++
	}
	has := true
	for _, rc := range key {
		if key_map[rc] > name_map[rc] {
			has = false
			break
		}
	}
	if has && len(ShortKeys) > 0 {
		short := []byte{}
		for _, rc := range str {
			if key_map[rc] != 0 {
				color_bytes[1] = rc
				short = append(short, color_bytes...)
			} else {
				short = append(short, rc)
			}
		}
		return string(short)
	}
	return string(str)
}

func parseTopicByTab(tab string) (ret []TopicInfo) {
	url := fmt.Sprintf("https://www.v2ex.com/?tab=%s", tab)
	resp, err := Session.Get(url, &requests.RequestOptions{
		UserAgent: userAgent,
	})
	if err != nil {
		log.Println(err)
		return
	}

	defer log.Println(url, "status_code", resp.StatusCode)
	doc, err := goquery.NewDocumentFromResponse(resp.RawResponse)
	if err != nil {
		log.Println(err)
		return
	}
	uiUser.User.Name = strings.TrimSpace(doc.Find("span.bigger a").Text())
	doc.Find("a.fade").Each(func(i int, s *goquery.Selection) {
		if v, has := s.Attr("href"); has && v == "/notifications" {
			uiUser.User.Notify = s.Text()
		}
	})
	sliverStr := doc.Find("a.balance_area").Text()
	sliverLst := strings.Split(sliverStr, " ")
	setSli := false
	for _, sli := range sliverLst {
		if len(sli) > 0 {
			if !setSli {
				uiUser.User.Silver, _ = strconv.Atoi(sli)
				setSli = true
			} else {
				uiUser.User.Bronze, _ = strconv.Atoi(sli)
				break
			}
		}
	}
	uiUser.Fresh()
	log.Println("UserInfo", uiUser.User)
	childList := [][]string{{}, {}}
	doc.Find("div.box div.cell").Each(func(i int, s *goquery.Selection) {
		if s.Next().HasClass("cell item") && s.Prev().HasClass("inner") {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				hrefSplit := strings.Split(href, "/")
				if hrefSplit[len(hrefSplit)-2] == "go" {
					href = hrefSplit[len(hrefSplit)-1]
					childList[1] = append(childList[1], href)
					childList[0] = append(childList[0], s.Text())
				}
			})
		}
	})
	uiTab.ChildList[uiTab.CurrTab] = childList
	doc.Find("div.cell.item").Each(func(i int, s *goquery.Selection) {
		topic := TopicInfo{}
		title := s.Find(".item_title a")
		topic.Title = title.Text()
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

func parseTopicByNode(node string) (ret []TopicInfo) {
	url := fmt.Sprintf("https://www.v2ex.com/go/%s", node)
	resp, err := Session.Get(url, &requests.RequestOptions{
		UserAgent: userAgent,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer log.Println(url, "status_code", resp.StatusCode)
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

func randID() []byte {
	ret := make([]byte, 2)
	ret[0] = byte(rand.Int()%26) + 'a'
	ret[1] = byte(rand.Int()%26) + 'a'
	return ret
}

func handleKey(e ui.Event) {
	switch CurrState {
	case StateDefault:
		log.Println(e.Data.(ui.EvtKbd).KeyStr, "default")
	case StateTab, StateBody:
		key := e.Data.(ui.EvtKbd).KeyStr
		if len(key) == 1 && ((key[0] >= '0' && key[0] <= '9') || (key[0] >= 'a' && key[0] <= 'z') || (key[0] >= 'A' && key[0] <= 'Z')) {
			MatchIndex = 0
			log.Println(e.Data.(ui.EvtKbd).KeyStr, "select")
			ShortKeys = append(ShortKeys, key[0])
			if CurrState == StateTab {
				uiTab.MatchTab()
				uiTab.UpdateLabel()
			} else if CurrState == StateBody {
				uiTopic.MatchTopic()
				uiTopic.UpdateLabel()
			}
		}
		if key == "<escape>" || key == "C-c" || key == "C-u" {
			MatchIndex = 0
			ShortKeys = ShortKeys[:0]
			if CurrState == StateTab {
				uiTab.MatchTab()
				uiTab.UpdateLabel()
			} else if CurrState == StateBody {
				uiTopic.MatchTopic()
				uiTopic.UpdateLabel()
			}
		}
		if key == "C-n" && len(MatchList) > 0 {
			MatchIndex++
			MatchIndex = MatchIndex % len(MatchList)
			if CurrState == StateTab {
				uiTab.MatchTab()
			} else {
				uiTopic.MatchTopic()
			}
		}
		if key == "<enter>" {
			if CurrState == StateTab {
				ui.Body.Rows[1] = TopicRows
				ui.Body.Align()
				ui.Render(ui.Body)

				uiTab.CurrChildTab = -1
				if len(MatchList) > 0 {
					tab := MatchList[MatchIndex]
					sz := len(uiTab.NameList[0])
					if tab >= sz {
						uiTab.CurrChildTab = tab - sz
					} else {
						uiTab.CurrTab = tab
					}
				} else {
					uiTab.CurrTab = 0
				}
				uiTopic.Fresh(uiTab.GetTabNode())
				switchState(BodyStateTopic)
			} else if CurrState == StateBody {
				if len(MatchList) == 0 {
					return
				}
				if CurrBodyState == BodyStateTopic {
					idx := MatchList[MatchIndex]
					idx += uiTopic.Index

					resetMatch()

					uiReply.Topic = uiTopic.AllTopicInfo[idx]
					parseReply(uiReply.Topic.Url, &uiReply.Reply)
					items := []string{}
					text := wordwrap.WrapString(uiReply.Topic.Title, uint(uiReply.InnerWidth()))
					items = append(items, strings.Split(text, "\n")...)
					items = append(items, strings.Repeat("=", uiReply.InnerWidth()-1))
					for i, content := range uiReply.Reply.Content {
						text := wordwrap.WrapString(content, uint(uiReply.InnerWidth()))
						items = append(items, strings.Split(text, "\n")...)
						if i != len(uiReply.Reply.Content)-1 {
							items = append(items, strings.Repeat("-", uiReply.InnerWidth()-1))
						} else {
							items = append(items, strings.Repeat("=", uiReply.InnerWidth()-1))
						}
					}
					for _, rep := range uiReply.Reply.List {
						floor := fmt.Sprintf("<%d>", rep.Floor)
						replyStr := fmt.Sprintf("%s %s", floor, rep.Reply)
						text := wordwrap.WrapString(replyStr, uint(uiReply.InnerWidth()))
						text = fmt.Sprintf("[%s](fg-green)%s", floor, text[len(floor):])
						items = append(items, strings.Split(text, "\n")...)
						items = append(items, "\n")
					}
					log.Println(strconv.Itoa(uiReply.Height))
					uiReply.SetAllItems(items, true)

					switchState(BodyStateReply)
				}
			}
		}
		ui.Render(ui.Body)
	}
}

func parseReply(url string, reply *ReplyList) error {
	*reply = ReplyList{}
	resp, err := Session.Get(url, &requests.RequestOptions{
		UserAgent: userAgent,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("resp.StatusCode=%d", resp.StatusCode))
	}
	defer log.Println(url, "status_code", resp.StatusCode)
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return err
	}
	reply.Url = url
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
									pList = append(pList, text)
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
			reply.Content = append(reply.Content, strings.Join(contentList, "\n\n"))
		} else {
			log.Println("+++")
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
			reply.Content = append(reply.Content, strings.Join(contentList, ""))
		}
	})

	head := doc.Find("small.gray").Text()
	head = strings.Replace(head, string([]rune{0xA0}), "", -1)
	head = strings.Replace(head, " ", "", -1)
	headList := strings.Split(head, "·")
	reply.Lz = headList[0]
	reply.PostTime = headList[1]
	reply.ClickNum = headList[2]

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
			info.Floor, _ = strconv.Atoi(sel.Find("span.no").Text())
			info.Member = sel.Find("a.dark").Text()
			info.Time = sel.Find("span.fade.small").Text()
			info.Up = sel.Find("span.small.fade").Text()
			reply.List = append(reply.List, info)
			// log.Println(info.Floor, info.Member, info.Time)
		}
	})

	return nil
}

func login(username, password string) error {

	/*
		ui.Close()
		reply := &ReplyList{}
		// parseReply("https://www.v2ex.com/t/296412", reply)
		parseReply("https://www.v2ex.com/t/296493", reply)
		for _, content := range reply.Content {
			log.Println(content)
		}
		log.Println(reply.Lz, reply.PostTime, reply.ClickNum)

		os.Exit(1)
		return nil
	*/

	resp, err := Session.Get("https://www.v2ex.com/signin", &requests.RequestOptions{
		UserAgent: userAgent,
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
		UserAgent: userAgent,
		Headers:   Headers,
	})

	return err
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	if err := ui.Init(); err != nil {
		log.Panic(err)
	}

	uiTab = NewTab()
	uiLog = NewLog()
	uiUser = NewUser()
	uiTopic = NewTopicList()
	uiReply = NewReplyList()

	Session = requests.NewSession(nil)
}

func main() {

	defer ui.Close()

	uiTopic.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
	uiReply.Height = uiTopic.Height

	TopicRows = ui.NewCol(12, 0, uiTopic)
	ReplyRows = ui.NewCol(12, 0, uiReply)

	ui.Body.AddRows(
		ui.NewCol(12, 0, uiTab),
		TopicRows,
		ui.NewRow(
			ui.NewCol(9, 0, uiLog),
			ui.NewCol(3, 0, uiUser)))
	ui.Body.Align()
	ui.Render(ui.Body)

	log.SetOutput(uiLog)

	user := os.Getenv("V2EX_NAME")
	pass := os.Getenv("V2EX_PASS")
	if len(user) > 0 && len(pass) > 0 {
		/*
			if err := login(user, pass); err != nil {
				log.Println(err)
			}
		*/
	} else {
		log.Println("$V2EX_NAME or $V2EX_PASS is empty")
	}

	switchState(StateTab)
	// 这里来回切换状态是为了在开始时候即对UIRow进行初始化(align, inner...)
	switchState(BodyStateReply)
	switchState(BodyStateTopic)

	ui.Handle("/sys/kbd/C-q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-w", func(ui.Event) {
		if LastCtrlW == 0 {
			LastCtrlW = time.Now().Unix()
		} else {
			now := time.Now().Unix()
			if now-LastCtrlW <= 2 {
				state := (CurrState + 1) % StateMax
				if state == StateDefault {
					state++
				}
				switchState(state)
				LastCtrlW = 0
			} else {
				LastCtrlW = now
			}
		}
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/C-t", func(ui.Event) {
		if CurrState != StateTab {
			switchState(StateTab)
		} else {
			switchState(StateDefault)
		}
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/C-r", func(ui.Event) {
		if CurrState == StateBody {
			if CurrBodyState == BodyStateTopic {
				uiTopic.Fresh(uiTab.GetTabNode())
			}
		}
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			uiReply.PageDown()
		} else {
			uiTopic.PageDown()
		}
	})
	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			uiReply.PageUp()
		} else {
			uiTopic.PageUp()
		}
	})
	ui.Handle("/sys/kbd/C-e", func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			uiReply.ScrollDown()
		} else {
			uiTopic.ScrollDown()
		}
	})
	ui.Handle("/sys/kbd/C-y", func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			uiReply.ScrollUp()
		} else {
			uiTopic.ScrollUp()
		}
	})
	ui.Handle("/sys/kbd/C-p", func(ui.Event) {
		if CurrBodyState != BodyStateReply {
			switchState(BodyStateReply)
		} else {
			switchState(BodyStateTopic)
		}
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd", handleKey)
	ui.Handle("/sys/kbd/C-l", func(e ui.Event) {
		if uiTopic.Height == ui.TermHeight()-uiTab.Height {
			uiTopic.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
		} else {
			uiTopic.Height = ui.TermHeight() - uiTab.Height
		}
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		uiTopic.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
		uiReply.Height = uiTopic.Height
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	firstLoad := true
	ui.Handle("/timer/1s", func(e ui.Event) {
		if firstLoad {
			firstLoad = false
			uiTopic.Fresh(uiTab.GetTabNode())
			switchState(BodyStateTopic)
			ui.Render(ui.Body)
		}
	})
	ui.Loop()
}
