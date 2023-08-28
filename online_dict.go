package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmespath/go-jmespath"
)

var v = `curl -XPOST --data 'q=symbol&le=en&t=3&client=web&keyfrom=webdict' 'https://dict.youdao.com/jsonapi_s?doctype=json&jsonversion=4'`

var dictUrl = "https://dict.youdao.com/jsonapi_s?doctype=json&jsonversion=4"

var soundUrl = "https://dict.youdao.com/dictvoice"

var WordNotFound = fmt.Errorf("WordNotFound")

func GetOnlineWord(word string) (Word, error) {
	postdata := url.Values{
		"q":       {word},
		"le":      {"en"},
		"t":       {"3"},
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

	tran, _ := jmespath.Search("ec.word.trs[*].join(' ', [pos||'', tran]) | join(', ', @)", data)

	if tran == nil {
		return Word{}, WordNotFound
	}

	phone, _ := jmespath.Search("ec.word.[usphone||ukphone][0]", data)

	return Word{Word: word, IsPhrase: strings.Contains(word, " "), PhoneticSymbol: phone.(string), Meaning: tran.(string)}, nil
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
