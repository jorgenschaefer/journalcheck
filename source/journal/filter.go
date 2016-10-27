package journal

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/jorgenschaefer/journalcheck/event"
)

type Filter []*regexp.Regexp

func NewFilter(filterfile string) (Filter, error) {
	var filters Filter
	file, err := os.Open(filterfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r, err := regexp.Compile(line)
		if err != nil {
			log.Printf("Bad regular expression, ignoring; err='%s' regexp='%s'", err, line)
		} else {
			filters = append(filters, r)
		}
	}
	return filters, nil
}

func NewEmptyFilter() Filter {
	return *new(Filter)
}

func (f Filter) Matches(ev event.Event) bool {
	line := ev.MatchString()
	for _, r := range f {
		if r.MatchString(line) {
			return true
		}
	}
	return false
}
