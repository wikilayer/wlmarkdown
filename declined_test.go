package wlmarkdown_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/wikilayer/wlmarkdown"
)

func declinedIn(t *testing.T, source string) []wlmarkdown.Declined {
	t.Helper()
	md := goldmark.New(goldmark.WithExtensions(wlmarkdown.New().Extensions()...))
	pc := parser.NewContext()
	md.Parser().Parse(text.NewReader([]byte(source)), parser.WithContext(pc))
	held := wlmarkdown.DeclinedIn(pc)
	if held == nil {
		return []wlmarkdown.Declined{}
	}
	return held
}

func TestASecondParseDoesNotInheritTheFirstsRefusals(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(wlmarkdown.New().Extensions()...))
	pc := parser.NewContext()

	md.Parser().Parse(text.NewReader([]byte("> [!MAP]\n> nowhere near a point\n")), parser.WithContext(pc))
	md.Parser().Parse(text.NewReader([]byte("Nothing here at all.\n")), parser.WithContext(pc))

	assert.Empty(t, wlmarkdown.DeclinedIn(pc),
		"the second document turned nothing down and is told otherwise")
}
