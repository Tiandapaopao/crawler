package douban

import (
	"fmt"
	"github.com/Tiandapaopao/crawler/collect"
	"regexp"
	"time"
)

const urlListRe1 = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
const ContentRe1 = `<div class="topic-content">[\s\S]*?地铁[\s\S]*?<div class="aside">`

var DoubangroupTask = &collect.Task{
	Property: collect.Property{
		Name:     "find_douban_sun_room",
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "__utma=30149280.609234316.1675058944.1676360711.1676443276.12; __utmb=30149280.2.10.1676443276; __utmc=30149280; __utmt=1; __utmv=30149280.19775; __utmz=30149280.1676013678.3.2.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/; push_doumail_num=0; push_noty_num=0; __gpi=UID=00000bb0bfe53e79:T=1675058959:RT=1676443268:S=ALNI_MYAKZt9MDINtpmS9LfugzL2iJA7uw; ap_v=0,6.0; _pk_id.100001.8cb4=14ce418015a1a446.1676013677.10.1676443265.1676360727.; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1676443265%2C%22https%3A%2F%2Ftime.geekbang.org%2F%22%5D; _pk_ses.100001.8cb4=*; ck=2kCv; ct=y; douban-fav-remind=1; dbcl2=\"197752134:rDFJ+M8i7fk\"; __gads=ID=87e1dea389dd72db-220535d3b2d9003f:T=1676013679:RT=1676013679:S=ALNI_MaT031PxpmNS81fN7vDuyMHqvkbGA; __yadk_uid=Hgf0RfykrqGvtLG6Yejj19BgmSf8ssDr; viewed=\"26416768_1007305_35720728\"; ll=\"118318\"; bid=xYSZqPEUrRo",
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			var roots []*collect.Request
			for i := 0; i <= 50; i += 25 {
				str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
				roots = append(roots, &collect.Request{
					Priority: 1,
					Url:      str,
					Method:   "GET",
					RuleName: "解析网站URL",
				})
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"解析网站URL": &collect.Rule{ParseURL},
			"解析阳台房":  &collect.Rule{GetSunRoom},
		},
	},
}

func ParseURL(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(urlListRe1)

	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requesrts = append(
			result.Requesrts, &collect.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      u,
				Depth:    ctx.Req.Depth + 1,
				RuleName: "解析阳台房",
			})
	}
	return result, nil
}

func GetSunRoom(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(ContentRe1)

	ok := re.Match(ctx.Body)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}, nil
	}
	result := collect.ParseResult{
		Items: []interface{}{ctx.Req.Url},
	}
	return result, nil
}
