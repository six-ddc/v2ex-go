package main

import (
	// "fmt"
	ui "github.com/six-ddc/termui"
	"log"
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

	LastCtrlW     int64
	CurrState     State
	CurrBodyState State
)

func switchState(st State) {
	log.Println("st", st, "CurrState", CurrState, "CurrBodyState", CurrBodyState)
	if st == BodyStateReply {
		if CurrBodyState != BodyStateReply {
			ui.Body.Rows[1] = ReplyRows
			ui.Body.Align()
			ui.Render(ui.Body)
			CurrBodyState = BodyStateReply
		}
		st = StateBody
	} else if st == BodyStateTopic {
		if CurrBodyState != BodyStateTopic {
			ui.Body.Rows[1] = TopicRows
			ui.Body.Align()
			ui.Render(ui.Body)
			CurrBodyState = BodyStateTopic
		}
		st = StateBody
	}
	ResetMatch()
	switch st {
	case StateDefault:
		uiTab.ResetTabList()
		uiTab.Highlight(false)
		uiTab.UpdateLabel()

		if CurrBodyState == BodyStateTopic {
			uiTopic.Highlight(false)
			uiTopic.UpdateLabel()
		} else if CurrBodyState == BodyStateReply {
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
		} else if CurrBodyState == BodyStateReply {
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
		} else if CurrBodyState == BodyStateReply {
			uiReply.Highlight(true)
			uiReply.UpdateLabel()
		}

		CurrState = StateBody
	}
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
				if CurrBodyState == BodyStateTopic {
					uiTopic.MatchTopic()
					uiTopic.UpdateLabel()
				}
			}
		} else if key == "<escape>" || key == "C-8" || key == "C-c" {
			// 这里可能是bug，C-8其实是<delete>
			MatchIndex = 0
			ShortKeys = ShortKeys[:0]
			if CurrState == StateTab {
				uiTab.MatchTab()
				uiTab.UpdateLabel()
			} else if CurrState == StateBody {
				uiTopic.MatchTopic()
				uiTopic.UpdateLabel()
			}
		} else if key == GetConfString("key.next", "C-n") && len(MatchList) > 0 {
			MatchIndex++
			MatchIndex = MatchIndex % len(MatchList)
			if CurrState == StateTab {
				uiTab.MatchTab()
			} else {
				uiTopic.MatchTopic()
			}
		} else if key == GetConfString("key.enter", "<enter>") {
			if CurrState == StateTab {
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

					uiReply.Fresh(&uiTopic.AllTopicInfo[idx], false)
					switchState(BodyStateReply)
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

	uiTab = NewTab()
	uiLog = NewLog()
	uiUser = NewUser()
	uiTopic = NewTopicList(uiTab, uiUser)
	uiReply = NewReplyList()

	SetConfFile("config.ini")
}

func uiResize() {
	ui.Body.Width = ui.TermWidth()
	if GetConfBool("ui.enable_log", false) {
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

	user := GetConfString("user.name", "")
	pass := GetConfString("user.pass", "")
	if len(user) > 0 && len(pass) > 0 {
		if err := Login(user, pass); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("$V2EX_NAME or $V2EX_PASS is empty")
	}

	switchState(StateTab)
	// 这里来回切换状态是为了在开始时候即对UIRow进行初始化(align, inner...)
	switchState(BodyStateReply)
	switchState(BodyStateTopic)

	ui.Handle("/sys/kbd/"+GetConfString("key.quit", "C-q"), func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.switch", "C-w"), func(ui.Event) {
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
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.tab", "C-t"), func(ui.Event) {
		if CurrState != StateTab {
			switchState(StateTab)
		} else {
			switchState(StateDefault)
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.update", "C-r"), func(ui.Event) {
		if CurrState == StateBody {
			if CurrBodyState == BodyStateTopic {
				uiTopic.Fresh(uiTab.GetTabNode())
			} else if CurrBodyState == BodyStateReply {
				uiReply.Fresh(nil, false)
			}
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.pagedown", "C-f"), func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			if uiReply.PageDown() {
			}
		} else if CurrBodyState == BodyStateTopic {
			if !uiTopic.PageDown() {
				uiTopic.LoadNext()
			}
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.pageup", "C-b"), func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			if !uiReply.PageUp() {
				uiReply.LoadPrev()
			}
		} else if CurrBodyState == BodyStateTopic {
			uiTopic.PageUp()
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.scrolldown", "C-e"), func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			uiReply.ScrollDown()
		} else if CurrBodyState == BodyStateTopic {
			if !uiTopic.ScrollDown() {
				uiTopic.LoadNext()
			}
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.scrollup", "C-y"), func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			if !uiReply.ScrollUp() {
				uiReply.LoadPrev()
			}
		} else if CurrBodyState == BodyStateTopic {
			uiTopic.ScrollUp()
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.topic2reply", "C-p"), func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			switchState(BodyStateTopic)
		} else if CurrBodyState == BodyStateTopic {
			switchState(BodyStateReply)
		}
	})
	// 这里其实是C-i
	ui.Handle("/sys/kbd/<tab>", func(ui.Event) {
		if CurrBodyState == BodyStateTopic {
			switchState(BodyStateReply)
		}
	})
	ui.Handle("/sys/kbd/C-o", func(ui.Event) {
		if CurrBodyState == BodyStateReply {
			switchState(BodyStateTopic)
		}
	})
	ui.Handle("/sys/kbd/"+GetConfString("key.log", "C-l"), func(e ui.Event) {
		if CurrBodyState == BodyStateReply {
			if uiReply.Height == ui.TermHeight()-uiTab.Height {
				uiReply.Height = ui.TermHeight() - uiLog.Height - uiTab.Height
			} else {
				uiReply.Height = ui.TermHeight() - uiTab.Height
			}
		} else if CurrBodyState == BodyStateTopic {
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
			switchState(BodyStateTopic)
		}
	})
	ui.Loop()
}
