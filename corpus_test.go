package wlmarkdown_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/wikilayer/wlmarkdown"
)

type corpus struct {
	Cases []struct {
		Name     string             `yaml:"name"`
		Markdown string             `yaml:"markdown"`
		Found    []wlmarkdown.Found `yaml:"found"`
	} `yaml:"cases"`
}

func TestTheDialectRecognisesWhatTheCorpusSays(t *testing.T) {
	read, err := os.ReadFile("corpus/dialect.yaml")
	if err != nil {
		t.Fatalf("the corpus every port answers to is unreadable: %v", err)
	}

	reader := yaml.NewDecoder(bytes.NewReader(read))
	reader.KnownFields(true)

	var held corpus
	if err := reader.Decode(&held); err != nil {
		t.Fatalf("the corpus does not parse, or asks for something no port would see: %v", err)
	}
	if len(held.Cases) == 0 {
		t.Fatal("the corpus holds no cases, so this test cannot fail")
	}

	for _, one := range held.Cases {
		t.Run(one.Name, func(t *testing.T) {
			if one.Markdown == "" {
				t.Fatal("the case carries no markdown, so it passes on an empty document")
			}
			got := wlmarkdown.New().Recognise([]byte(one.Markdown))
			want := one.Found
			if want == nil {
				want = []wlmarkdown.Found{}
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("recognised %+v, the corpus says %+v", got, want)
			}
		})
	}
}
