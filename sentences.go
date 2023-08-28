package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type SentenceInfo struct {
	Sentence string
	Words    []string
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
				if !IsValidWord(line) {
					panic(fmt.Errorf("word `%s` is not valid at line %d", line, lineno))
				}
				si.Words = append(si.Words, line)
				wordSpaceNum = 0
			}
		}
	}

	if state == STATE_WORDS {
		ret = append(ret, si)
	}
	return ret
}
