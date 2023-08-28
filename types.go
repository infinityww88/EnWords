package main

import "fmt"

type Word struct {
	Wid            int64  `gorm:"column:wid;primarykey;autoincrement"`
	Word           string `gorm:"column:word;not null"`
	IsPhrase       bool   `gorm:"column:is_phrase;default:0"`
	PhoneticSymbol string `gorm:"column:phonetic_symbol"`
	Meaning        string `gorm:"column:meaning"`
}

func (self Word) String() string {
	return fmt.Sprintf("Wid\t\t=\t%d\nWord\t\t=\t%s\nIsPhrase\t=\t%t\nPhoneticSymbol\t=\t%s\n%s\n", self.Wid, self.Word, self.IsPhrase, self.PhoneticSymbol, self.Meaning)
}

type Sentence struct {
	Sid      int64  `gorm:"column:sid;primarykey;autoincrement"`
	sentence string `gorm:"column:sentence;not null"`
}

type WordSentence struct {
	WSid int64 `gorm:"column:sid;primarykey;autoincrement"`
	Wid  int64 `gorm:"column:wid;not null"`
	Sid  int64 `gorm:"column:sid;not null"`
}

type Note struct {
	Nid  int64  `gorm:"column:sid;primarykey;autoincrement"`
	note string `gorm:"column:note;not null"`
}

type WordNote struct {
	WNid int64 `gorm:"column:sid;primarykey;autoincrement"`
	Wid  int64 `gorm:"column:wid;not null"`
	Nid  int64 `gorm:"column:nid;not null"`
}

type OEWord struct {
	WordId  int64  `gorm:"column:word_id;primarykey;autoincrement"`
	Word    string `gorm:"column:word;not null"`
	Meaning int64  `gorm:"column:meaning;not null"`
}
