package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type NoteInfo struct {
	Note  string
	Words []WordInfo
}

func LoadNoteInfo(nis []NoteInfo) {
	for _, ni := range nis {
		nid := InsertNoteAlways(ni.Note)
		for _, wi := range ni.Words {
			var w Word
			var err error
			var ok bool
			if w, ok = GetWord(wi.word); !ok {
				if w, err = LoadWord(wi.word); err != nil {
					panic(fmt.Errorf("load word %s error %w at line [%d]", wi.word, err, wi.lineno))
				}
			}
			InsertWordNoteAlways(w.Wid, nid)
		}
	}
}

func ReadNoteInfo(reader io.Reader) []NoteInfo {
	const (
		STATE_NOTE  = "note"
		STATE_WORDS = "words"
	)

	state := STATE_NOTE
	scanner := bufio.NewScanner(reader)

	var noteInfo NoteInfo
	var ret []NoteInfo
	lineno := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineno++
		if state == STATE_NOTE {
			if line != "-" {
				noteInfo.Note += line + "\n"
			} else {
				state = STATE_WORDS
			}
		} else if state == STATE_WORDS {
			if line == "" {
				continue
			} else if line != "--" {
				noteInfo.Words = append(noteInfo.Words, WordInfo{word: line, lineno: lineno})
			} else {
				state = STATE_NOTE
				ret = append(ret, noteInfo)
				noteInfo = NoteInfo{}
			}
		}
	}

	if state == STATE_WORDS {
		ret = append(ret, noteInfo)
	}

	return ret
}
