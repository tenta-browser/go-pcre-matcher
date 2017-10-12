package replacer

import (
	"github.com/tenta-browser/go-pcre-matcher"
)

type replacerProto struct {
	replacementFunc func(match matcher.Match) string
}

func (rp *replacerProto) Replacement(match matcher.Match) string {
	return rp.replacementFunc(match)
}

// NewReplacer creates a simple struct that holds the provided function and
// implements the Replacer interface
func NewReplacer(rf func(match matcher.Match) string) matcher.Replacer {
	return &replacerProto{rf}
}
