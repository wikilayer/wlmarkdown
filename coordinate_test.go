package wlmarkdown

import "testing"

func TestASignIsACharacterOfTheAlphabet(t *testing.T) {
	written := coordinate{Signs: "+-−", Digits: "0123456789", Point: "."}

	for _, spoken := range []string{"+44.7866", "-44.7866", "−44.7866", "44.7866"} {
		if !written.reads(spoken) {
			t.Errorf("%q is a coordinate under these rules and was refused, which is how the Go port "+
				"drifts from a port that walks characters rather than bytes", spoken)
		}
	}
}
