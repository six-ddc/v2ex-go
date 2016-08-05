package g2ex

import (
	"fmt"
	rw "github.com/mattn/go-runewidth"
	ui "github.com/six-ddc/termui"
	"log"
	"math/rand"
	"strings"
)

var (
	ShortKeys  []byte
	MatchList  []int
	MatchIndex int
)

func ResetMatch() {
	ShortKeys = ShortKeys[:0]
	MatchList = MatchList[:0]
	MatchIndex = 0
}

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

func (l *ScrollList) SetAllItems(items []string) {
	l.AllItems = items
	sz := l.InnerHeight()
	log.Println("sz", sz, "len(items)", len(items), "index", l.Index)
	if len(items)-l.Index < sz {
		sz = len(items) - l.Index
	}
	l.Items = make([]string, sz)
	// 显示index后面的部分
	copy(l.Items, items[l.Index:]) // 复制长度以较小的slice为准
	l.ResetBgColor()
	ui.Render(l)
}

func (l *ScrollList) SetBgColor(i int, color ui.Attribute) {
	if len(l.ItemBg) != len(l.Items) {
		l.ResetBgColor()
	}
	l.ItemBg[i] = color
	ui.Render(l)
}

func (l *ScrollList) Highlight(b bool) {
	if b {
		l.BorderFg = ui.ColorRed
	} else {
		l.BorderFg = ui.ColorDefault
	}
	ui.Render(l)
}

func (l *ScrollList) ResetBgColor() {
	l.ItemBg = make([]ui.Attribute, len(l.Items))
	for i, _ := range l.ItemBg {
		l.ItemBg[i] = ui.ThemeAttr("list.item.bg")
	}
	ui.Render(l)
}

func (l *ScrollList) ScrollDown() bool {
	sz := len(l.AllItems)
	screen_heigth := l.InnerHeight()
	if sz > screen_heigth+l.Index {
		l.ResetBgColor()
		l.Index++
		l.SetAllItems(l.AllItems)
		ui.Render(l)
		return true
	}
	return false
}

func (l *ScrollList) PageDown() bool {
	sz := len(l.AllItems)
	screen_heigth := l.InnerHeight()
	if sz < screen_heigth {
		return false
	}
	index := l.Index + screen_heigth
	if index >= sz {
		return false
	}
	l.Index = index
	l.SetAllItems(l.AllItems)
	ui.Render(l)
	return true
}

func (l *ScrollList) PageUp() bool {
	if l.Index == 0 {
		return false
	}
	screen_heigth := l.InnerHeight()
	index := l.Index - screen_heigth
	if index < 0 {
		index = 0
	}
	l.Index = index
	l.SetAllItems(l.AllItems)
	ui.Render(l)
	return true
}

func (l *ScrollList) ScrollUp() bool {
	if l.Index > 0 {
		l.ResetBgColor()
		l.Index--
		l.SetAllItems(l.AllItems)
		ui.Render(l)
		return true
	}
	return false
}

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
	str := fmt.Sprintf("[%d] %s", l.Index+1, p)
	if l.Items[len(l.Items)-1] != "" {
		i := 0
		for ; i < len(l.Items)-1; i++ {
			l.Items[i] = l.Items[i+1]
		}
		l.Items[i] = str
	} else {
		for i, item := range l.Items {
			if item == "" {
				l.Items[i] = str
				break
			}
		}
	}
	l.Index++
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
	ui.Render(t)
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
	ui.Render(t)
}

func (t *UITab) UpdateLabel() {
	str := t.Label
	if len(ShortKeys) > 0 {
		str = fmt.Sprintf("%s (%s)", str, ShortKeys)
	}
	t.BorderLabel = str
	ui.Render(t)
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
		str := MatchKey([]byte(names), ShortKeys)
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
		str := MatchKey([]byte(names), ShortKeys)
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
	ui.Render(t)
}

type UITopicList struct {
	ScrollList
	uiTab        *UITab
	uiUser       *UIUser
	Label        string
	Name         string
	Type         TopicType
	AllTopicInfo []TopicInfo
	Page         int
}

func NewTopicList(tab *UITab, user *UIUser) *UITopicList {
	l := &UITopicList{Label: "Topic [C-p]", Type: TopicTab, Page: 1, uiTab: tab, uiUser: user}
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
	ui.Render(l)
}

func (l *UITopicList) MatchTopic() {
	count := 0
	MatchList = MatchList[:0]
	l.ResetBgColor()
	log.Println("+", len(l.Items), len(l.AllItems), l.Index)
	for i := 0; i < len(l.Items); i++ {
		item := l.AllItems[i+l.Index]
		match_str := []byte(item)[:10]
		str := MatchKey(match_str, ShortKeys)
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
	ui.Render(l)
}

func (l *UITopicList) Fresh(cate, node string) {
	l.Index = 0
	log.Println(cate, node)
	ResetMatch()
	l.Name = "..."
	l.UpdateLabel()
	ui.Render(l)
	if len(node) > 0 {
		l.Name = node
		l.Type = TopicNode
		l.AllTopicInfo = ParseTopicByNode(l.Name, 1)
	} else {
		l.Name = cate
		l.Type = TopicTab
		tabList := [][]string{}
		l.AllTopicInfo = ParseTopicByTab(l.Name, l.uiUser.User, tabList)
		l.uiTab.ChildList[l.uiTab.CurrTab] = tabList
	}
	l.DrawTopic()
	l.UpdateLabel()
	ui.Render(l)
}

func (l *UITopicList) LoadNext() {
	if l.Type == TopicTab {
		// tab 不支持翻页
		return
	}
	l.Page += 1
	name := l.Name
	log.Println(l.Name, l.Page)
	ResetMatch()
	l.Name = "..."
	l.UpdateLabel()
	ui.Render(l)
	l.Name = name
	tpList := ParseTopicByNode(l.Name, l.Page)
	l.AllTopicInfo = append(l.AllTopicInfo, tpList...)
	l.DrawTopic()
	l.UpdateLabel()
	ui.Render(l)
}

func randID() []byte {
	ret := make([]byte, 2)
	ret[0] = byte(rand.Int()%26) + 'a'
	ret[1] = byte(rand.Int()%26) + 'a'
	return ret
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
	l.SetAllItems(lst)
}

type UIReplyList struct {
	ScrollList
	Topic *TopicInfo
	Reply ReplyList
	Label string
}

func NewReplyList() *UIReplyList {
	r := &UIReplyList{ScrollList: *NewScrollList(), Label: "Reply [C-p]"}
	return r
}

func (l *UIReplyList) UpdateLabel() {
	str := l.Label
	if l.Topic != nil && len(l.Topic.Url) > 0 {
		str = fmt.Sprintf("%s [%s](fg-cyan) (%s)", str, l.Reply.Lz, l.Topic.Url)
	}
	/*
		if len(ShortKeys) > 0 {
			str = fmt.Sprintf("%s (%s)", str, ShortKeys)
		}
	*/
	l.BorderLabel = str
	ui.Render(l)
}

func (l *UIReplyList) Fresh(topic *TopicInfo, addToHead bool) {
	l.Index = 0
	if topic != nil {
		l.Topic = topic
	}
	if l.Topic == nil {
		return
	}
	if addToHead { // 加载上一页的情况
		var reply ReplyList
		if ParseReply(l.Topic.Url, &reply) != nil {
			return
		}
		l.Reply.List = append(reply.List, l.Reply.List...)
	} else {
		if ParseReply(l.Topic.Url, &l.Reply) != nil {
			return
		}
	}
	log.Println("addToHead", addToHead)
	items := []string{"\n"}
	text := WrapString(l.Topic.Title, l.InnerWidth()-1)
	items = append(items, strings.Split(text, "\n")...)
	items = append(items, "\n")
	if len(l.Reply.Content) > 0 {
		items = append(items, strings.Repeat("=", l.InnerWidth()-1))
	}
	for i, content := range l.Reply.Content {
		text := WrapString(content, l.InnerWidth()-1)
		items = append(items, strings.Split(text, "\n")...)
		if i != len(l.Reply.Content)-1 {
			items = append(items, strings.Repeat("-", l.InnerWidth()-1))
		}
	}
	if len(l.Reply.Content) > 0 {
		items = append(items, strings.Repeat("=", l.InnerWidth()-1))
		items = append(items, "\n")
	}
	for i, rep := range l.Reply.List {
		source := strings.Replace(rep.Source, "♥", " [♥](fg-red)", 1)
		if rep.Member == l.Reply.Lz {
			items = append(items, fmt.Sprintf("[%d](fg-blue) [%s](fg-cyan) %s", rep.Floor, rep.Member, source))
		} else {
			items = append(items, fmt.Sprintf("[%d](fg-blue) [%s](fg-green) %s", rep.Floor, rep.Member, source))
		}
		text := WrapString(rep.Reply, l.InnerWidth()-1)
		text = strings.Replace(text, "@"+l.Reply.Lz, fmt.Sprintf("[@%s](fg-cyan)", l.Reply.Lz), 1)
		items = append(items, strings.Split(text, "\n")...)
		if i != len(l.Reply.List)-1 {
			items = append(items, "\n")
		}
	}
	l.SetAllItems(items)
	ui.Render(l)
}

func (l *UIReplyList) LoadPrev() {
	if len(l.Reply.List) == 0 || l.Reply.List[0].Floor/100 == 0 {
		return
	}
	page := l.Reply.List[0].Floor / 100
	idx := strings.Index(l.Topic.Url, "#reply")
	if idx > -1 {
		l.Topic.Url = l.Topic.Url[:idx]
	}
	url := fmt.Sprintf("%s?p=%d", l.Topic.Url, page)
	l.Topic.Url = "..."
	l.UpdateLabel()
	l.Topic.Url = url
	l.Fresh(l.Topic, true)
	l.UpdateLabel()
}
