package douban

import (
	"github.com/Tiandapaopao/crawler/collect"
	"regexp"
)

const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	re := regexp.MustCompile(urlListRe)

	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		//fmt.Println(u)
		result.Requesrts = append(
			result.Requesrts, &collect.Request{
				Url:    u,
				Cookie: req.Cookie,
				ParseFunc: func(c []byte, request *collect.Request) collect.ParseResult {
					return GetContent(c, u)
				},
			})
	}
	return result
}

const ContentRe = `<div class="topic-content">[\s\S]*?商场[\s\S]*?<div`

func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)
	//fmt.Println(url)
	ok := re.Match(contents)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}

	result := collect.ParseResult{
		Items: []interface{}{url},
	}

	return result
}
