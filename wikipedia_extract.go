package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func ExtractWikipedia(url string) {
	resp, err := http.Get(url)
	Must(err)
	defer resp.Body.Close()
	ParseWikipedia(resp.Body)
}

func ParseWikipedia(r io.Reader) {
	doc, err := goquery.NewDocumentFromReader(r)
	Must(err)
	sel := doc.Find(".mw-parser-output h2,h3,h4,p")
	re, err := regexp.Compile(`\[(\d+|edit)\]`)
	Must(err)
	for i := 0; i < sel.Length(); i++ {
		s := sel.Slice(i, i+1)
		text := re.ReplaceAllString(s.Text(), "")
		if goquery.NodeName(s) == "h2" && text == "See also" {
			break
		}
		fmt.Println(text)
	}
}
