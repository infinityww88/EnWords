package main

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func ServerMain() {
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
