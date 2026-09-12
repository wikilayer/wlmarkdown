package wlmarkdown_test

import (
	"testing"

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
	return wlmarkdown.DeclinedIn(pc)
}

func TestAPlaceWhoseCoordinatesDoNotReadIsReportedDeclined(t *testing.T) {
	declined := declinedIn(t, "> [!MAP]\n> somewhere near the river\n")
	if len(declined) != 1 || declined[0].Marker != "[!MAP]" {
		t.Errorf("the author cannot see why no map appeared, and neither can the host: %+v", declined)
	}
}

func TestACalloutInsideACalloutIsReportedDeclined(t *testing.T) {
	declined := declinedIn(t, "> [!NOTE]\n> Outer.\n>\n> > [!TIP]\n> > Inner.\n")
	if len(declined) != 1 || declined[0].Marker != "[!TIP]" {
		t.Errorf("the inner quote stands unclaimed and should say so: %+v", declined)
	}
}

func TestWhatTheDialectMadeIsNotReportedDeclined(t *testing.T) {
	declined := declinedIn(t, "> [!NOTE]\n> Body.\n\n> [!MAP]\n> 44.7866, 20.4489\n")
	if len(declined) != 0 {
		t.Errorf("both were made, so there is nothing to report: %+v", declined)
	}
}

func TestAMarkerSharingItsLineIsNotReportedDeclined(t *testing.T) {
	declined := declinedIn(t, "> [!NOTE] see below\n> Body.\n")
	if len(declined) != 0 {
		t.Errorf("that quote was never a candidate, so calling it declined cries wolf: %+v", declined)
	}
}
