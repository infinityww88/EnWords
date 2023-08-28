package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("db/english.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	db.AutoMigrate(&Word{})
	db.AutoMigrate(&Sentence{})
	db.AutoMigrate(&WordSentence{})
	db.AutoMigrate(&Note{})
	db.AutoMigrate(&WordNote{})
	return nil
}

func GetWord(word string) (Word, bool) {
	w := Word{}
	result := db.Model(&Word{}).Where(Word{Word: word}).First(&w)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return w, false
	}
	return w, true
}

func InsertWordAlways(word Word) Word {
	ret := Word{}
	word.IsPhrase = strings.Contains(word.Word, " ")
	db.Where(Word{Word: word.Word}).Attrs(word).FirstOrCreate(&ret)
	return ret
}

func IsValidWord(word string) bool {
	var count int64
	db.Model(&OEWord{}).Where("word = ?", word).Count(&count)
	return count > 0
}

func CheckInvalidWords(words []string) []string {
	ret := []string{}
	for _, w := range words {
		if !IsValidWord(w) {
			ret = append(ret, w)
		}
	}
	return ret
}

func InsertSentenceAlways(sentence string) int64 {
	ret := Sentence{}
	db.Where(Sentence{Digit: digit(sentence)}).
		Attrs(Sentence{Sentence: sentence}).
		FirstOrCreate(&ret)
	return ret.Sid
}

/*
func InsertSentenceInfo(si *SentenceInfo) {
	sid := InsertSentenceAlways(si.Sentence)
	for _, w := range si.Words {
		wid := InsertWordAlways(w)
		db.Create(WordSentence{Wid: wid, Sid: sid})
	}
}
*/

func digit(content string) string {
	d := md5.Sum([]byte(content))
	return hex.EncodeToString(d[:])
}

func InsertNoteAlways(note string) int64 {
	ret := Note{}
	db.Where(Note{Digit: digit(note)}).
		Attrs(Note{Note: note}).
		FirstOrCreate(&ret)
	return ret.Nid
}

/*
func InsertNoteInfo(ni *NoteInfo) {
	nid := InsertNoteAlways(ni.Note)
	for _, w := range ni.Words {
		wid := InsertWordAlways(w)
		db.Create(WordNote{Wid: wid, Nid: nid})
	}
}
*/

func InsertWordSentenceAlways(wid int64, sid int64) int64 {
	ws := WordSentence{}
	db.FirstOrCreate(&ws, WordSentence{Wid: wid, Sid: sid})
	return ws.WSid
}

func InsertWordNoteAlways(wid int64, nid int64) int64 {
	wn := WordNote{}
	db.FirstOrCreate(&wn, WordNote{Wid: wid, Nid: nid})
	return wn.WNid
}

func UpdateWord(word string, meaning string) int64 {
	ret := db.Model(&Word{}).Where("word = ?", word).Update("meaning", meaning)
	return ret.RowsAffected
}

func UpdateSentence(sid int64, sentence string) int64 {
	ret := db.Model(&Sentence{}).Where("sid = ?", sid).Updates(Sentence{Sentence: sentence, Digit: digit(sentence)})
	return ret.RowsAffected
}

func UpdateNote(nid int64, note string) int64 {
	ret := db.Model(&Note{}).Where("nid = ?", nid).Updates(Note{Note: note, Digit: digit(note)})
	return ret.RowsAffected
}
