package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/dlclark/regexp2"
)

type Cmd struct {
	File string `arg:"-f,env:EN_DRAFT" default:"draft.toml"`
}

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

func main() {
	var cmd Cmd
	arg.MustParse(&cmd)

	//os.Setenv("HTTPS_PROXY", "http://127.0.0.1:1080")
	//fmt.Println(ExtractWikipedia("https://en.wikipedia.org/wiki/MG_42"))
	/*
		InitDB()
		LoadDataF(cmd.File)
	*/
	InitDB()
	f, err := os.Open("wiki.txt")
	Must(err)
	_t, err := io.ReadAll(f)
	Must(err)
	text := string(_t)
	for _, s := range SplitSentence(text) {
		fmt.Println("insert", s)
		InsertDocs(s, "MG_42", "wiki")
	}
}
