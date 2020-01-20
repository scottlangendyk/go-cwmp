package soap

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestDecodeFaultEmptyDetail(t *testing.T) {
	d := xml.NewDecoder(strings.NewReader(`<soapenv:Fault><faultcode>faultcodehere</faultcode><faultstring>faultstringhere</faultstring><faultfactor>faultfactorhere</faultfactor><detail></detail></soapenv:Fault>`))

	f := &Fault{}

	err := d.Decode(&f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if f.Code != "faultcodehere" {
		t.Errorf("Expected (faultcodehere) got (%s)", f.Code)
	}

	if f.String != "faultstringhere" {
		t.Errorf("Expected (faultstringhere) got (%s)", f.String)
	}

	if f.Factor != "faultfactorhere" {
		t.Errorf("Expected (faultfactorhere) got (%s)", f.Factor)
	}
}

func TestDecodeFaultTabs(t *testing.T) {
	d := xml.NewDecoder(strings.NewReader("<soapenv:Envelope>\n\t<soapenv:Body>\n\t\t<soapenv:Fault>\n\t\t\t<faultcode>faultcodehere</faultcode>\n\t\t\t<faultstring>faultstringhere</faultstring>\n\t\t\t<faultfactor>faultfactorhere</faultfactor>\n\t\t\t<detail>\n\t\t\t\t<string>detailhere</string>\n\t\t\t</detail>\n\t\t</soapenv:Fault>\n\t</soapenv:Body>\n</soapenv:Envelope>"))

	var detail string

	f := &Fault{
		Detail: &detail,
	}

	e := &Envelope{
		Body: f,
	}

	err := d.Decode(e)
	if err != nil {
		t.Errorf("%s", err)
	}

	if f.Code != "faultcodehere" {
		t.Errorf("Expected (faultcodehere) got (%s)", f.Code)
	}

	if f.String != "faultstringhere" {
		t.Errorf("Expected (faultstringhere) got (%s)", f.String)
	}

	if f.Factor != "faultfactorhere" {
		t.Errorf("Expected (faultfactorhere) got (%s)", f.Factor)
	}

	if detail != "detailhere" {
		t.Errorf("Expected (detailhere) got (%s)", detail)
	}
}

func TestDecodeFault(t *testing.T) {
	d := xml.NewDecoder(strings.NewReader(`<soapenv:Fault><faultcode>faultcodehere</faultcode><faultstring>faultstringhere</faultstring><faultfactor>faultfactorhere</faultfactor><detail><string>detailhere</string></detail></soapenv:Fault>`))

	var detail string

	f := &Fault{
		Detail: &detail,
	}

	err := d.Decode(&f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if f.Code != "faultcodehere" {
		t.Errorf("Expected (faultcodehere) got (%s)", f.Code)
	}

	if f.String != "faultstringhere" {
		t.Errorf("Expected (faultstringhere) got (%s)", f.String)
	}

	if f.Factor != "faultfactorhere" {
		t.Errorf("Expected (faultfactorhere) got (%s)", f.Factor)
	}

	if detail != "detailhere" {
		t.Errorf("Expected (detailhere) got (%s)", detail)
	}
}

func TestEncodeFault(t *testing.T) {
	fault := Fault{
		Code:   "faultcodehere",
		String: "faultstringhere",
		Factor: "faultfactorhere",
		Detail: "detailhere",
	}

	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	err := e.Encode(&fault)
	if err != nil {
		t.Errorf("%s", err)
	}

	expected := `<Fault xmlns="http://schemas.xmlsoap.org/soap/envelope/"><faultcode>faultcodehere</faultcode><faultstring>faultstringhere</faultstring><faultfactor>faultfactorhere</faultfactor><detail><string>detailhere</string></detail></Fault>`

	if b.String() != expected {
		t.Errorf("Got (%s) Expected (%s)", b.String(), expected)
	}
}
