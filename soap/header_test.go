package soap

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestEncodeHeader(t *testing.T) {
	h := Header{}

	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	err := e.Encode(&h)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<soapenv:Header></soapenv:Header>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
