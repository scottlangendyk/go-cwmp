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
		h := &header{
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

	b := &body{
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
		h := &header{
			Contents: env.Header,
		}

		err = e.Encode(h)
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

type header struct {
	Contents interface{}
}

func (h *header) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if h.Contents == nil {
		return d.Skip()
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch el := t.(type) {
		case xml.EndElement:
			if el.Name.Local == "Header" {
				return nil
			}

			return fmt.Errorf("soap: Unexpected EndElement")
		case xml.StartElement:
			err = d.DecodeElement(h.Contents, &el)
			if err != nil {
				return err
			}
		}
	}
}

func (h header) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Space: XMLSpaceEnvelope,
			Local: "Header",
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	if h.Contents != nil {
		err = e.Encode(h.Contents)
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

type body struct {
	Contents interface{}
}

func (b *body) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local != "Body" {
		return fmt.Errorf("soap: Expected (Body) got (%s)", start.Name.Local)
	}

	if b.Contents == nil {
		return d.Skip()
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch el := t.(type) {
		case xml.EndElement:
			if el.Name.Local == "Body" {
				return nil
			}

			return fmt.Errorf("soap: Unexpected EndElement")
		case xml.StartElement:
			err = d.DecodeElement(b.Contents, &el)
			if err != nil {
				return err
			}
		}
	}
}

func (b body) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Space: XMLSpaceEnvelope,
			Local: "Body",
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
