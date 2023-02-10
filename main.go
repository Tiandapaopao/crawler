package main

import (
	"github.com/Tiandapaopao/crawler/collect"
	"github.com/Tiandapaopao/crawler/log"
	"github.com/Tiandapaopao/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// tag v0.0.5
//var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	plugin, c := log.NewFilePlugin("./log/log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")
	proxyURLs := []string{"http://172.19.0.2:8888", "http://172.19.0.2:8888"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}

	url := "https://www.zhishew.com"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		logger.Error("read content failed",
			zap.Error(err),
		)
		return
	}
	//fmt.Println(string(body))

	logger.Info("get content", zap.Int("len", len(body)))
}
