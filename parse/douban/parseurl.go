package douban

import (
	"fmt"
	"github.com/Tiandapaopao/crawler/collect"
	"regexp"
)

//const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
//
//func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
//	re := regexp.MustCompile(urlListRe)
//
//	matches := re.FindAllSubmatch(contents, -1)
//	result := collect.ParseResult{}
//
//	i := 0
//	for _, m := range matches {
//		//fmt.Println(i)
//		u := string(m[1])
//		//fmt.Println(u)
//		result.Requesrts = append(
//			result.Requesrts, &collect.Request{
//				Task:   req.Task,
//				Url:    u,
//				Depth:  req.Depth + 1,
//				Method: "GET",
//				//ParseFunc: func(c []byte, request *collect.Request) collect.ParseResult {
//				//	return GetContent(c, u)
//				//},
//			})
//		i++
//	}
//	return result
//}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div class="aside">`

func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)
	ok := re.Match(contents)
	fmt.Println("parse!!!")
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
