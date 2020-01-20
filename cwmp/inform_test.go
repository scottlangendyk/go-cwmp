package cwmp

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestEncodeRebootResponse(t *testing.T) {
	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	r := &RebootResponse{}

	err := e.Encode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<RebootResponse xmlns="urn:dslforum-org:cwmp-1-0"></RebootResponse>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
