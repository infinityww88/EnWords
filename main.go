package main

import (
	"fmt"
	"log"
	"os"
)

func main0() {

	if e := InitDB(); e != nil {
		log.Fatal("init db failed", e)
		os.Exit(1)
	}

	f, err := os.Open("demo/sentences.txt")
	defer f.Close()

	if err != nil {
		log.Fatal("read sentences failed", err)
		os.Exit(1)
	}
	_ = ReadSentenceInfo(f)
}

func main() {
	fmt.Println(GetOnlineWord("symbol"))
}
