package main

import (
	"fmt"
	"net/http"

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
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	Must(err)
	defer resp.Body.Close()
	sel := doc.Find(".mw-parser-output")
	fmt.Println(sel.Text())
}
