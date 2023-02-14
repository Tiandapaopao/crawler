package collect

type Request struct {
	Url        string
	ParseFunc  func([]byte, *Request) ParseResult
	Cookie     string
	ParseTopic func([]byte, string) string
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}
