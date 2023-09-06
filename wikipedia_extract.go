package main

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func ExtractWikipedia(url string) string {
	resp, err := http.Get(url)
	Must(err)
	defer resp.Body.Close()
	return ParseWikipedia(resp.Body)
}

func ParseWikipedia(r io.Reader) string {
	doc, err := goquery.NewDocumentFromReader(r)
	Must(err)
	sel := doc.Find(".mw-parser-output h2,h3,h4,p")
	re, err := regexp.Compile(`\[(\d+|edit)\]`)
	Must(err)
	sb := strings.Builder{}
	for i := 0; i < sel.Length(); i++ {
		s := sel.Slice(i, i+1)
		text := re.ReplaceAllString(s.Text(), "")
		if goquery.NodeName(s) == "h2" && text == "See also" {
			break
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}
	return sb.String()
}
