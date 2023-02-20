package collect

import (
	"errors"
	"time"
)

type Task struct {
	Url        string
	Cookie     string
	ParseTopic func([]byte, string) string
	WaitTime   time.Duration
	MaxDepth   int
	Fetcher    Fetcher
	RootReq    *Request
}

type Request struct {
	Url       string
	ParseFunc func([]byte, *Request) ParseResult
	Depth     int
	Task      *Task
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
