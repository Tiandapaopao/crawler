package main

import (
	"fmt"
	"github.com/Tiandapaopao/crawler/collect"
	"regexp"
)

// tag v0.0.5
var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	url := "https://book.douban.com/subject/26416768/"
	var f collect.Fetcher = collect.BrowserFetch{}
	body, err := f.Get(url)

	fmt.Println(string(body))

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	matches := headerRe.FindAllSubmatch(body, -1)
	for _, m := range matches {
		fmt.Println("fetch card news:", string(m[1]))
	}
}
