package main

import (
	"errors"
	"log"

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

func WordExists(word string) bool {
	w := Word{}
	result := db.Model(&Word{}).Where(Word{Word: word}).First(&w)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func InsertWordIfNotExist(word string) int64 {
	w := Word{}
	db.FirstOrCreate(&w, Word{Word: word})
	return w.Wid
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

func InsertSentence(sentence string) int64 {
	s := Sentence{sentence: sentence}
	db.Create(&s)
	return s.Sid
}

func InsertSentenceInfo(si *SentenceInfo) {
	sid := InsertSentence(si.Sentence)
	for _, w := range si.Words {
		wid := InsertWordIfNotExist(w)
		db.Create(WordSentence{Wid: wid, Sid: sid})
	}
}

func InsertNote(note string) int64 {
	n := Note{note: note}
	db.Create(&n)
	return n.Nid
}

func InsertNoteInfo(ni *NoteInfo) {
	nid := InsertNote(ni.Note)
	for _, w := range ni.Words {
		wid := InsertWordIfNotExist(w)
		db.Create(WordNote{Wid: wid, Nid: nid})
	}
}
