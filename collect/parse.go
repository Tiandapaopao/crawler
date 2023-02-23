package collect

type RuleTree struct {
	Root  func() ([]*Request, error)
	Trunk map[string]*Rule
}

type Rule struct {
	ParseFunc func(*Context) (ParseResult, error)
}
