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

func main() {

	p := arg.MustParse(&args)

	switch {
	case args.BatchLoad != nil:
		switch args.BatchLoad.LoadType {
		case "sentences":
			fmt.Println("batch load sentences")
		case "notes":
			fmt.Println("batch load notes")
		default:
			p.Fail("batch load must be \"sentences\" or \"notes\"")
		}
	case args.InsertWord != nil:
		fmt.Printf("insert word %s %s %s\n", args.InsertWord.Word,
			args.InsertWord.Phonetic, args.InsertWord.Meaning)
	}

	if e := InitDB(); e != nil {
		log.Fatal("init db failed", e)
		os.Exit(1)
	}

	/*
		f, _ := os.Open("demo/sentences.txt")
		sis := ReadSentenceInfo(f)
		for _, si := range sis {
			fmt.Println(si.Sentence)
			for _, w := range si.Words {
				fmt.Println("\t", w)
			}
			fmt.Println(strings.Repeat("-", 10))
		}
		LoadSentenceInfo(sis)
	*/
	/*
		f, _ := os.Open("demo/notes.txt")
		nis := ReadNoteInfo(f)
		LoadNoteInfo(nis)
	*/
}
