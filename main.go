package main

import (
	// "fmt"
	"bytes"
	"html/template"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	// "goblog/helper"
	"goblog/route"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	g := gin.New()

	g.SetFuncMap(helperFuncs)

	g.LoadHTMLGlob(filepath.Join("", "./view/*"))
	g.Static("/static", filepath.Join("", "./static"))

	route.Router(g)

	g.Run(":8090")
}

// template function
var helperFuncs = template.FuncMap{
	"string": func(b []byte) string {
		return string(b)
	},
	"datetime": func(format string, param ...string) string {
		var duration string
		if len(param) > 0 {
			duration = param[0]
		}

		t, _ := strconv.ParseInt(duration, 10, 64)
		return time.Unix(t, 0).Format(format)
	},
	"htmlMore": func(b []byte) bool {
		return bytes.Contains(b, []byte("<!--more-->"))
	},
	"htmlLess": func(b []byte) template.HTML {

		b = bytes.Replace(b, []byte("<!--markdown-->"), []byte(""), -1)

		unsafe := blackfriday.MarkdownCommon(b)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

		return template.HTML(string(html))
	},
}
