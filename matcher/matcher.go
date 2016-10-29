// Package matcher implements a file-based regexp matcher. A list of
// regular expressions are loaded from a file, one per line, and then
// strings can be checked whether they match any of the regular
// expressions.
package matcher

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Matcher []*regexp.Regexp

// New returns an object that can be used to match against regular
// expressions in the given filterfile.
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

// Matches returns true if any regular expression in the matcher
// matches the needle string.
func (m *Matcher) Matches(needle string) bool {
	for _, r := range *m {
		if r.MatchString(needle) {
			return true
		}
	}
	return false
}
