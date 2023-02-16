package collect

import (
	"time"
)

type Request struct {
	Url        string
	ParseFunc  func([]byte, *Request) ParseResult
	Cookie     string
	ParseTopic func([]byte, string) string
	WaitTime   time.Duration
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}
