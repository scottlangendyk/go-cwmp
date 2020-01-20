package cwmp

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	input := `<soapenv:Envelope><soapenv:Header><cwmp:ID>MyID123Here</cwmp:ID><cwmp:SessionTimeout>Two</cwmp:SessionTimeout></soapenv:Header><soapenv:Body></soapenv:Body></soapenv:Envelope>`

	d := xml.NewDecoder(strings.NewReader(input))

	e, err := Decode(d)
	if err != nil {
		t.Errorf(err.Error())
	}

	h, ok := e.Header.(*Header)
	if !ok {
		t.Errorf("Not a header")
	}

	if len(*h) < 1 {
		t.Errorf("Invalid header")
	}

	for _, hdr := range *h {
		switch el := hdr.(type) {
		case ID:
			if el != "MyID123Here" {
				t.Errorf("Wrong ID")
			}

			return
		}
	}

	t.Errorf("Missing ID header")
}

func TestEncodeInformResponse(t *testing.T) {
	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	r := &InformResponse{}

	err := e.Encode(r)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<InformResponse xmlns="urn:dslforum-org:cwmp-1-0"><MaxEnvelopes>1</MaxEnvelopes></InformResponse>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
