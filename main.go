package main

import (
	// "fmt"
	ui "github.com/six-ddc/termui"
	"github.com/six-ddc/v2ex-go/lib"
	"log"
	"time"
)

var (
	uiTab   *g2ex.UITab
	uiTopic *g2ex.UITopicList
	uiReply *g2ex.UIReplyList
	uiLog   *g2ex.UILog
	uiUser  *g2ex.UIUser

	TopicRows *ui.Row
	ReplyRows *ui.Row

	LastCtrlW     int64
	CurrState     g2ex.State
	CurrBodyState g2ex.State
)

func switchState(st g2ex.State) {
	log.Println("st", st, "CurrState", CurrState, "CurrBodyState", CurrBodyState)
	if st == g2ex.BodyStateReply {
		if CurrBodyState != g2ex.BodyStateReply {
			ui.Body.Rows[1] = ReplyRows
			ui.Body.Align()
			ui.Render(ui.Body)
			CurrBodyState = g2ex.BodyStateReply
		}
		st = g2ex.StateBody
	} else if st == g2ex.BodyStateTopic {
		if CurrBodyState != g2ex.BodyStateTopic {
			ui.Body.Rows[1] = TopicRows
			ui.Body.Align()
			ui.Render(ui.Body)
			CurrBodyState = g2ex.BodyStateTopic
		}
		st = g2ex.StateBody
	}
	g2ex.ResetMatch()
	switch st {
	case g2ex.StateDefault:
		uiTab.ResetTabList()
		uiTab.Highlight(false)
		uiTab.UpdateLabel()

		if CurrBodyState == g2ex.BodyStateTopic {
			uiTopic.Highlight(false)
			uiTopic.UpdateLabel()
		} else {
			uiReply.Highlight(false)
			uiReply.UpdateLabel()
		}

		CurrState = g2ex.StateDefault
	case g2ex.StateTab:
		uiTab.MatchTab()
		uiTab.Highlight(true)
		uiTab.UpdateLabel()

		if CurrBodyState == g2ex.BodyStateTopic {
			uiTopic.Highlight(false)
			uiTopic.UpdateLabel()
		} else {
			uiReply.Highlight(false)
			uiReply.UpdateLabel()
		}

		CurrState = g2ex.StateTab
	case g2ex.StateBody:
		uiTab.ResetTabList()
		uiTab.Highlight(false)
		uiTab.UpdateLabel()

		if CurrBodyState == g2ex.BodyStateTopic {
			uiTopic.Highlight(true)
			uiTopic.UpdateLabel()
		} else {
			uiReply.Highlight(true)
			uiReply.UpdateLabel()
		}

		CurrState = g2ex.StateBody
	}
}

func handleKey(e ui.Event) {
	switch CurrState {
	case g2ex.StateDefault:
		log.Println(e.Data.(ui.EvtKbd).KeyStr, "default")
	case g2ex.StateTab, g2ex.StateBody:
		key := e.Data.(ui.EvtKbd).KeyStr
		if len(key) == 1 && ((key[0] >= '0' && key[0] <= '9') || (key[0] >= 'a' && key[0] <= 'z') || (key[0] >= 'A' && key[0] <= 'Z')) {
			g2ex.MatchIndex = 0
			log.Println(e.Data.(ui.EvtKbd).KeyStr, "select")
			g2ex.ShortKeys = append(g2ex.ShortKeys, key[0])
			if CurrState == g2ex.StateTab {
				uiTab.MatchTab()
				uiTab.UpdateLabel()
			} else if CurrState == g2ex.StateBody {
				if CurrBodyState == g2ex.BodyStateTopic {
					uiTopic.MatchTopic()
					uiTopic.UpdateLabel()
				}
			}
		} else if key == "<escape>" || key == "C-8" || key == "C-c" {
			// 这里可能是bug，C-8其实是<delete>
			g2ex.MatchIndex = 0
			g2ex.ShortKeys = g2ex.ShortKeys[:0]
			if CurrState == g2ex.StateTab {
				uiTab.MatchTab()
				uiTab.UpdateLabel()
			} else if CurrState == g2ex.StateBody {
				uiTopic.MatchTopic()
				uiTopic.UpdateLabel()
			}
		} else if key == g2ex.GetConfString("key.next", "C-n") && len(g2ex.MatchList) > 0 {
			g2ex.MatchIndex++
			g2ex.MatchIndex = g2ex.MatchIndex % len(g2ex.MatchList)
			if CurrState == g2ex.StateTab {
				uiTab.MatchTab()
			} else {
				uiTopic.MatchTopic()
			}
		} else if key == g2ex.GetConfString("key.enter", "<enter>") {
			if CurrState == g2ex.StateTab {
				uiTab.CurrChildTab = -1
				if len(g2ex.MatchList) > 0 {
					tab := g2ex.MatchList[g2ex.MatchIndex]
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
				switchState(g2ex.BodyStateTopic)
			} else if CurrState == g2ex.StateBody {
				if len(g2ex.MatchList) == 0 {
					return
				}
				if CurrBodyState == g2ex.BodyStateTopic {
					idx := g2ex.MatchList[g2ex.MatchIndex]
					idx += uiTopic.Index

					uiReply.Fresh(&uiTopic.AllTopicInfo[idx], false)
					switchState(g2ex.BodyStateReply)
				}
			}
		}
	}
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	if err := ui.Init(); err != nil {
		log.Panic(err)
	}

	uiTab = g2ex.NewTab()
	uiLog = g2ex.NewLog()
	uiUser = g2ex.NewUser()
	uiTopic = g2ex.NewTopicList(uiTab, uiUser)
	uiReply = g2ex.NewReplyList()

	g2ex.SetConfFile("config.ini")
}

func uiResize() {
	ui.Body.Width = ui.TermWidth()
	if g2ex.GetConfBool("ui.enable_log", false) {
		uiTopic.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
	} else {
		uiTopic.Height = ui.TermHeight() - uiTab.Height
	}
	uiReply.Height = uiTopic.Height
}

func main() {

	defer ui.Close()

	uiResize()

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

	user := g2ex.GetConfString("user.name", "")
	pass := g2ex.GetConfString("user.pass", "")
	if len(user) > 0 && len(pass) > 0 {
		if err := g2ex.Login(user, pass); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("$V2EX_NAME or $V2EX_PASS is empty")
	}

	switchState(g2ex.StateTab)
	// 这里来回切换状态是为了在开始时候即对UIRow进行初始化(align, inner...)
	switchState(g2ex.BodyStateReply)
	switchState(g2ex.BodyStateTopic)

	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.quit", "C-q"), func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.switch", "C-w"), func(ui.Event) {
		if LastCtrlW == 0 {
			LastCtrlW = time.Now().Unix()
		} else {
			now := time.Now().Unix()
			if now-LastCtrlW <= 2 {
				state := (CurrState + 1) % g2ex.StateMax
				if state == g2ex.StateDefault {
					state++
				}
				switchState(state)
				LastCtrlW = 0
			} else {
				LastCtrlW = now
			}
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.tab", "C-t"), func(ui.Event) {
		if CurrState != g2ex.StateTab {
			switchState(g2ex.StateTab)
		} else {
			switchState(g2ex.StateDefault)
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.update", "C-r"), func(ui.Event) {
		if CurrState == g2ex.StateBody {
			if CurrBodyState == g2ex.BodyStateTopic {
				uiTopic.Fresh(uiTab.GetTabNode())
			} else {
				uiReply.Fresh(nil, false)
			}
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.pagedown", "C-f"), func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			if uiReply.PageDown() {
			}
		} else {
			if !uiTopic.PageDown() {
				uiTopic.LoadNext()
			}
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.pageup", "C-b"), func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			if !uiReply.PageUp() {
				uiReply.LoadPrev()
			}
		} else {
			uiTopic.PageUp()
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.scrolldown", "C-e"), func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			uiReply.ScrollDown()
		} else {
			if !uiTopic.ScrollDown() {
				uiTopic.LoadNext()
			}
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.scrollup", "C-y"), func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			if !uiReply.ScrollUp() {
				uiReply.LoadPrev()
			}
		} else {
			uiTopic.ScrollUp()
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.topic2reply", "C-p"), func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			switchState(g2ex.BodyStateTopic)
		} else {
			switchState(g2ex.BodyStateReply)
		}
	})
	// 这里其实是C-i
	ui.Handle("/sys/kbd/<tab>", func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateTopic {
			switchState(g2ex.BodyStateReply)
		}
	})
	ui.Handle("/sys/kbd/C-o", func(ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			switchState(g2ex.BodyStateTopic)
		}
	})
	ui.Handle("/sys/kbd/"+g2ex.GetConfString("key.log", "C-l"), func(e ui.Event) {
		if CurrBodyState == g2ex.BodyStateReply {
			if uiReply.Height == ui.TermHeight()-uiTab.Height {
				uiReply.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
			} else {
				uiReply.Height = ui.TermHeight() - uiTab.Height
			}
		} else {
			if uiTopic.Height == ui.TermHeight()-uiTab.Height {
				uiTopic.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
			} else {
				uiTopic.Height = ui.TermHeight() - uiTab.Height
			}
		}
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	ui.Handle("/sys/kbd", handleKey)
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		uiResize()
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	firstLoad := true
	ui.Handle("/timer/1s", func(e ui.Event) {
		if firstLoad {
			firstLoad = false
			uiTopic.Fresh(uiTab.GetTabNode())
			switchState(g2ex.BodyStateTopic)
		}
	})
	ui.Loop()
}
