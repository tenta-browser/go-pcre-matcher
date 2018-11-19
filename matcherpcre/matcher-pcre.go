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
 * matcher-pcre.go: PCRE matcher implementation
 */

package matcherpcre

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/tenta-browser/go-pcre-matcher"

	"github.com/gijsbers/go-pcre"
)

type pcreEngine struct {
}

type pcreRegexp struct {
	wre pcre.Regexp
}

type pcreMatch struct {
	wm      *pcre.Matcher
	subject string
}

// NewEngine creates a direct libpcre regexp engine instance
func NewEngine() matcher.Engine {
	return &pcreEngine{}
}

func (e *pcreEngine) Compile(pattern string, flags int) (matcher.Regexp, error) {
	wre, err := pcre.Compile(pattern, flags)
	if err != nil {
		return nil, err
	}
	return &pcreRegexp{wre}, nil
}

func (e *pcreEngine) Quote(s string) string {
	return regexp.QuoteMeta(s)
}

func (e *pcreEngine) FlagDotAll() int {
	return pcre.DOTALL
}

func (e *pcreEngine) FlagExtended() int {
	return pcre.EXTENDED
}

func (e *pcreEngine) FlagUnicode() int {
	return pcre.UCP
}

func (e *pcreEngine) FlagCaseInsensitive() int {
	return pcre.CASELESS
}

func (e *pcreEngine) FlagMultiline() int {
	return pcre.MULTILINE
}

func (re *pcreRegexp) Search(subject string) matcher.Match {
	wm := re.wre.MatcherString(subject, 0)
	if !wm.Matches() {
		return nil
	}
	return &pcreMatch{wm, subject}
}

func (re *pcreRegexp) Replace(subject, repl string) string {
	//
	// goals: to work like the one in Java
	// non-goals: to be fast, since it's for testing only
	//
	ln := len(repl)
	return re.replaceFuncCommon(subject, func(m matcher.Match) string {
		parts := []string{}
		for c := 0; ; {
			b := c
			for ; c < ln && repl[c] != '\\' && repl[c] != '$'; c++ {
			}
			parts = append(parts, repl[b:c])
			if c == ln {
				break
			}
			if repl[c] == '\\' {
				c++
				if c < ln {
					parts = append(parts, string(repl[c]))
					c++
				} else {
					panic("Character to be escaped is missing")
				}
			} else if repl[c] == '$' {
				c++
				grp := ""
				part := ""
				if c < ln && repl[c] == '{' {
					c++
					for ; c < ln && repl[c] != '}'; c++ {
						grp += string(repl[c])
					}
					if c < ln {
						c++
					}
					if grp == "" {
						panic("Group name is missing")
					}
					part = m.GroupByName(grp)
				} else {
					for ; c < ln && repl[c] >= '0' && repl[c] <= '9'; c++ {
						ngrp := grp + string(repl[c])
						if grpn, _ := strconv.Atoi(ngrp); grpn > m.Groups() {
							break
						}
						grp = ngrp
					}
					if grp == "" {
						panic("Group index is missing")
					}
					grpn, _ := strconv.Atoi(grp)
					part = m.GroupByIdx(grpn)
				}
				parts = append(parts, part)
			}
		}
		return strings.Join(parts, "")
	})
}

func (re *pcreRegexp) ReplaceFunc(subject string, replacer matcher.Replacer) string {
	return re.replaceFuncCommon(subject, replacer.Replacement)
}

func (re *pcreRegexp) replaceFuncCommon(subject string, replacementFunc func(matcher.Match) string) string {
	wm := re.wre.MatcherString(subject, 0)
	parts := []string{}
	for wm.Matches() {
		// append substring up to match
		idxs := wm.Index()
		parts = append(parts, subject[:idxs[0]])

		// append replacement
		parts = append(parts, replacementFunc(&pcreMatch{wm, subject}))

		// find next match
		subject = subject[idxs[1]:]
		wm.MatchString(subject, 0)
	}
	parts = append(parts, subject)
	return strings.Join(parts, "")
}

func (m *pcreMatch) Groups() int {
	return m.wm.Groups()
}

func (m *pcreMatch) GroupByIdx(idx int) string {
	return m.wm.GroupString(idx)
}

func (m *pcreMatch) GroupPresentByIdx(idx int) bool {
	return m.wm.Present(idx)
}

func (m *pcreMatch) GroupByName(name string) string {
	group, err := m.wm.NamedString(name)
	if err != nil {
		return ""
	}
	return group
}

func (m *pcreMatch) GroupPresentByName(name string) bool {
	present, err := m.wm.NamedPresent(name)
	return err == nil && present
}

func (m *pcreMatch) Next() bool {
	matchRegion := m.wm.Index()
	if matchRegion == nil {
		// there's no current match, so there's no next one either
		return false
	}
	m.subject = m.subject[matchRegion[1]:] // trim matched part of subject
	return m.wm.MatchString(m.subject, 0)
}
