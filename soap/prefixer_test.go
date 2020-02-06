package soap

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrefixerAttrs(t *testing.T) {
	input := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/"><Header xmlns="http://schemas.xmlsoap.org/soap/envelope/"><ID xmlns="urn:dslforum-org:cwmp-1-0" xmlns:envelope="http://schemas.xmlsoap.org/soap/envelope/" envelope:mustUnderstand="1">1234</ID><SessionTimeout xmlns="urn:dslforum-org:cwmp-1-0">2</SessionTimeout><SupportedCWMPVersions xmlns="urn:dslforum-org:cwmp-1-0">1.0,1.1,1.4</SupportedCWMPVersions></Header><Body xmlns="http://schemas.xmlsoap.org/soap/envelope/"><InformResponse xmlns="urn:dslforum-org:cwmp-1-0"><MaxEnvelopes>1</MaxEnvelopes></InformResponse></Body></Envelope>`

	var b bytes.Buffer

	p := NewPrefixer(&b, map[string]string{
		"http://schemas.xmlsoap.org/soap/envelope/": "soapenv",
		"urn:dslforum-org:cwmp-1-0":                 "cwmp",
	})

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Error() != nil {
		t.Fatalf("err: %v", err)
	}

	want := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cwmp="urn:dslforum-org:cwmp-1-0"><soapenv:Header><cwmp:ID soapenv:mustUnderstand="1">1234</cwmp:ID><cwmp:SessionTimeout>2</cwmp:SessionTimeout><cwmp:SupportedCWMPVersions>1.0,1.1,1.4</cwmp:SupportedCWMPVersions></soapenv:Header><soapenv:Body><cwmp:InformResponse><MaxEnvelopes>1</MaxEnvelopes></cwmp:InformResponse></soapenv:Body></soapenv:Envelope>`
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}

func TestPrefixerMultipleNamespaces(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, map[string]string{
		"http://schemas.xmlsoap.org/soap/envelope/": "soapenv",
		"urn:dslforum-org:cwmp-1-0":                 "cwmp",
		"http://schemas.xmlsoap.org/soap/encoding/": "soap",
	})

	input := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/"><Body xmlns="http://schemas.xmlsoap.org/soap/envelope/"><InformResponse xmlns="urn:dslforum-org:cwmp-1-0"><MaxEnvelopes>1</MaxEnvelopes></InformResponse></Body></Envelope>`

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Error() != nil {
		t.Fatalf("err: %v", err)
	}

	want := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/"><soapenv:Body><cwmp:InformResponse><MaxEnvelopes>1</MaxEnvelopes></cwmp:InformResponse></soapenv:Body></soapenv:Envelope>`
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}

func TestPrefixer(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, map[string]string{
		"mynamespace": "yo",
	})

	input := `<test xmlns="mynamespace">Hey</test>`

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := `<yo:test xmlns:yo="mynamespace">Hey</yo:test>`
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}

func TestPrefixerNesting(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, map[string]string{
		"mynamespace": "yo",
	})

	input := `<test xmlns="mynamespace"><one><two>Hey</two><two>Man</two><two><three>3</three></two></one></test>`

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := `<yo:test xmlns:yo="mynamespace"><one><two>Hey</two><two>Man</two><two><three>3</three></two></one></yo:test>`
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}

func TestPrefixerRepeatedNesting(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, map[string]string{
		"mynamespace": "yo",
	})

	input := `<test xmlns="mynamespace"><one xmlns="mynamespace"><two>Hey</two><two>Man</two><two><three>3</three></two></one></test>`

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	want := `<yo:test xmlns:yo="mynamespace"><yo:one><two>Hey</two><two>Man</two><two><three>3</three></two></yo:one></yo:test>`
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}

func TestPrefixerPassthrough(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, nil)

	input := `<test xmlns="mynamespace">Hey</test>`

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	want := input
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}

func TestPrefixerMultiWrite(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, nil)

	input := `<multi xmlns="mynamespace">Hey</multi>`

	_, err := fmt.Fprint(p, input[:len(input)-11])
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = fmt.Fprint(p, input[len(input)-11:])
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.Error() != nil {
		t.Fatalf("err: %v", p.Error())
	}

	want := input
	got := b.String()

	if want != got {
		t.Fatalf("Doesn't match\nwant: %s\ngot:  %s", want, got)
	}
}
