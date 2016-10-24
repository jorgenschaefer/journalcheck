package filter

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/coreos/go-systemd/sdjournal"

	"github.com/jorgenschaefer/journalcheck/formatter"
)

type RegexpFilter struct {
	filters []*regexp.Regexp
}

func NewRegexpFilter(filterFile string) Filter {
	filter := new(RegexpFilter)
	filter.filters = loadFilters(filterFile)
	return filter
}

func (f *RegexpFilter) Matches(entry *sdjournal.JournalEntry) bool {
	line := formatter.Match(entry)
	for _, r := range f.filters {
		if r.MatchString(line) {
			return true
		}
	}
	return false
}

func loadFilters(filterFile string) []*regexp.Regexp {
	var filters []*regexp.Regexp
	file, err := os.Open(filterFile)
	if err != nil {
		log.Fatal(err)
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

	return filters
}
