package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
)

type InsertWordCmd struct {
	Word     string `arg:"-w,required"`
	Phonetic string `arg:"-p,required"`
	Meaning  string `arg:"-m,required"`
}

type BatchLoadCmd struct {
	LoadType string `arg:"required,positional"`
}

var args struct {
	InsertWord *InsertWordCmd `arg:"subcommand:insertword"`
	BatchLoad  *BatchLoadCmd  `arg:"subcommand:batchload"`
}

func loadSentences() {
	if f, err := os.Open("draft/sentences.txt"); err == nil {
		sis := ReadSentenceInfo(f)
		LoadSentenceInfo(sis)
	} else {
		panic(err)
	}
}

func loadNotes() {
	if f, err := os.Open("draft/notes.txt"); err == nil {
		nis := ReadNoteInfo(f)
		LoadNoteInfo(nis)
	} else {
		panic(err)
	}
}

func entry() {

	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("no subcommand specifed")
	}

	if e := InitDB(); e != nil {
		log.Fatal("init db failed", e)
		os.Exit(1)
	}

	switch {
	case args.BatchLoad != nil:
		switch args.BatchLoad.LoadType {
		case "sentences":
			loadSentences()
		case "notes":
			loadNotes()
		default:
			p.Fail("batch load must be \"sentences\" or \"notes\"")
		}
	case args.InsertWord != nil:
		w := Word{
			Word:           args.InsertWord.Word,
			PhoneticSymbol: args.InsertWord.Phonetic,
			Meaning:        args.InsertWord.Meaning}
		w = InsertWordAlways(w)
		fmt.Printf("insert at id %d\n", w.Wid)
	}
}

func main() {
	ExtractWikipedia("https://en.wikipedia.org/wiki/World_War_II")
}
