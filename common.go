package main

import (
	"bytes"
	"github.com/go-ini/ini"
	rw "github.com/mattn/go-runewidth"
	"strings"
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

var UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.82 Safari/537.36"

func MatchKey(str, key []byte) string {
	colorRc := "[c](fg-green)"
	colorBytes := []byte(colorRc)
	keyMap := make(map[byte]uint16)
	for _, c := range key {
		keyMap[c]++
	}
	nameMap := make(map[byte]uint16)
	for _, c := range str {
		nameMap[c]++
	}
	has := true
	for _, rc := range key {
		if keyMap[rc] > nameMap[rc] {
			has = false
			break
		}
	}
	if has && len(ShortKeys) > 0 {
		short := []byte{}
		for _, rc := range str {
			if keyMap[rc] != 0 {
				colorBytes[1] = rc
				short = append(short, colorBytes...)
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
			// ui.Lise的最后一个rune会被显示...
			if wid+w >= limit-3 {
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

var iniCfg *ini.File

func SetConfFile(f string) (err error) {
	iniCfg, err = ini.Load(f)
	return err
}

func GetConfString(secKey string, defau string) string {
	if iniCfg == nil {
		return defau
	}
	sk := strings.Split(secKey, ".")
	s, k := sk[0], sk[1]
	sec, err := iniCfg.GetSection(s)
	if err != nil {
		return defau
	}
	key, err := sec.GetKey(k)
	if err != nil {
		return defau
	}
	return key.String()
}

func GetConfInt(secKey string, defau int) int {
	if iniCfg == nil {
		return defau
	}
	sk := strings.Split(secKey, ".")
	s, k := sk[0], sk[1]
	sec, err := iniCfg.GetSection(s)
	if err != nil {
		return defau
	}
	key, err := sec.GetKey(k)
	if err != nil {
		return defau
	}
	i, _ := key.Int()
	return i
}

func GetConfBool(secKey string, defau bool) bool {
	if iniCfg == nil {
		return defau
	}
	sk := strings.Split(secKey, ".")
	s, k := sk[0], sk[1]
	sec, err := iniCfg.GetSection(s)
	if err != nil {
		return defau
	}
	key, err := sec.GetKey(k)
	if err != nil {
		return defau
	}
	b, _ := key.Bool()
	return b
}
