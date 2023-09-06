package main

import "errors"

var WordNotFound = errors.New("Word Not Found")

type Recall struct {
	Word        string `gorm:"column:word;primarykey;not null"`
	RecallScore int32  `gorm:"column:recall_score;default:0"`
	LastTime    uint32 `gorm:"column:last_time;default:0"`
}

type OEWord struct {
	WordExp string `gorm:"column:word_exp"`
	Usage   string `gorm:"column:usage"`
}

func (OEWord) TableName() string {
	return "oe_dict"
}

type ZHWord struct {
	Word    string `gorm:"column:word"`
	Meaning string `gorm:"column:meaning"`
	Phonic  string `gorm:"column:phonic"`
}

type Doc struct {
	Body    string `gorm:"column:body"`
	Title   string `gorm:"column:title"`
	DocType string `gorm:"column:type"`
}
