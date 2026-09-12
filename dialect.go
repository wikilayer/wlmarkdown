package wlmarkdown

import (
	_ "embed"
	"regexp"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
)

//go:embed corpus/rules.yaml
var writtenRules []byte

type Dialect struct {
	calloutClassByMarker map[string]string
	mapMarker            string
	coordinate           *regexp.Regexp
	refSchemes           []string
}

type rules struct {
	CalloutClassByMarker map[string]string `yaml:"callout_class_by_marker"`
	MapMarker            string            `yaml:"map_marker"`
	Coordinate           string            `yaml:"coordinate"`
	RefSchemes           []string          `yaml:"ref_schemes"`
}

var read = sync.OnceValue(func() Dialect {
	var held rules
	if err := yaml.Unmarshal(writtenRules, &held); err != nil {
		panic("corpus/rules.yaml does not parse: " + err.Error())
	}
	return Dialect{
		calloutClassByMarker: held.CalloutClassByMarker,
		mapMarker:            held.MapMarker,
		coordinate:           regexp.MustCompile(held.Coordinate),
		refSchemes:           held.RefSchemes,
	}
})

func New() Dialect {
	return read()
}

func (d Dialect) Extensions() []goldmark.Extender {
	return []goldmark.Extender{
		extension.GFM,
		&calloutExtension{classByMarker: d.calloutClassByMarker},
		&mapEmbedExtension{marker: d.mapMarker, coordinate: d.coordinate},
	}
}

func (d Dialect) Parser() parser.Parser {
	return goldmark.New(goldmark.WithExtensions(d.Extensions()...)).Parser()
}
