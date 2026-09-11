package wlmarkdown

import "testing"

func TestADialectRecognisesOnlyTheMarkersItCarries(t *testing.T) {
	own := Dialect{
		calloutClassByMarker: map[string]string{"[!ASIDE]": "aside"},
		mapMarker:            "[!PLACE]",
	}

	found := own.Recognise([]byte("> [!ASIDE]\n> Spoken aside.\n"))
	if len(found) != 1 || found[0].Class != "aside" {
		t.Errorf("the marker this dialect declares went unrecognised: %+v", found)
	}

	found = own.Recognise([]byte("> [!NOTE]\n> A marker of another dialect.\n"))
	if len(found) != 0 {
		t.Errorf("a marker this dialect never declared was recognised: %+v", found)
	}

	found = own.Recognise([]byte("> [!PLACE]\n> 44.8032201, 20.4726824\n"))
	if len(found) != 1 || found[0].Kind != "map" {
		t.Errorf("the place marker this dialect declares went unrecognised: %+v", found)
	}
}
