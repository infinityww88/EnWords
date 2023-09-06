package main

import (
	"os"

	"github.com/alexflint/go-arg"
)

type Cmd struct {
	File string `arg:"-f,env:EN_DRAFT" default:"draft.toml"`
}

func main() {
	var cmd Cmd
	arg.MustParse(&cmd)

	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:1080")
	InitDB()
	LoadDataF(cmd.File)
}
