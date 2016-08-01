package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	ui "github.com/six-ddc/termui"
	// requests "github.com/levigross/grequests"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ui_list *UITopicList
	ui_log  *UILog
	ui_tab  *UITab

	LastCtrlW int64
	CurrState State

	ShortKeys  []byte
	MatchList  []int
	MatchIndex int
)

type UILog struct {
	Widget *ui.List
	Label  string
	Index  int
}

func NewLog() *UILog {
	l := &UILog{Index: 0, Label: "Log [C-l]"}
	l.Widget = ui.NewList()
	l.Widget.Height = 20
	l.Widget.BorderLabel = l.Label
	l.Widget.Items = make([]string, l.Widget.Height-2)
	return l
}

type UITab struct {
	Widget       *ui.List
	Label        string
	NameList     [][]string
	ChildList    [][][]string
	CurrTab      int
	CurrChildTab int
}

func NewTab() *UITab {
	t := &UITab{Label: "Tab [C-t]", CurrTab: 0, CurrChildTab: -1}
	t.NameList = [][]string{
		{"技术", "创意", "好玩", "Apple", "酷工作", "交易", "城市", "问与答", "最热", "全部", "R2", "节点", "关注"},
		{"tech", "creative", "play", "apple", "jobs", "deals", "city", "qna", "hot", "all", "r2", "nodes", "members"}}
	t.ChildList = [][][]string{
		{
			{"程序员", "Python", "iDev", "Android", "Linux", "node.js", "云计算", "宽带症候群"},
			{"programmer", "python", "idev", "android", "linux", "nodejs", "cloud", "bb"},
		},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
		{},
	}
	t.Widget = ui.NewList()
	t.Widget.BorderLabel = t.Label
	t.Widget.Height = 2 + 2
	t.Widget.Items = make([]string, 2)
	t.ResetTabList()
	return t
}

func (t *UITab) Highlight(b bool) {
	if b {
		t.Widget.BorderFg = ui.ColorRed
	} else {
		t.Widget.BorderFg = ui.ColorDefault
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
	t.Widget.Items[0] = strings.Join(strList, " ")
	childList := (t.ChildList)[t.CurrTab]
	if len(childList) == 0 {
		t.Widget.Items[1] = ""
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
	t.Widget.Items[1] = strings.Join(strList, " ")
}

func (t *UITab) UpdateLabel() {
	str := t.Label
	if len(ShortKeys) > 0 {
		str = fmt.Sprintf("%s (%s)", str, ShortKeys)
	}
	t.Widget.BorderLabel = str
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
	t.Widget.Items[0] = strings.Join(strList, " ")
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
	t.Widget.Items[1] = strings.Join(strList, " ")
}

type TopicType uint16

const (
	TopicTab TopicType = iota
	TopicNode
)

type UITopicList struct {
	Tab           *UITab
	Widget        *ui.List
	Label         string
	Name          string
	Type          TopicType
	AllTopicItems []string
	TopicFirst    int
	AllTopicInfo  []TopicInfo
}

func NewTopicList(t *UITab) *UITopicList {
	l := &UITopicList{Tab: t, Label: "Title [C-p]", TopicFirst: 0, Type: TopicTab}
	l.Widget = ui.NewList()
	l.Widget.BorderLabel = l.Label
	return l
}

func (l *UITopicList) Highlight(b bool) {
	if b {
		l.Widget.BorderFg = ui.ColorRed
	} else {
		l.Widget.BorderFg = ui.ColorDefault
	}
}

func (l *UITopicList) Height() int {
	return l.Widget.Height - 2
}

func (l *UITopicList) SetItem(i int, item string) {
	l.Widget.Items[i] = item
}

func (l *UITopicList) SetBgColor(i int, color ui.Attribute) {
	if len(l.Widget.ItemBg) != len(l.Widget.Items) {
		l.ResetBgColor()
	}
	l.Widget.ItemBg[i] = color
}

func (l *UITopicList) SetItems(items []string, updateLastList bool) {
	if updateLastList {
		l.AllTopicItems = items
	}
	sz := l.Height()
	if len(items) < sz {
		sz = len(items)
	}
	l.Widget.Items = make([]string, sz)
	copy(l.Widget.Items, items) // 复制长度以较小的slice为准
	l.ResetBgColor()
}

func (l *UITopicList) ResetBgColor() {
	l.Widget.ItemBg = make([]ui.Attribute, len(l.Widget.Items))
	for i, _ := range l.Widget.ItemBg {
		l.Widget.ItemBg[i] = ui.ThemeAttr("list.item.bg")
	}
}

func (l *UITopicList) UpdateLabel() {
	str := l.Label
	if len(l.Name) > 0 {
		str = fmt.Sprintf("%s (%s)", str, l.Name)
	}
	if len(ShortKeys) > 0 {
		str = fmt.Sprintf("%s (%s)", str, ShortKeys)
	}
	l.Widget.BorderLabel = str
}

func (l *UITopicList) MatchTopic() {
	count := 0
	MatchList = MatchList[:0]
	l.ResetBgColor()
	uiPrintln("+", len(l.Widget.Items), len(l.AllTopicItems), l.TopicFirst)
	for i := 0; i < len(l.Widget.Items); i++ {
		item := l.AllTopicItems[i+l.TopicFirst]
		str := matchKey([]byte(item), ShortKeys)
		if str != item {
			if count == MatchIndex {
				l.SetBgColor(i, ui.ColorRed)
			} else {
				l.SetBgColor(i, ui.ColorBlue)
			}
			count++
			MatchList = append(MatchList, i)
		}
		l.SetItem(i, str)
	}
}

func (l *UITopicList) Fresh() {
	cate, node := l.Tab.GetTabNode()
	uiPrintln(cate, node)
	resetMatch()
	l.Name = "..."
	l.UpdateLabel()
	ui.Render(ui.Body)
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
		lst[i] = fmt.Sprintf("<%2d> <%s> %s|[%s](fg-yellow)", i, randID(), info.Title, info.Author)
	}
	l.SetItems(lst, true)
}

func (l *UITopicList) ScrollDown() {
	sz := len(l.AllTopicItems)
	screen_heigth := l.Height()
	if sz > screen_heigth+l.TopicFirst {
		l.ResetBgColor()
		l.TopicFirst++
		l.SetItems(l.AllTopicItems[l.TopicFirst:], false)
		ui.Render(ui.Body)
	}
}

func (l *UITopicList) PageDown() {
	sz := len(l.AllTopicItems)
	screen_heigth := l.Height()
	if sz < screen_heigth {
		return
	}
	index := l.TopicFirst + screen_heigth
	if index > sz-screen_heigth {
		index = sz - screen_heigth
		if index == l.TopicFirst {
			return
		}
	}
	l.TopicFirst = index
	l.SetItems(l.AllTopicItems[l.TopicFirst:], false)
	ui.Render(ui.Body)
}

func (l *UITopicList) PageUp() {
	screen_heigth := l.Height()
	index := l.TopicFirst - screen_heigth
	if index < 0 {
		index = 0
		if index == l.TopicFirst {
			return
		}
	}
	l.TopicFirst = index
	l.SetItems(l.AllTopicItems[l.TopicFirst:], false)
	ui.Render(ui.Body)
}

func (l *UITopicList) ScrollUp() {
	if l.TopicFirst > 0 {
		l.ResetBgColor()
		l.TopicFirst--
		l.SetItems(l.AllTopicItems[l.TopicFirst:], false)
		ui.Render(ui.Body)
	}
}

type TopicInfo struct {
	Title      string
	Author     string
	AuthorImg  string
	Node       string
	Time       string
	LastReply  string
	ReplyCount int
}

type State int

const (
	StateDefault State = iota
	StateTab
	StateTopic
	StateMax
)

func uiLog(str string) {
	if ui_log == nil {
		log.Println(str)
		return
	}
	str = fmt.Sprintf("[%d] %s", ui_log.Index+1, str)
	if ui_log.Widget.Items[len(ui_log.Widget.Items)-1] != "" {
		i := 0
		for ; i < len(ui_log.Widget.Items)-1; i++ {
			ui_log.Widget.Items[i] = ui_log.Widget.Items[i+1]
		}
		ui_log.Widget.Items[i] = str
	} else {
		for i, item := range ui_log.Widget.Items {
			if item == "" {
				ui_log.Widget.Items[i] = str
				break
			}
		}
	}
	ui_log.Index++
	ui.Render(ui.Body)
}

// Write(p []byte) (n int, err error)
func uiPrintln(a ...interface{}) {
	uiLog(fmt.Sprint(a))
}

func init() {
	LastCtrlW = 0
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	if err := ui.Init(); err != nil {
		log.Panic(err)
	}
	log.SetOutput(os.Stdout)
}

func resetMatch() {
	ShortKeys = ShortKeys[:0]
	MatchList = MatchList[:0]
	MatchIndex = 0
}

func switchState(st State) {
	resetMatch()
	switch st {
	case StateDefault:
		ui_tab.ResetTabList()
		ui_tab.Highlight(false)
		ui_tab.UpdateLabel()

		ui_list.Highlight(false)
		ui_list.UpdateLabel()

		CurrState = StateDefault
	case StateTab:
		ui_tab.MatchTab()
		ui_tab.Highlight(true)
		ui_tab.UpdateLabel()

		ui_list.Highlight(false)
		ui_list.UpdateLabel()

		CurrState = StateTab
	case StateTopic:
		ui_tab.ResetTabList()
		ui_tab.Highlight(false)
		ui_tab.UpdateLabel()

		ui_list.Highlight(true)
		ui_list.UpdateLabel()

		CurrState = StateTopic
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
	uiPrintln(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		uiPrintln(err)
		return
	}
	doc.Find("div.cell.item").Each(func(i int, s *goquery.Selection) {
		topic := TopicInfo{}
		topic.Title = s.Find(".item_title a").Text()
		info := s.Find(".small.fade").Text()
		infoList := strings.Split(info, "•")
		topic.Node = strings.TrimSpace(infoList[0])
		topic.Author = strings.TrimSpace(infoList[1])
		if len(infoList) > 2 {
			topic.Time = strings.TrimSpace(infoList[2])
			if len(infoList) > 3 {
				topic.LastReply = strings.TrimSpace(infoList[3])
			}
		}
		replyCount := s.Find("a count_livid").Text()
		if replyCount != "" {
			topic.ReplyCount, _ = strconv.Atoi(replyCount)
		}
		ret = append(ret, topic)
		// log.Println(s.Find(".small.fade a.node").Text())
		// log.Println(s.Find(".small.fade strong a").Text())
	})
	return
}

func parseTopicByNode(node string) (ret []TopicInfo) {
	url := fmt.Sprintf("https://www.v2ex.com/go/%s", node)
	uiPrintln(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		uiPrintln(err)
		return
	}
	doc.Find("div.cell").Each(func(i int, s *goquery.Selection) {
		info := s.Find(".small.fade").Text()
		info = strings.TrimSpace(info)
		if len(info) == 0 {
			return
		}
		topic := TopicInfo{}
		topic.Title = s.Find(".item_title a").Text()
		infoList := strings.Split(info, "•")
		topic.Node = node
		topic.Author = strings.TrimSpace(infoList[0])
		if len(infoList) > 2 {
			topic.Time = strings.TrimSpace(infoList[1])
			if len(infoList) > 3 {
				topic.LastReply = strings.TrimSpace(infoList[2])
			}
		}
		replyCount := s.Find("a count_livid").Text()
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
		uiPrintln(e.Data.(ui.EvtKbd).KeyStr, "default")
	case StateTab, StateTopic:
		key := e.Data.(ui.EvtKbd).KeyStr
		if len(key) == 1 && ((key[0] >= '0' && key[0] <= '9') || (key[0] >= 'a' && key[0] <= 'z') || (key[0] >= 'A' && key[0] <= 'Z')) {
			MatchIndex = 0
			uiPrintln(e.Data.(ui.EvtKbd).KeyStr, "select")
			ShortKeys = append(ShortKeys, key[0])
			if CurrState == StateTab {
				ui_tab.MatchTab()
				ui_tab.UpdateLabel()
			} else if CurrState == StateTopic {
				ui_list.MatchTopic()
				ui_list.UpdateLabel()
			}
		}
		if key == "<escape>" || key == "C-c" || key == "C-u" {
			MatchIndex = 0
			ShortKeys = ShortKeys[:0]
			if CurrState == StateTab {
				ui_tab.MatchTab()
				ui_tab.UpdateLabel()
			} else if CurrState == StateTopic {
				ui_list.MatchTopic()
				ui_list.UpdateLabel()
			}
		}
		if key == "C-n" && len(MatchList) > 0 {
			MatchIndex++
			MatchIndex = MatchIndex % len(MatchList)
			if CurrState == StateTab {
				ui_tab.MatchTab()
			} else {
				ui_list.MatchTopic()
			}
		}
		if key == "<enter>" {
			if CurrState == StateTab {
				ui_tab.CurrChildTab = -1
				if len(MatchList) > 0 {
					tab := MatchList[MatchIndex]
					sz := len(ui_tab.NameList[0])
					if tab >= sz {
						ui_tab.CurrChildTab = tab - sz
					} else {
						ui_tab.CurrTab = tab
					}
				} else {
					ui_tab.CurrTab = 0
				}
				uiPrintln("---", MatchList, MatchIndex, ui_tab.CurrChildTab, ui_tab.CurrTab)
				ui_list.Fresh()
				switchState(StateTopic)
			}
		}
	}
	ui.Render(ui.Body)
}

func main() {
	/*
		parseTopicByNode("programmer")
		os.Exit(1)
	*/

	defer ui.Close()

	ui_tab = NewTab()

	ui_log = NewLog()

	ui_list = NewTopicList(ui_tab)
	ui_list.Widget.Height = ui.TermHeight() - ui_log.Widget.Height - ui_tab.Widget.Height

	ui.Body.AddRows(
		ui.NewCol(12, 0, ui_tab.Widget),
		ui.NewCol(12, 0, ui_list.Widget),
		ui.NewCol(12, 0, ui_log.Widget))
	ui.Body.Align()
	ui.Render(ui.Body)

	switchState(StateTab)

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
		if CurrState == StateTopic {
			ui_list.Fresh()
		}
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		ui_list.PageDown()
	})
	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		ui_list.PageUp()
	})
	ui.Handle("/sys/kbd/C-e", func(ui.Event) {
		ui_list.ScrollDown()
	})
	ui.Handle("/sys/kbd/C-y", func(ui.Event) {
		ui_list.ScrollUp()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd/C-p", func(ui.Event) {
		if CurrState != StateTopic {
			switchState(StateTopic)
		} else {
			switchState(StateDefault)
		}
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd", handleKey)
	ui.Handle("/sys/kbd/C-l", func(e ui.Event) {
		if ui_list.Widget.Height == ui.TermHeight()-ui_tab.Widget.Height {
			ui_list.Widget.Height = ui.TermHeight() - ui_log.Widget.Height - ui_tab.Widget.Height
		} else {
			ui_list.Widget.Height = ui.TermHeight() - ui_tab.Widget.Height
		}
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui_list.Widget.Height = ui.TermHeight() - ui_log.Widget.Height - ui_tab.Widget.Height
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	firstLoad := true
	ui.Handle("/timer/1s", func(e ui.Event) {
		if firstLoad {
			firstLoad = false
			ui_list.Fresh()
			switchState(StateTopic)
			ui.Render(ui.Body)
		}
	})
	ui.Loop()
}
