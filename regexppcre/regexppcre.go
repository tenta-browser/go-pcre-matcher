package regexppcre

import (
	"goutils"
	gore "regexp"
	"strconv"
	"strings"

	"github.com/gijsbers/go-pcre"
)

type pcreEngine struct {
}

type pcreRegexp struct {
	wre pcre.Regexp
}

type pcreMatcher struct {
	wm      *pcre.Matcher
	subject string
}

// NewEngine creates a direct libpcre regexp engine instance
func NewEngine() goutils.Engine {
	return &pcreEngine{}
}

func (e *pcreEngine) Compile(pattern string, flags int) (goutils.Regexp, error) {
	wre, err := pcre.Compile(pattern, flags)
	if err != nil {
		return nil, err
	}
	return &pcreRegexp{wre}, nil
}

func (e *pcreEngine) Quote(s string) string {
	return gore.QuoteMeta(s)
}

func (e *pcreEngine) FlagDotAll() int {
	return pcre.DOTALL
}

func (re *pcreRegexp) Search(subject string) goutils.Matcher {
	wm := re.wre.MatcherString(subject, 0)
	if !wm.Matches() {
		return nil
	}
	return &pcreMatcher{wm, subject}
}

func (re *pcreRegexp) Replace(subject, repl string) string {
	//
	// goals: to work like the one in Java
	// non-goals: to be fast, since it's for testing only
	//
	wm := re.wre.MatcherString(subject, 0)
	parts := []string{}
	for wm.Matches() {
		// append substring up to match
		idxs := wm.Index()
		parts = append(parts, subject[:idxs[0]])

		// append replacement
		ln := len(repl)
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
				for ; c < ln && repl[c] >= '0' && repl[c] <= '9'; c++ {
					ngrp := grp + string(repl[c])
					if grpn, _ := strconv.Atoi(ngrp); grpn > wm.Groups() {
						break
					}
					grp = ngrp
				}
				if grp == "" {
					panic("Group index is missing")
				}
				grpn, _ := strconv.Atoi(grp)
				parts = append(parts, wm.GroupString(grpn))
			}
		}

		// find next match
		subject = subject[idxs[1]:]
		wm.MatchString(subject, 0)
	}
	parts = append(parts, subject)
	return strings.Join(parts, "")
}

func (m *pcreMatcher) Groups() int {
	return m.wm.Groups()
}

func (m *pcreMatcher) GroupByIdx(idx int) string {
	return m.wm.GroupString(idx)
}

func (m *pcreMatcher) GroupPresentByIdx(idx int) bool {
	return m.wm.Present(idx)
}

func (m *pcreMatcher) GroupByName(name string) string {
	group, err := m.wm.NamedString(name)
	if err != nil {
		return ""
	}
	return group
}

func (m *pcreMatcher) GroupPresentByName(name string) bool {
	present, err := m.wm.NamedPresent(name)
	return err == nil && present
}

func (m *pcreMatcher) Next() bool {
	matchRegion := m.wm.Index()
	if matchRegion == nil {
		// there's no current match, so there's no next one either
		return false
	}
	m.subject = m.subject[matchRegion[1]:] // trim matched part of subject
	return m.wm.MatchString(m.subject, 0)
}
