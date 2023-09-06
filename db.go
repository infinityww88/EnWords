package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dictDB   *gorm.DB
	recallDB *gorm.DB
	docsDB   *gorm.DB
)

func InitDB() {
	var err error
	dictDB, err = gorm.Open(sqlite.Open("db/dict.db"), &gorm.Config{})
	Must(err)
	recallDB, err = gorm.Open(sqlite.Open("db/recall.db"), &gorm.Config{})
	Must(err)
	docsDB, err = gorm.Open(sqlite.Open("db/docs.db"), &gorm.Config{})
	Must(err)
}

func GetRecall(word string) (Recall, error) {
	recall := Recall{}
	ret := dictDB.Table("recall").Where(`where word = ?`, word).First(&recall)
	if ret.RowsAffected == 0 {
		return recall, WordNotFound
	}
	return recall, nil
}

func GetOEWord(word string) []OEWord {
	oeword := OEWord{}
	ret := []OEWord{}
	sql := fmt.Sprintf("SELECT * FROM oe_dict WHERE oe_dict MATCH 'word_exp:^\"%s\"'", word)
	rows, err := dictDB.Raw(sql).Rows()
	Must(err)
	defer rows.Close()
	for rows.Next() {
		dictDB.ScanRows(rows, &oeword)
		ret = append(ret, oeword)
	}

	return ret
}

func GetZHWord(word string) (ZHWord, error) {
	zhword := ZHWord{}
	ret := dictDB.Table("zh_dict").Where(`word = ?`, word).First(&zhword)
	if ret.RowsAffected == 0 {
		return zhword, WordNotFound
	}
	return zhword, nil
}

func InsertRecallAlways(word string) bool {
	recall := Recall{Word: word}
	ret := recallDB.Table("recall").FirstOrCreate(&recall)
	return ret.RowsAffected > 0
}

func InsertDocs(body, title, docType string) {
	doc := Doc{
		Body:    body,
		Title:   title,
		DocType: docType,
	}
	docsDB.Create(&doc)
}

var DocNotFound = fmt.Errorf("doc not found")

func GetDoc(title, docType string) (Doc, error) {
	doc := Doc{
		Title:   title,
		DocType: docType,
	}
	ret := docsDB.Where(&doc).First(&doc)
	if ret.RowsAffected == 0 {
		return doc, DocNotFound
	}
	return doc, nil
}
