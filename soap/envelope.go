package soap

import (
	"encoding/xml"
	"fmt"
)

const (
	XMLSpaceEnvelope = "http://schemas.xmlsoap.org/soap/envelope/"
	XMLSpaceEncoding = "http://schemas.xmlsoap.org/soap/encoding/"
)

type Header interface {
	MustUnderstand() bool
	Name() xml.Name
}

type Envelope struct {
	Header interface{}
	Body   interface{}
}

func (e *Envelope) startElement(d *xml.Decoder) (*xml.StartElement, error) {
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}

		switch el := t.(type) {
		case xml.EndElement:
			return nil, fmt.Errorf("soap: Unexpected EndElement (%s)", el.Name.Local)
		case xml.StartElement:
			return &el, nil
		}
	}
}

func (e *Envelope) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local != "Envelope" {
		return fmt.Errorf("soap: Expected (Envelope) got (%s)", start.Name.Local)
	}

	el, err := e.startElement(d)
	if err != nil {
		return err
	}

	if el.Name.Local == "Header" {
		h := &element{
			Name: "Header",
			Contents: e.Header,
		}

		err = d.DecodeElement(h, el)
		if err != nil {
			return err
		}

		el, err = e.startElement(d)
		if err != nil {
			return err
		}
	}

	if el.Name.Local != "Body" {
		return fmt.Errorf("soap: Expected (Body) got (%s)", el.Name.Local)
	}

	b := &element{
		Name: "Body",
		Contents: &e.Body,
	}

	err = d.DecodeElement(b, el)
	if err != nil {
		return err
	}

	return d.Skip()
}

func (env Envelope) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Space: XMLSpaceEnvelope,
			Local: "Envelope",
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	if env.Header != nil {
		h := &element{
			Contents: env.Header,
			Name: "Header",
		}

		err = e.Encode(h)
		if err != nil {
			return err
		}
	}

	b := &element{
		Name: "Body",
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

type element struct {
	Name string
	Contents interface{}
}

func (e *element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if e.Contents == nil {
		return d.Skip()
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch el := t.(type) {
		case xml.EndElement:
			if el.Name.Local == e.Name {
				return nil
			}

			return fmt.Errorf("soap: Unexpected EndElement")
		case xml.StartElement:
			err = d.DecodeElement(e.Contents, &el)
			if err != nil {
				return err
			}
		}
	}
}

func (el element) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Space: XMLSpaceEnvelope,
			Local: el.Name,
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	if el.Contents != nil {
		err = e.Encode(el.Contents)
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