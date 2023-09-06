package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dlclark/regexp2"
)

type Data struct {
	Recall   []string
	Wiki     []string
	Note     []string
	Sentence []string
}

var SPLIT_NEWLINE = "\r\n"
var JOIN_NEWLINE = "\n"

func Split(text string, re *regexp2.Regexp) []string {
	ret := []string{}
	start := 0
	runes := []rune(text)
	m, _ := re.FindRunesMatch(runes)
	for m != nil {
		ret = append(ret, string(runes[start:m.Index]))
		start = m.Index + len(m.Runes())
		m, _ = re.FindNextMatch(m)
	}
	ret = append(ret, string(runes[start:]))
	return ret
}

func SplitSentence(text string) []string {
	re := regexp2.MustCompile(`\.\s+(?=[A-Z])`, 0)

	ret := []string{}
	for _, line := range strings.Split(text, "\n") {
		_r := Split(line, re)
		for _, s := range _r {
			if len([]rune(s)) >= 50 {
				ret = append(ret, s)
			}
		}
	}
	return ret
}

func LoadRecall(words []string) {
	for _, w := range words {
		slog.Info("Load recall " + w)
		_, err := GetZHWord(w)
		Must(err)
		InsertRecallAlways(w)
	}
}

func LoadWiki(wiki []string) {
	for _, w := range wiki {
		slog.Info("Load wiki " + w)
		title := path.Base(w)
		_, err := GetDoc(title, "wiki")
		if err == nil {
			slog.Info(fmt.Sprintf("wiki exists, ignore: %s", title))
			continue
		}
		if strings.TrimSpace(title) == "" {
			panic(fmt.Errorf("wiki title is empty: %s", w))
		}
		body := ExtractWikipedia(w)
		for _, s := range SplitSentence(body) {
			InsertDocs(s, title, "wiki")
		}

	}
}

func LoadNote(note []string) {
	type tnote struct {
		title string
		body  string
	}
	tns := []tnote{}
	for _, c := range note {
		slog.Info("Load note " + c)
		c = strings.TrimSpace(c)
		_t := strings.Split(c, SPLIT_NEWLINE)
		if len(_t) < 2 {
			panic(fmt.Errorf("note must contain more than two lines: %s", c))
		}
		if _t[0][:1] != ">" {
			panic(fmt.Errorf("note title must start with \">\":\n%s", c))
		}
		title := _t[0][1:]
		body := strings.TrimSpace(strings.Join(_t[1:], JOIN_NEWLINE))
		tns = append(tns, tnote{title: title, body: body})
	}
	for _, n := range tns {
		InsertDocs(n.body, n.title, "note")
	}
}

func LoadSentence(sentences []string) {
	for _, line := range sentences {
		slog.Info("Load sentence " + line)
		InsertDocs(line, "", "sentence")
	}
}

func LoadData(r io.Reader) {
	data := Data{}
	_, err := toml.NewDecoder(r).Decode(&data)
	Must(err)
	LoadRecall(data.Recall)
	LoadWiki(data.Wiki)
	LoadSentence(data.Sentence)
	LoadNote(data.Note)
}

func LoadDataF(file string) {
	f, err := os.Open(file)
	Must(err)
	LoadData(f)
}
