package wlmarkdown_test

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/wikilayer/wlmarkdown"
)

type writtenRules struct {
	CalloutClassByMarker map[string]string `yaml:"callout_class_by_marker"`
	MapMarker            string            `yaml:"map_marker"`
	RefSchemes           []string          `yaml:"ref_schemes"`
}

func rules(t *testing.T) writtenRules {
	t.Helper()
	read, err := os.ReadFile("corpus/rules.yaml")
	if err != nil {
		t.Fatalf("the rules every port answers to are unreadable: %v", err)
	}
	var held writtenRules
	if err := yaml.Unmarshal(read, &held); err != nil {
		t.Fatalf("the rules do not parse: %v", err)
	}
	if len(held.CalloutClassByMarker) == 0 || held.MapMarker == "" || len(held.RefSchemes) == 0 {
		t.Fatal("the rules are empty, so this test cannot fail")
	}
	return held
}

func TestEveryMarkerTheRulesNameCarriesItsClass(t *testing.T) {
	for marker, class := range rules(t).CalloutClassByMarker {
		t.Run(marker, func(t *testing.T) {
			source := fmt.Sprintf("%s\n> Body.\n", "> "+marker)
			found := wlmarkdown.New().Recognise([]byte(source))
			if len(found) != 1 || found[0].Kind != "callout" || found[0].Class != class {
				t.Errorf("%s should carry class %q, recognised %+v", marker, class, found)
			}
		})
	}
}

func TestTheClassesOnOfferAreTheOnesTheRulesName(t *testing.T) {
	written := rules(t).CalloutClassByMarker
	want := make([]string, 0, len(written))
	for _, class := range written {
		want = append(want, class)
	}
	slices.Sort(want)
	want = slices.Compact(want)

	if got := wlmarkdown.New().Classes(); !slices.Equal(got, want) {
		t.Errorf("offered %v, the rules name %v", got, want)
	}
}

func TestTheSchemesOnOfferAreTheOnesTheRulesName(t *testing.T) {
	want := slices.Clone(rules(t).RefSchemes)
	slices.Sort(want)

	if got := wlmarkdown.New().Schemes(); !slices.Equal(got, want) {
		t.Errorf("offered %v, the rules name %v", got, want)
	}
}

func TestTheMarkerTheRulesNameOpensAPlace(t *testing.T) {
	source := fmt.Sprintf("> %s\n> 44.7866, 20.4489\n", rules(t).MapMarker)
	found := wlmarkdown.New().Recognise([]byte(source))
	if len(found) != 1 || found[0].Kind != "map" {
		t.Errorf("the map marker went unrecognised: %+v", found)
	}
}

func TestEverySchemeTheRulesNameIsReadAsOne(t *testing.T) {
	for _, scheme := range rules(t).RefSchemes {
		t.Run(scheme, func(t *testing.T) {
			source := fmt.Sprintf("A [label](%s:1).\n", scheme)
			found := wlmarkdown.New().Recognise([]byte(source))
			if len(found) != 1 || found[0].Scheme != scheme || found[0].Destination != scheme+":1" {
				t.Errorf("%s: should be read as a scheme, recognised %+v", scheme, found)
			}
		})
	}
}
