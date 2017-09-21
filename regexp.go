package goutils

import (
	"strings"
)

// Matcher represents the result of succesful match
type Matcher interface {
	// Groups returns the number of capturing groups in the pattern
	Groups() int
	// GroupByIdx returns the value of a group by index
	// Important: GroupByIdx(0) returns the full match, the first group has index 1
	GroupByIdx(idx int) string
	// GroupPresentByIdx returns whether the specified group is present in the succesful match
	GroupPresentByIdx(idx int) bool
	// GroupByName returns the value of a group by name
	GroupByName(name string) string
	// GroupPresentByName returns whether the specified group is present in the succesful match
	GroupPresentByName(name string) bool
	// Next tries to find a next match, returns true if it succeeds;
	// after calling this method all info about the previous match gets lost
	Next() bool
}

// Regexp represents compiled regular expression
type Regexp interface {
	// Search returns a matcher for the leftmost successful match or nil if the pattern failed to match
	Search(subject string) Matcher
	// Replace replaces every occurrence of the pattern with repl
	// (group references in repl are in Perl/Java style: $0, $1, $2, ..)
	Replace(subject, repl string) string
}

// Engine represents a regexp engine
type Engine interface {
	Compile(pattern string, flags int) (Regexp, error)
	Quote(s string) string

	FlagDotAll() int
}

// ReEngine contains the regexp engine which should be used for all matching
// For local dev runs init this with 'regexppcre.NewEngine()'; from Java it will be inited on app-start
var ReEngine Engine

// ReTest runs a quick match, it should be 23,7. What else?
func ReTest() string {
	re, err := ReEngine.Compile("(?<!redherring )'(?<num>[0-9]+)'", 0)
	if err != nil {
		return "errored"
	}
	m := re.Search("My favourite redherring '7' number is '23' and also '7'. Yeah!")
	var nums []string
	if m != nil {
		for {
			nums = append(nums, m.GroupByName("num"))
			if !m.Next() {
				break
			}
		}
	}
	return strings.Join(nums, ",")
}
