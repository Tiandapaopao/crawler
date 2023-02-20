package main

import (
	"fmt"
	"github.com/Tiandapaopao/crawler/collect"
	"github.com/Tiandapaopao/crawler/engine"
	"github.com/Tiandapaopao/crawler/log"
	"github.com/Tiandapaopao/crawler/parse/douban"
	"go.uber.org/zap/zapcore"
	"time"
)

// tag v0.0.5
//var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {

	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	cookie := "__utma=30149280.609234316.1675058944.1676360711.1676443276.12; __utmb=30149280.2.10.1676443276; __utmc=30149280; __utmt=1; __utmv=30149280.19775; __utmz=30149280.1676013678.3.2.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/; push_doumail_num=0; push_noty_num=0; __gpi=UID=00000bb0bfe53e79:T=1675058959:RT=1676443268:S=ALNI_MYAKZt9MDINtpmS9LfugzL2iJA7uw; ap_v=0,6.0; _pk_id.100001.8cb4=14ce418015a1a446.1676013677.10.1676443265.1676360727.; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1676443265%2C%22https%3A%2F%2Ftime.geekbang.org%2F%22%5D; _pk_ses.100001.8cb4=*; ck=2kCv; ct=y; douban-fav-remind=1; dbcl2=\"197752134:rDFJ+M8i7fk\"; __gads=ID=87e1dea389dd72db-220535d3b2d9003f:T=1676013679:RT=1676013679:S=ALNI_MaT031PxpmNS81fN7vDuyMHqvkbGA; __yadk_uid=Hgf0RfykrqGvtLG6Yejj19BgmSf8ssDr; viewed=\"26416768_1007305_35720728\"; ll=\"118318\"; bid=xYSZqPEUrRo"

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 10000 * time.Millisecond,
		Logger:  logger,
	}
	var seeds = make([]*collect.Task, 0, 1000)
	for i := 0; i <= 0; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &collect.Task{
			Url:      str,
			WaitTime: 1 * time.Second,
			Cookie:   cookie,
			MaxDepth: 5,
			Fetcher:  f,
			RootReq: &collect.Request{
				Method:    "GET",
				ParseFunc: douban.ParseURL,
			},
		})
	}

	s := engine.NewEngine(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)
	s.Run()
}
