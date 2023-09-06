package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

type Data struct {
	Recall   []string
	Wiki     []string
	Note     []string
	Sentence []string
}

var NEWLINE = "\r\n"

func LoadRecall(words []string) {
	for _, w := range words {
		_, err := GetZHWord(w)
		Must(err)
		InsertRecallAlways(w)
	}
}

func LoadWiki(wiki []string) {
	for _, w := range wiki {
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
		InsertDocs(body, title, "wiki")
	}
}

func LoadNote(note []string) {
	type tnote struct {
		title string
		body  string
	}
	tns := []tnote{}
	for _, c := range note {
		c = strings.TrimSpace(c)
		_t := strings.Split(c, NEWLINE)
		if len(_t) < 2 {
			panic(fmt.Errorf("note must contain more than two lines: %s", c))
		}
		if _t[0][:1] != ">" {
			panic(fmt.Errorf("note title must start with \">\":\n%s", c))
		}
		title := _t[0][1:]
		body := strings.TrimSpace(strings.Join(_t[1:], NEWLINE))
		tns = append(tns, tnote{title: title, body: body})
	}
	for _, n := range tns {
		InsertDocs(n.body, n.title, "note")
	}
}

func LoadSentence(sentences []string) {
	for _, line := range sentences {
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
