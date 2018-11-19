/**
 * Go PCRE Matcher
 *
 *    Copyright 2017 Tenta, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * For any questions, please contact developer@tenta.io
 *
 * matcher.go: Matcher interface
 */

// Package matcher provides a simple and compatible matcher interface around [Go PCRE](https://github.com/gijsbers/go-pcre).
package matcher

import (
	"strings"
)

// Match represents the result of a successful match
type Match interface {
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
	// Search returns the leftmost successful match or nil if the pattern failed to match
	Search(subject string) Match
	// Replace replaces every occurrence of the pattern with repl
	// (group references in repl are in Perl/Java style: $0, $1, $2, ..)
	Replace(subject, repl string) string
	// ReplaceFunc replaces every occurrence of the pattern with
	// replacements produced by invoking the replacer on every match
	ReplaceFunc(subject string, replacer Replacer) string
}

// Replacer produces replacements for matches
type Replacer interface {
	Replacement(match Match) string
}

// Engine represents a regexp engine
type Engine interface {
	Compile(pattern string, flags int) (Regexp, error)
	Quote(s string) string

	FlagDotAll() int
	FlagExtended() int
	FlagUnicode() int
	FlagCaseInsensitive() int
	FlagMultiline() int
}

// ReEngine contains the regexp engine which should be used for all matching;
// For local dev runs init this with 'matcherpcre.NewEngine()'.
// For production uses init it from Java with its implementation.
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
