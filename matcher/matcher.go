package matcher

type Matcher struct{}

func New(filterfile string) *Matcher {
	return &Matcher{}
}

func (m *Matcher) Matches(needle string) (bool, error) {
	return false, nil
}
