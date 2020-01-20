package soap

import (
	"bytes"
	"fmt"
	"testing"
)

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
		t.Errorf(err.Error())
	}

	if p.Error() != nil {
		t.Errorf(p.Error().Error())
	}

	expected := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/encoding/"><soapenv:Body><cwmp:InformResponse><MaxEnvelopes>1</MaxEnvelopes></cwmp:InformResponse></soapenv:Body></soapenv:Envelope>`

	if b.String() != expected {
		t.Errorf(b.String())
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
		t.Errorf(err.Error())
	}

	if b.String() != `<yo:test xmlns:yo="mynamespace">Hey</yo:test>` {
		t.Errorf(b.String())
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
		t.Errorf(err.Error())
	}

	if b.String() != `<yo:test xmlns:yo="mynamespace"><one><two>Hey</two><two>Man</two><two><three>3</three></two></one></yo:test>` {
		t.Errorf(b.String())
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
		t.Errorf(err.Error())
	}

	if b.String() != `<yo:test xmlns:yo="mynamespace"><yo:one><two>Hey</two><two>Man</two><two><three>3</three></two></yo:one></yo:test>` {
		t.Errorf(b.String())
	}
}

func TestPrefixerPassthrough(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, nil)

	input := `<test xmlns="mynamespace">Hey</test>`

	_, err := fmt.Fprint(p, input)
	if err != nil {
		t.Errorf(err.Error())
	}

	if b.String() != input {
		t.Errorf(b.String())
	}
}

func TestPrefixerMultiWrite(t *testing.T) {
	var b bytes.Buffer

	p := NewPrefixer(&b, nil)

	input := `<multi xmlns="mynamespace">Hey</multi>`

	_, err := fmt.Fprint(p, input[:len(input)-11])
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = fmt.Fprint(p, input[len(input)-11:])
	if err != nil {
		t.Errorf(err.Error())
	}

	if p.Error() != nil {
		t.Errorf(p.Error().Error())
	}

	if b.String() != input {
		t.Errorf(b.String())
	}
}
