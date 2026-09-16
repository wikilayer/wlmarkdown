package wlmarkdown_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikilayer/wlmarkdown"
	"gopkg.in/yaml.v3"
)

func TestStripAnswersTheSharedPlainTextCorpus(t *testing.T) {
	written, err := os.ReadFile("corpus/plain_text.yaml")
	require.NoError(t, err)

	var corpus struct {
		Cases []struct {
			Name     string `yaml:"name"`
			Markdown string `yaml:"markdown"`
			Plain    string `yaml:"plain"`
		} `yaml:"cases"`
	}
	decoder := yaml.NewDecoder(bytes.NewReader(written))
	decoder.KnownFields(true)
	require.NoError(t, decoder.Decode(&corpus))
	require.NotEmpty(t, corpus.Cases)

	for _, one := range corpus.Cases {
		t.Run(one.Name, func(t *testing.T) {
			assert.Equal(t, one.Plain, wlmarkdown.Strip([]byte(one.Markdown)))
		})
	}
}
