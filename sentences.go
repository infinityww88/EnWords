package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type WordInfo struct {
	word   string
	lineno int
}

type SentenceInfo struct {
	Sentence string
	Words    []WordInfo
}

func LoadWord(word string) (Word, error) {
	DownloadWordSound(word)
	w, err := GetOnlineWordApi(word)
	if err != nil {
		return Word{}, err
	}
	w = InsertWordAlways(w)
	return w, nil
}

func LoadSentenceInfo(sis []SentenceInfo) {
	for _, si := range sis {
		sid := InsertSentenceAlways(si.Sentence)
		for _, wi := range si.Words {
			var w Word
			var err error
			var ok bool
			if w, ok = GetWord(wi.word); !ok {
				if w, err = LoadWord(wi.word); err != nil {
					panic(fmt.Errorf("load word %s error %w at line [%d]", wi.word, err, wi.lineno))
				}
			}
			InsertWordSentenceAlways(w.Wid, sid)
		}
	}
}

func ReadSentenceInfo(reader io.Reader) []SentenceInfo {
	const (
		STATE_SENTENCE = "sentence"
		STATE_WORDS    = "words"
	)
	state := STATE_SENTENCE
	ret := []SentenceInfo{}
	var si SentenceInfo

	p := bufio.NewScanner(reader)
	lineno := 0
	wordSpaceNum := 0
	for p.Scan() {
		line := strings.TrimSpace(p.Text())
		lineno++
		switch state {
		case STATE_SENTENCE:
			if line != "" {
				state = STATE_WORDS
				wordSpaceNum = 0
				si = SentenceInfo{Sentence: p.Text()}
			}
		case STATE_WORDS:
			if line == "" {
				wordSpaceNum++
				if wordSpaceNum == 2 {
					ret = append(ret, si)
					state = STATE_SENTENCE
				}
			} else {
				si.Words = append(si.Words,
					WordInfo{word: line, lineno: lineno})
				wordSpaceNum = 0
			}
		}
	}

	if state == STATE_WORDS {
		ret = append(ret, si)
	}
	return ret
}
