package wlmarkdown

import (
	_ "embed"
	"slices"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
)

//go:embed corpus/rules.yaml
var writtenRules []byte

type Dialect struct {
	calloutClassByMarker map[string]string
	mapMarker            string
	coordinate           coordinate
	blanks               string
	refSchemes           []string
}

type coordinate struct {
	Signs           string `yaml:"signs"`
	Digits          string `yaml:"digits"`
	Point           string `yaml:"point"`
	LatitudeWithin  string `yaml:"latitude_within"`
	LongitudeWithin string `yaml:"longitude_within"`
}

type rules struct {
	CalloutClassByMarker map[string]string `yaml:"callout_class_by_marker"`
	MapMarker            string            `yaml:"map_marker"`
	Coordinate           coordinate        `yaml:"coordinate"`
	Blanks               string            `yaml:"blanks"`
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
		coordinate:           held.Coordinate,
		blanks:               held.Blanks,
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
		&mapEmbedExtension{marker: d.mapMarker, point: d.coordinate},
		&declinedExtension{markers: d.Markers()},
	}
}

func (d Dialect) Markers() []string {
	markers := make([]string, 0, len(d.calloutClassByMarker)+1)
	for marker := range d.calloutClassByMarker {
		markers = append(markers, marker)
	}
	markers = append(markers, d.mapMarker)
	slices.Sort(markers)
	return slices.Compact(markers)
}

func (d Dialect) Classes() []string {
	classes := make([]string, 0, len(d.calloutClassByMarker))
	for _, class := range d.calloutClassByMarker {
		classes = append(classes, class)
	}
	slices.Sort(classes)
	return slices.Compact(classes)
}

func (d Dialect) Schemes() []string {
	schemes := slices.Clone(d.refSchemes)
	slices.Sort(schemes)
	return schemes
}

func Kinds() []ast.NodeKind {
	return []ast.NodeKind{KindCallout, KindMapEmbed, KindUnreadable}
}

func (d Dialect) Parser() parser.Parser {
	return goldmark.New(goldmark.WithExtensions(d.Extensions()...)).Parser()
}
