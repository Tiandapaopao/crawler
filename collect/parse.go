package collect

type RuleTree struct {
	Root  func() []*Request
	Trunk map[string]*Rule
}

type Rule struct {
	ParseFunc func(*Context) ParseResult
}
