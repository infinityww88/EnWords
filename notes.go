package main

import (
	"bufio"
	"io"
	"strings"
)

type NoteInfo struct {
	Note  string
	Words []string
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

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if state == STATE_NOTE {
			if line != "-" {
				noteInfo.Note += line
			} else {
				state = STATE_WORDS
			}
		} else if state == STATE_WORDS {
			if line == "" {
				continue
			} else if line != "--" {
				noteInfo.Words = append(noteInfo.Words, line)
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
