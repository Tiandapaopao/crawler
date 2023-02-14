package main

import (
	"fmt"
	"github.com/Tiandapaopao/crawler/collect"
	"github.com/Tiandapaopao/crawler/log"
	"github.com/Tiandapaopao/crawler/parse/douban"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// tag v0.0.5
//var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	//plugin, c := log.NewFilePlugin("./log/log.txt", zapcore.InfoLevel)
	//defer c.Close()
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")
	//proxyURLs := []string{"http://172.19.0.2:8888", "http://172.19.0.2:8888"}
	//p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	//if err != nil {
	//	logger.Error("RoundRobinProxySwitcher failed")
	//}

	//url := "https://www.zhishew.com"
	//var f collect.Fetcher = collect.BrowserFetch{
	//	Timeout: 3000 * time.Millisecond,
	//	//Proxy:   p,
	//}

	cookie := "__utma=30149280.609234316.1675058944.1676280379.1676346325.9; __utmb=30149280.35.5.1676346451768; __utmc=30149280; __utmv=30149280.19775; __utmz=30149280.1676013678.3.2.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/; push_doumail_num=0; push_noty_num=0; _pk_id.100001.8cb4=14ce418015a1a446.1676013677.7.1676346451.1676280376.; _pk_ses.100001.8cb4=*; __gpi=UID=00000bb0bfe53e79:T=1675058959:RT=1676346324:S=ALNI_MYAKZt9MDINtpmS9LfugzL2iJA7uw; ap_v=0,6.0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1676346323%2C%22https%3A%2F%2Ftime.geekbang.org%2F%22%5D; ck=2kCv; ct=y; douban-fav-remind=1; dbcl2=\"197752134:rDFJ+M8i7fk\"; __gads=ID=87e1dea389dd72db-220535d3b2d9003f:T=1676013679:RT=1676013679:S=ALNI_MaT031PxpmNS81fN7vDuyMHqvkbGA; __yadk_uid=Hgf0RfykrqGvtLG6Yejj19BgmSf8ssDr; viewed=\"26416768_1007305_35720728\"; ll=\"118318\"; bid=xYSZqPEUrRo"
	var workList []*collect.Request
	for i := 0; i <= 25; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workList = append(workList, &collect.Request{
			Url:       str,
			ParseFunc: douban.ParseURL,
			Cookie:    cookie,
		})
	}

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		//Proxy:   p,
	}

	items := workList
	workList = nil
	for _, item := range items {
		body, err := f.Get(item)
		time.Sleep(1 * time.Second)
		if err != nil {
			logger.Error("read content failed",
				zap.Error(err),
			)
			continue
		}
		res := item.ParseFunc(body, item)
		for _, v := range res.Requesrts {
			topicStruct := &collect.Request{
				Url:        v.Url,
				ParseTopic: douban.GetContent,
				Cookie:     cookie,
			}
			content, _ := f.Get(topicStruct)
			time.Sleep(1 * time.Second)
			topicUrl := topicStruct.ParseTopic(content, topicStruct.Url)
			if topicUrl != "" {
				logger.Info("result", zap.String("get url:", topicUrl))
			}
		}
	}

}
