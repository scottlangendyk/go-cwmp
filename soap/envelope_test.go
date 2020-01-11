package soap

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestDecodeEnvelope(t *testing.T) {
	r := strings.NewReader(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Body><string>test</string></soapenv:Body></soapenv:Envelope>`)

	d := xml.NewDecoder(r)

	var s string

	e := Envelope{
		Body: &s,
	}

	err := d.Decode(&e)
	if err != nil {
		t.Errorf("%s", err)
	}

	if s != "test" {
		t.Errorf("Expected (test), got (%s)", s)
	}
}

func TestEncodeEnvelope(t *testing.T) {
	env := Envelope{
		Body: "test",
	}

	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	err := e.Encode(&env)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Body><string>test</string></soapenv:Body></soapenv:Envelope>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
