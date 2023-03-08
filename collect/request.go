package collect

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/Tiandapaopao/crawler/limiter"
	"github.com/Tiandapaopao/crawler/storage"
	"go.uber.org/zap"
	"math/rand"
	"regexp"
	"sync"
	"time"
)

type Property struct {
	Name     string `json:"name"` // 任务名称，应保证唯一性
	Url      string `json:"url"`
	Cookie   string `json:"cookie"`
	WaitTime int64  `json:"wait_time"`
	Reload   bool   `json:"reload"` // 网站是否可以重复爬取
	MaxDepth int64  `json:"max_depth"`
}

type Task struct {
	Property
	Fetcher     Fetcher
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Storage     storage.Storage
	Rule        RuleTree
	Logger      *zap.Logger
	Limit       limiter.RateLimiter
	//RootReq  *Request
}

type Request struct {
	unique string
	Method string
	Url    string
	//ParseFunc func([]byte, *Request) ParseResult
	Depth    int64
	Priority int64
	Task     *Task
	RuleName string
	TemData  *Temp
}

type Context struct {
	Body []byte
	Req  *Request
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("Max Depth limit")
	}
	return nil
}

func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}

func (c *Context) ParseJSReg(name string, reg string) ParseResult {
	re := regexp.MustCompile(reg)

	matches := re.FindAllSubmatch(c.Body, -1)
	result := ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requesrts = append(
			result.Requesrts, &Request{
				Method:   "GET",
				Task:     c.Req.Task,
				Url:      u,
				Depth:    c.Req.Depth + 1,
				RuleName: name,
			})
	}
	return result
}

func (c *Context) OutputJS(reg string) ParseResult {
	re := regexp.MustCompile(reg)
	ok := re.Match(c.Body)
	if !ok {
		return ParseResult{
			Items: []interface{}{},
		}
	}
	result := ParseResult{
		Items: []interface{}{c.Req.Url},
	}
	return result
}

func (c *Context) Output(data interface{}) *storage.DataCell {
	res := &storage.DataCell{}
	res.Data = make(map[string]interface{})
	res.Data["Task"] = c.Req.Task.Name
	res.Data["rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["Url"] = c.Req.Url
	res.Data["Time"] = time.Now().Format("2006-01-02 15:04:05")
	return res
}

func (r *Request) Fetch() ([]byte, error) {
	if err := r.Task.Limit.Wait(context.Background()); err != nil {
		return nil, err
	}
	// 随机休眠，模拟人类行为
	sleeptime := rand.Int63n(r.Task.WaitTime * 1000)
	time.Sleep(time.Duration(sleeptime) * time.Millisecond)
	return r.Task.Fetcher.Get(r)
}
