package engine

import (
	"fmt"
	"github.com/Tiandapaopao/crawler/collect"
	"go.uber.org/zap"
	"time"
)

type Schedule struct {
	requestCh chan *collect.Request
	workerCh  chan *collect.Request
	out       chan collect.ParseResult
	options
}

func NewSchedule(opts ...Option) *Schedule {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &Schedule{}
	s.options = options
	return s
}

func (s *Schedule) Run() {
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	out := make(chan collect.ParseResult)
	s.requestCh = requestCh
	s.workerCh = workerCh
	s.out = out
	go s.Schedule()
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}
	time.Sleep(2 * time.Second)
	s.HandleResult()
}

func (s *Schedule) Schedule() {
	var reqQueue []*collect.Request
	for _, seed := range s.Seeds {
		seed.RootReq.Task = seed
		seed.RootReq.Url = seed.Url
		reqQueue = append(reqQueue, seed.RootReq)
	}

	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request

			if len(reqQueue) > 0 {
				req = reqQueue[0]
				//fmt.Println(req.Url)
				reqQueue = reqQueue[1:]
				//time.Sleep(1 * time.Second)
				ch = s.workerCh
			}
			select {
			case r := <-s.requestCh:
				//fmt.Println(r.Url)
				reqQueue = append(reqQueue, r)

			case ch <- req:
				//fmt.Println(req.Url) //条数减少

			}
		}
	}()
}

func (s *Schedule) CreateWork() {
	for {
		r := <-s.workerCh

		if err := r.Check(); err != nil {
			s.Logger.Error("check failed", zap.Error(err))
			continue
		}
		//fmt.Println(r.Url)
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			continue
		}
		fmt.Println("begin parse content : ", r.Url)
		result := r.ParseFunc(body, r)
		//fmt.Printf("there")
		//for _, res := range result.Requesrts {
		//	fmt.Println(res.Url)
		//}
		s.out <- result
	}
}

func (s *Schedule) HandleResult() {
	for {
		select {
		case result := <-s.out:
			//fmt.Println("handle result")
			for _, req := range result.Requesrts {
				//fmt.Println(req.Url)
				s.requestCh <- req
			}
			for _, item := range result.Items {
				//fmt.Println("handle result")
				// todo: store
				s.Logger.Sugar().Info("get result : ", item)
			}
		}
	}
}
