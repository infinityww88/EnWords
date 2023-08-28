package main

import "fmt"

type Word struct {
	Wid            int64  `gorm:"column:wid;primarykey;autoincrement"`
	Word           string `gorm:"column:word;unique;not null"`
	IsPhrase       bool   `gorm:"column:is_phrase;default:0"`
	PhoneticSymbol string `gorm:"column:phonetic_symbol"`
	Meaning        string `gorm:"column:meaning"`
}

func (self Word) String() string {
	return fmt.Sprintf("Wid\t\t=\t%d\nWord\t\t=\t%s\nIsPhrase\t=\t%t\nPhoneticSymbol\t=\t%s\n%s\n", self.Wid, self.Word, self.IsPhrase, self.PhoneticSymbol, self.Meaning)
}

type Sentence struct {
	Sid      int64  `gorm:"column:sid;primarykey;autoincrement"`
	Sentence string `gorm:"column:sentence;not null"`
	Digit    string `gorm:"column:digit;unique;not null"`
}

type WordSentence struct {
	WSid int64 `gorm:"column:wsid;primarykey;autoincrement"`
	Wid  int64 `gorm:"column:wid;not null;uniqueIndex:ws_idx"`
	Sid  int64 `gorm:"column:sid;not null;uniqueIndex:ws_idx"`
}

type Note struct {
	Nid   int64  `gorm:"column:nid;primarykey;autoincrement"`
	Note  string `gorm:"column:note;not null"`
	Digit string `gorm:"column:digit;unique;not null"`
}

type WordNote struct {
	WNid int64 `gorm:"column:wnid;primarykey;autoincrement"`
	Wid  int64 `gorm:"column:wid;not null;uniqueIndex:wn_idx"`
	Nid  int64 `gorm:"column:nid;not null;uniqueIndex:wn_idx"`
}

type OEWord struct {
	WordId  int64  `gorm:"column:word_id;primarykey;autoincrement"`
	Word    string `gorm:"column:word;not null"`
	Meaning int64  `gorm:"column:meaning;not null"`
}
