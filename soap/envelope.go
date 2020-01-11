package soap

import (
	"encoding/xml"
	"fmt"
)

type Envelope struct {
	Header *Header
	Body   interface{}
}

func (e *Envelope) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local != "Envelope" {
		return fmt.Errorf("Expected Envelope got (%s)", start.Name.Local)
	}

	b := &body{
		Contents: &e.Body,
	}

	err := d.Decode(b)
	if err != nil {
		return err
	}

	return d.Skip()
}

func (env *Envelope) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Local: "soapenv:Envelope",
		},
		Attr: []xml.Attr{
			xml.Attr{Name: xml.Name{Local: "xmlns:soapenv"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	if env.Header != nil {
		err = e.Encode(env.Header)
		if err != nil {
			return err
		}
	}

	b := &body{
		Contents: &env.Body,
	}

	err = e.Encode(b)
	if err != nil {
		return err
	}

	err = e.EncodeToken(t.End())
	if err != nil {
		return err
	}

	return e.Flush()
}

type body struct {
	Contents interface{}
}

func (b *body) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local != "Body" {
		return fmt.Errorf("Expected Body got (%s)", start.Name.Local)
	}

	if b.Contents != nil {
		err := d.DecodeElement(b.Contents, nil)
		if err != nil {
			return err
		}
	}

	return d.Skip()
}

func (b *body) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Local: "soapenv:Body",
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	if b.Contents != nil {
		err = e.Encode(b.Contents)
		if err != nil {
			return err
		}
	}

	err = e.EncodeToken(t.End())
	if err != nil {
		return err
	}

	return e.Flush()
}
