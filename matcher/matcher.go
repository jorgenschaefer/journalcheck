package matcher

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Matcher []*regexp.Regexp

func New(filterfile string) (Matcher, error) {
	m := Matcher{}
	f, err := os.Open(filterfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue

		}
		r, err := regexp.Compile(line)
		if err != nil {
			return nil, err
		} else {
			m = append(m, r)
		}
	}
	return m, err
}

func (m *Matcher) Matches(needle string) bool {
	for _, r := range *m {
		if r.MatchString(needle) {
			return true
		}
	}
	return false
}
