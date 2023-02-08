package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Tiandapaopao/crawler/collect"
	"github.com/Tiandapaopao/crawler/proxy"
	"time"
)

// tag v0.0.5
//var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	proxyURLs := []string{"http://172.19.0.2:8888", "http://172.19.0.2:8888"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		fmt.Println("RoundRobinProxySwitcher failed")
	}

	url := "https://www.zhishew.com"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	fmt.Println(string(body))

	// 加载HTML文档
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("read content failed:%v", err)
	}

	doc.Find("div.news_li h2 a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		// 获取匹配元素的文本
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
