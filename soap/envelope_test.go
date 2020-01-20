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

func TestDecodeEnvelopeWithHeader(t *testing.T) {
	r := strings.NewReader(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"><soapenv:Header><string>One</string><string>Two</string></soapenv:Header><soapenv:Body><string>test</string></soapenv:Body></soapenv:Envelope>`)

	d := xml.NewDecoder(r)

	var s string
	var h []string

	e := Envelope{
		Body:   &s,
		Header: &h,
	}

	err := d.Decode(&e)
	if err != nil {
		t.Errorf("%s", err)
	}

	if s != "test" {
		t.Errorf("Expected (test), got (%s)", s)
	}

	if len(h) != 2 {
		t.Errorf("Expected 2 header items, got (%d)", len(h))
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

	expected := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/"><Body xmlns="http://schemas.xmlsoap.org/soap/envelope/"><string>test</string></Body></Envelope>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}

func TestEncodeEnvelopeWithHeader(t *testing.T) {
	env := Envelope{
		Body:   "test",
		Header: []string{"header"},
	}

	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	err := e.Encode(&env)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/"><Header xmlns="http://schemas.xmlsoap.org/soap/envelope/"><string>header</string></Header><Body xmlns="http://schemas.xmlsoap.org/soap/envelope/"><string>test</string></Body></Envelope>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
