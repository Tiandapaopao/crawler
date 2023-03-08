package engine

import (
	"github.com/Tiandapaopao/crawler/collect"
	"github.com/Tiandapaopao/crawler/collector"
	"github.com/Tiandapaopao/crawler/parse/doubanbook"
	"github.com/Tiandapaopao/crawler/parse/doubangroup"
	"github.com/Tiandapaopao/crawler/parse/doubangroupjs"
	"github.com/robertkrimen/otto"
	"go.uber.org/zap"
	"sync"
	"time"
)

func init() {
	Store.Add(doubangroup.DoubangroupTask)
	Store.Add(doubanbook.DoubanBookTask)
	Store.AddJSTask(doubangroupjs.DoubangroupJSTask)
}

// 用于动态规则添加请求。
func AddJsReqs(jreqs []map[string]interface{}) []*collect.Request {
	reqs := make([]*collect.Request, 0)

	for _, jreq := range jreqs {
		req := &collect.Request{}
		u, ok := jreq["Url"].(string)
		if !ok {
			return nil
		}
		req.Url = u
		req.RuleName, _ = jreq["RuleName"].(string)
		req.Method, _ = jreq["Method"].(string)
		req.Priority, _ = jreq["Priority"].(int64)
		reqs = append(reqs, req)
	}
	return reqs
}

// 用于动态规则添加请求。
func AddJsReq(jreq map[string]interface{}) []*collect.Request {
	reqs := make([]*collect.Request, 0)
	req := &collect.Request{}
	u, ok := jreq["Url"].(string)
	if !ok {
		return nil
	}
	req.Url = u
	req.RuleName, _ = jreq["RuleName"].(string)
	req.Method, _ = jreq["Method"].(string)
	req.Priority, _ = jreq["Priority"].(int64)
	reqs = append(reqs, req)
	return reqs
}

func (c *CrawlerStore) AddJSTask(m *collect.TaskModle) {
	task := &collect.Task{
		Property: m.Property,
	}

	task.Rule.Root = func() ([]*collect.Request, error) {
		vm := otto.New()
		vm.Set("AddJsReq", AddJsReqs)
		v, err := vm.Eval(m.Root)
		if err != nil {
			return nil, err
		}
		e, err := v.Export()
		if err != nil {
			return nil, err
		}
		return e.([]*collect.Request), nil
	}

	for _, r := range m.Rules {
		paesrFunc := func(parse string) func(ctx *collect.Context) (collect.ParseResult, error) {
			return func(ctx *collect.Context) (collect.ParseResult, error) {
				vm := otto.New()
				vm.Set("ctx", ctx)
				v, err := vm.Eval(parse)
				if err != nil {
					return collect.ParseResult{}, err
				}
				e, err := v.Export()
				if err != nil {
					return collect.ParseResult{}, err
				}
				if e == nil {
					return collect.ParseResult{}, err
				}
				return e.(collect.ParseResult), err
			}
		}(r.ParseFunc)
		if task.Rule.Trunk == nil {
			task.Rule.Trunk = make(map[string]*collect.Rule, 0)
		}
		task.Rule.Trunk[r.Name] = &collect.Rule{
			ParseFunc: paesrFunc,
		}
	}

	c.Hash[task.Name] = task
	c.list = append(c.list, task)
}

type CrawlerStore struct {
	list []*collect.Task
	Hash map[string]*collect.Task
}

var Store = &CrawlerStore{
	list: []*collect.Task{},
	Hash: map[string]*collect.Task{},
}

func (c *CrawlerStore) Add(task *collect.Task) {
	c.list = append(c.list, task)
	c.Hash[task.Name] = task
}

type Crawler struct {
	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*collect.Request
	failureLock sync.Mutex
	options
}

type Scheduler interface {
	Schedule()
	Push(...*collect.Request)
	Pull() *collect.Request
}

type Schedule struct {
	requestCh   chan *collect.Request
	workerCh    chan *collect.Request
	Logger      *zap.Logger
	reqQueue    []*collect.Request
	priReqQueue []*collect.Request
}

func GetFields(taskName string, ruleName string) []string {
	return Store.Hash[taskName].Rule.Trunk[ruleName].ItemFields
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	e := &Crawler{}
	e.options = options
	e.Visited = make(map[string]bool, 100)
	out := make(chan collect.ParseResult)
	e.out = out
	failures := make(map[string]*collect.Request)
	e.failures = failures
	return e
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	s.requestCh = requestCh
	s.workerCh = workerCh
	return s
}

func (e *Crawler) Schedule() {
	var reqs []*collect.Request
	for _, seed := range e.Seeds {
		task := Store.Hash[seed.Name]
		task.Fetcher = seed.Fetcher
		task.Storage = seed.Storage
		task.Logger = e.Logger
		//rootReqs := task.Rule.Root()
		rootReqs, err := task.Rule.Root()
		if err != nil {
			e.Logger.Error("get root failed",
				zap.Error(err),
			)
			continue
		}
		for _, req := range rootReqs {
			req.Task = task
		}
		reqs = append(reqs, rootReqs...)
	}
	go e.scheduler.Schedule()
	go e.scheduler.Push(reqs...)
}

func (e *Crawler) Run() {
	go e.Schedule()
	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}
	time.Sleep(2 * time.Second)
	e.HandleResult()
}

func (s *Schedule) Schedule() {
	var req *collect.Request
	var ch chan *collect.Request
	for {
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerCh
		}
		if len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			//fmt.Println(123)
			req = nil
			ch = nil
		}
	}
}

func (s *Crawler) CreateWork() {
	for {
		r := s.scheduler.Pull()
		if err := r.Check(); err != nil {
			s.Logger.Error("check failed",
				zap.Error(err),
			)
			continue
		}

		if !r.Task.Reload && s.HasVisited(r) {
			s.Logger.Debug("request has visited",
				zap.String("url:", r.Url),
			)

			continue
		}
		s.StoreVisited(r)
		body, err := r.Task.Fetcher.Get(r)
		//fmt.Println(body)
		//if len(body) < 6000 {
		//	s.Logger.Error("can't fetch ",
		//		zap.Int("length", len(body)),
		//		zap.String("url", r.Url),
		//	)
		//	continue
		//}
		if err != nil {
			s.Logger.Error("can't fetch ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			s.SetFailure(r)
			continue
		}

		rule := r.Task.Rule.Trunk[r.RuleName]

		result, err := rule.ParseFunc(&collect.Context{
			body,
			r,
		})

		if err != nil {
			s.Logger.Error("ParseFunc failed ",
				zap.Error(err),
				zap.String("url", r.Url),
			)
			continue
		}

		if len(result.Requesrts) > 0 {
			go s.scheduler.Push(result.Requesrts...)
		}

		s.out <- result
	}
}

func (s *Crawler) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, item := range result.Items {
				// todo: store
				switch d := item.(type) {
				case *collector.DataCell:
					name := d.GetTaskName()
					task := Store.Hash[name]
					task.Storage.Save(d)
				}
				s.Logger.Sugar().Info("get result: ", item)
			}
		}
	}
}

func (s *Schedule) Push(reqs ...*collect.Request) {
	for _, req := range reqs {
		s.requestCh <- req
	}
}

func (s *Schedule) Pull() *collect.Request {
	r := <-s.workerCh
	return r
}

func (s *Schedule) Output() *collect.Request {
	r := <-s.workerCh
	return r
}

func (e *Crawler) HasVisited(r *collect.Request) bool {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()
	unique := r.Unique()
	return e.Visited[unique]
}

func (e *Crawler) StoreVisited(reqs ...*collect.Request) {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		e.Visited[unique] = true
	}
}

func (e *Crawler) SetFailure(req *collect.Request) bool {
	if !req.Task.Reload {
		e.VisitedLock.Lock()
		unique := req.Unique()
		delete(e.Visited, unique)
		e.VisitedLock.Unlock()
	}
	e.failureLock.Lock()
	defer e.failureLock.Unlock()
	if _, ok := e.failures[req.Unique()]; !ok {
		e.failures[req.Unique()] = req
		e.scheduler.Push(req)
	}
	return true
}
