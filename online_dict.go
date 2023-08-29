package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jmespath/go-jmespath"
	"golang.org/x/time/rate"
)

var v = `curl -XPOST --data 'q=symbol&le=en&t=3&client=web&keyfrom=webdict' 'https://dict.youdao.com/jsonapi_s?doctype=json&jsonversion=4'`

var dictUrl = "https://dict.youdao.com/jsonapi_s?doctype=json&jsonversion=4"

var dictPageUrl = "https://www.youdao.com/result?word=%s&lang=en"

var soundUrl = "https://dict.youdao.com/dictvoice"

var WordNotFound = errors.New("Word Not Found")

var limiter = rate.NewLimiter(2, 1)

func GetOnlineWord(word string) (Word, error) {
	if w, err := GetOnlineWordApi(word); err != nil {
		return GetOnlineWordOnPage(word)
	} else {
		return w, nil
	}
}

func GetOnlineWordOnPage(word string) (Word, error) {
	link := fmt.Sprintf(dictPageUrl, url.QueryEscape(word))
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Get word on page failed with status code %d", resp.StatusCode))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(fmt.Errorf("Parse document for word %s failed", word))
	}
	ret := []string{}
	doc.Find(".catalogue_author .trans-container .word-exp").Each(func(i int, s *goquery.Selection) {
		ret = append(ret, s.Text())
	})

	meaning := strings.TrimSpace(strings.Join(ret, "\n"))
	if meaning == "" {
		return Word{}, WordNotFound
	}

	sel := doc.Find(".trans-container .phonetic")
	phonetic := ""
	if len(sel.Nodes) > 0 {
		sel.Nodes = sel.Nodes[len(sel.Nodes)-1:]
		phonetic = sel.Text()
	}
	return Word{
		Word:           word,
		IsPhrase:       strings.Contains(word, " "),
		PhoneticSymbol: phonetic,
		Meaning:        meaning,
	}, nil
}

func GetOnlineWordApi(word string) (Word, error) {
	limiter.Wait(context.Background())
	postdata := url.Values{
		"q":       {word},
		"le":      {"en"},
		"client":  {"web"},
		"keyfrom": {"webdict"}}

	resp, err := http.PostForm(dictUrl, postdata)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data any
	json.Unmarshal(body, &data)

	jpath := "ec.word.trs[*].join(' ', [pos||'', tran]) | join(', ', @)"
	tran, _ := jmespath.Search(jpath, data)

	if tran == nil {
		return Word{}, WordNotFound
	}

	phone, _ := jmespath.Search("ec.word.[usphone||ukphone][0]", data)

	return Word{
			Word:           word,
			IsPhrase:       strings.Contains(word, " "),
			PhoneticSymbol: phone.(string),
			Meaning:        tran.(string)},
		nil
}

func DownloadWordSound(word string) {
	u := fmt.Sprintf("%s?audio=%s&type=2", soundUrl, url.QueryEscape(word))
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	destFile, _ := os.Create(filepath.Join("sound", strings.ReplaceAll(word, " ", "_")+".mp3"))
	io.Copy(destFile, resp.Body)
}
