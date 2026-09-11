package markdown_test

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/wikilayer/markdown"
)

type corpus struct {
	Cases []struct {
		Name     string           `yaml:"name"`
		Markdown string           `yaml:"markdown"`
		Found    []markdown.Found `yaml:"found"`
	} `yaml:"cases"`
}

func TestTheDialectRecognisesWhatTheCorpusSays(t *testing.T) {
	read, err := os.ReadFile("corpus/dialect.yaml")
	if err != nil {
		t.Fatalf("the corpus every port answers to is unreadable: %v", err)
	}

	var held corpus
	if err := yaml.Unmarshal(read, &held); err != nil {
		t.Fatalf("the corpus does not parse: %v", err)
	}
	if len(held.Cases) == 0 {
		t.Fatal("the corpus holds no cases, so this test cannot fail")
	}

	for _, one := range held.Cases {
		t.Run(one.Name, func(t *testing.T) {
			got := markdown.Recognise([]byte(one.Markdown))
			want := one.Found
			if want == nil {
				want = []markdown.Found{}
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("recognised %+v, the corpus says %+v", got, want)
			}
		})
	}
}
