package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
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

func TryWiki() {
	/*
		os.Setenv("HTTPS_PROXY", "http://127.0.0.1:1080")
		ExtractWikipedia("https://en.wikipedia.org/wiki/NSA")
	*/
	/*
		f, err := os.Open("wiki.html")
		Must(err)
		defer f.Close()
		ParseWikipedia(f)
	*/
}

func MyPlugin() gin.HandlerFunc {
	return func(c *gin.Context) {
		dir := path.Dir(c.Request.URL.Path)
		parts := strings.Split(dir, "/")
		if parts[1] == "assets" {
			path := c.Request.URL.Path
			if ".gz" == path[len(path)-3:] {
				c.Header("Content-Encoding", "gzip")
				if len(path) >= 8 && ".wasm.gz" == path[len(path)-8:] {
					c.Header("Content-Type", "application/wasm")
				} else {
					c.Header("Content-Type", "application/x-gzip")
				}
			}
		}
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(MyPlugin())
	r.GET("/api", func(c *gin.Context) {
		data := map[string]any{
			"lang": "golang",
			"tag":  "<br>",
		}
		c.AsciiJSON(http.StatusOK, data)
	})
	r.Static("/assets", "./assets")
	// remove server.key password
	r.RunTLS(":8083", "./server.crt", "./server.key")
	//r.Run(":8080")
}
