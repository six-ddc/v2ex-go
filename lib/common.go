package g2ex

import (
	"bytes"
	rw "github.com/mattn/go-runewidth"
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
	Source string
}

type ReplyList struct {
	Title    string
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

var UserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.82 Safari/537.36"

func MatchKey(str, key []byte) string {
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

func WrapString(str string, limit int) string {
	wid := 0
	buf := bytes.NewBuffer([]byte{})
	for _, ch := range str {
		w := rw.RuneWidth(ch)
		if ch == '\n' {
			wid = 0
		} else {
			if wid+w > limit {
				buf.WriteRune('\n')
				wid = 0
			} else {
				wid += w
			}
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}
