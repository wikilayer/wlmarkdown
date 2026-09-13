package wlmarkdown_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/wikilayer/wlmarkdown"
)

type corpus struct {
	Cases []struct {
		Name     string             `yaml:"name"`
		Markdown string             `yaml:"markdown"`
		Found    []wlmarkdown.Found `yaml:"found"`
		Declined []string           `yaml:"declined"`
	} `yaml:"cases"`
}

func TestTheDialectRecognisesWhatTheCorpusSays(t *testing.T) {
	read, err := os.ReadFile("corpus/dialect.yaml")
	require.NoError(t, err, "the corpus every port answers to is unreadable")

	reader := yaml.NewDecoder(bytes.NewReader(read))
	reader.KnownFields(true)

	var held corpus
	require.NoError(t, reader.Decode(&held),
		"the corpus does not parse, or asks for something no port would see")
	require.NotEmpty(t, held.Cases, "the corpus holds no cases, so this test cannot fail")

	for _, one := range held.Cases {
		t.Run(one.Name, func(t *testing.T) {
			require.NotEmpty(t, one.Markdown,
				"the case carries no markdown, so it passes on an empty document")

			want := one.Found
			if want == nil {
				want = []wlmarkdown.Found{}
			}
			assert.Equal(t, want, wlmarkdown.New().Recognise([]byte(one.Markdown)),
				"the dialect recognised something other than what the corpus names")

			turnedDown := make([]wlmarkdown.Declined, 0, len(one.Declined))
			for _, marker := range one.Declined {
				turnedDown = append(turnedDown, wlmarkdown.Declined{Marker: marker})
			}
			assert.Equal(t, turnedDown, declinedIn(t, one.Markdown),
				"the dialect turned down something other than what the corpus names")
		})
	}
}
