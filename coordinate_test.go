package wlmarkdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestASignIsACharacterOfTheAlphabet(t *testing.T) {
	written := coordinate{Signs: "+-−", Digits: "0123456789", Point: "."}

	for _, spoken := range []string{"+44.7866", "-44.7866", "−44.7866", "44.7866"} {
		assert.True(t, written.reads(spoken),
			"%q is a coordinate under these rules and was refused, which is how the Go port "+
				"drifts from a port that walks characters rather than bytes", spoken)
	}
}
