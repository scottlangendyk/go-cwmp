package cwmp

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestEncodeInformResponse(t *testing.T) {
	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	r := &InformResponse{}

	err := e.Encode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<InformResponse><MaxEnvelopes>1</MaxEnvelopes></InformResponse>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}

func TestEncodeInformResponseMaxEnvelopes(t *testing.T) {
	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	r := &InformResponse{
		MaxEnvelopes: 99,
	}

	err := e.Encode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<InformResponse><MaxEnvelopes>1</MaxEnvelopes></InformResponse>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
