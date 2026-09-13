package wlmarkdown_test

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err, "the rules every port answers to are unreadable")

	var held writtenRules
	require.NoError(t, yaml.Unmarshal(read, &held), "the rules do not parse")
	require.NotEmpty(t, held.CalloutClassByMarker, "the rules are empty, so this test cannot fail")
	require.NotEmpty(t, held.MapMarker, "the rules are empty, so this test cannot fail")
	require.NotEmpty(t, held.RefSchemes, "the rules are empty, so this test cannot fail")
	return held
}

func TestEveryMarkerTheRulesNameCarriesItsClass(t *testing.T) {
	for marker, class := range rules(t).CalloutClassByMarker {
		t.Run(marker, func(t *testing.T) {
			source := fmt.Sprintf("%s\n> Body.\n", "> "+marker)
			found := wlmarkdown.New().Recognise([]byte(source))
			require.Len(t, found, 1, "%s opened nothing", marker)
			assert.Equal(t, "callout", found[0].Kind)
			assert.Equal(t, class, found[0].Class, "%s should carry the class the rules give it", marker)
		})
	}
}

func TestTheMarkersOnOfferAreEveryOneTheRulesName(t *testing.T) {
	written := rules(t)
	want := make([]string, 0, len(written.CalloutClassByMarker)+1)
	for marker := range written.CalloutClassByMarker {
		want = append(want, marker)
	}
	want = append(want, written.MapMarker)
	slices.Sort(want)

	assert.Equal(t, want, wlmarkdown.New().Markers(),
		"a host that wants to notice a marker nothing was made of has to keep its own list and drift")
}

func TestTheClassesOnOfferAreTheOnesTheRulesName(t *testing.T) {
	written := rules(t).CalloutClassByMarker
	want := make([]string, 0, len(written))
	for _, class := range written {
		want = append(want, class)
	}
	slices.Sort(want)
	want = slices.Compact(want)

	assert.Equal(t, want, wlmarkdown.New().Classes())
}

func TestTheSchemesOnOfferAreTheOnesTheRulesName(t *testing.T) {
	want := slices.Clone(rules(t).RefSchemes)
	slices.Sort(want)

	assert.Equal(t, want, wlmarkdown.New().Schemes())
}

func TestTheMarkerTheRulesNameOpensAPlace(t *testing.T) {
	source := fmt.Sprintf("> %s\n> 44.7866, 20.4489\n", rules(t).MapMarker)
	found := wlmarkdown.New().Recognise([]byte(source))
	require.Len(t, found, 1, "the map marker went unrecognised")
	assert.Equal(t, "map", found[0].Kind)
}

func TestEverySchemeTheRulesNameIsReadAsOne(t *testing.T) {
	for _, scheme := range rules(t).RefSchemes {
		t.Run(scheme, func(t *testing.T) {
			source := fmt.Sprintf("A [label](%s:1).\n", scheme)
			found := wlmarkdown.New().Recognise([]byte(source))
			require.Len(t, found, 1, "%s: went unrecognised", scheme)
			assert.Equal(t, scheme, found[0].Scheme)
			assert.Equal(t, scheme+":1", found[0].Destination)
		})
	}
}
