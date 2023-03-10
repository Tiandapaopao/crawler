package doubanbook

import (
	"github.com/Tiandapaopao/crawler/collect"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

var DoubanBookTask = &collect.Task{
	Property: collect.Property{
		Name:     "douban_book_list",
		WaitTime: 20,
		MaxDepth: 5,
		Cookie:   "__utma=30149280.609234316.1675058944.1676360711.1676443276.12; __utmb=30149280.2.10.1676443276; __utmc=30149280; __utmt=1; __utmv=30149280.19775; __utmz=30149280.1676013678.3.2.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/; push_doumail_num=0; push_noty_num=0; __gpi=UID=00000bb0bfe53e79:T=1675058959:RT=1676443268:S=ALNI_MYAKZt9MDINtpmS9LfugzL2iJA7uw; ap_v=0,6.0; _pk_id.100001.8cb4=14ce418015a1a446.1676013677.10.1676443265.1676360727.; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1676443265%2C%22https%3A%2F%2Ftime.geekbang.org%2F%22%5D; _pk_ses.100001.8cb4=*; ck=2kCv; ct=y; douban-fav-remind=1; dbcl2=\"197752134:rDFJ+M8i7fk\"; __gads=ID=87e1dea389dd72db-220535d3b2d9003f:T=1676013679:RT=1676013679:S=ALNI_MaT031PxpmNS81fN7vDuyMHqvkbGA; __yadk_uid=Hgf0RfykrqGvtLG6Yejj19BgmSf8ssDr; viewed=\"26416768_1007305_35720728\"; ll=\"118318\"; bid=xYSZqPEUrRo",
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				&collect.Request{
					Priority: 1,
					Url:      "https://book.douban.com",
					Method:   "GET",
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag":  &collect.Rule{ParseFunc: ParseTag},
			"书籍列表": &collect.Rule{ParseFunc: ParseBookList},
			"书籍简介": &collect.Rule{
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

func ParseTag(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(regexpStr)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		result.Requesrts = append(result.Requesrts, &collect.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      "https://book.douban.com" + string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍列表",
		})
	}
	//result.Requesrts = result.Requesrts[:1]
	return result, nil
}

const BookListRe = `<a.*?href="([^"]+)" title="([^"]+)"`

func ParseBookList(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(BookListRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		req := &collect.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Priority: 100,
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TemData = &collect.Temp{}
		req.TemData.Set("book_name", string(m[2]))
		result.Requesrts = append(result.Requesrts, req)
	}
	//result.Requesrts = result.Requesrts[:1]
	zap.S().Debugln("parse book list,count:", len(result.Requesrts))
	return result, nil

}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>[\d\D]*?<a.*?>([^<]+)</a>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

func ParseBookDetail(ctx *collect.Context) (collect.ParseResult, error) {
	bookName := ctx.Req.TemData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))
	book := map[string]interface{}{
		"书名":   bookName,
		"作者":   ExtraString(ctx.Body, autoRe),
		"页数":   page,
		"出版社": ExtraString(ctx.Body, public),
		"得分":   ExtraString(ctx.Body, scoreRe),
		"价格":   ExtraString(ctx.Body, priceRe),
		"简介":   ExtraString(ctx.Body, intoRe),
	}

	data := ctx.Output(book)
	result := collect.ParseResult{
		Items: []interface{}{data},
	}
	//result.Requesrts = result.Requesrts[:3]
	zap.S().Debugln("parse book detail", data)
	return result, nil
}

func ExtraString(content []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(content)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
