package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type Prefixer interface {
	Error() error
	Write(w []byte) (int, error)
}

func NewPrefixer(w io.Writer, p map[string]string) Prefixer {
	return &prefixer{
		e: xml.NewEncoder(w),
		p: p,
	}
}

type prefixer struct {
	p     map[string]string
	b     bytes.Buffer
	e     *xml.Encoder
	start []xml.StartElement
	err   error
}

func (p *prefixer) Error() error {
	return p.err
}

func (p *prefixer) prefixStartElement(start xml.StartElement) xml.StartElement {
	s := xml.StartElement{
		Name: start.Name,
	}

	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" && p.p[attr.Value] != "" {
			s.Name.Local = fmt.Sprintf("%s:%s", p.p[attr.Value], start.Name.Local)

		} else {
			s.Attr = append(s.Attr, attr)
		}
	}

	if len(p.start) > 0 {
		return s
	}

	for space, prefix := range p.p {
		attr := xml.Attr{
			Name: xml.Name{
				Local: fmt.Sprintf("xmlns:%s", prefix),
			},
			Value: space,
		}

		s.Attr = append(s.Attr, attr)
	}

	return s
}

func (p *prefixer) Write(w []byte) (int, error) {
	d := xml.NewDecoder(&p.b)

	n, err := p.b.Write(w)
	if err != nil {
		return n, err
	}

	for {
		t, err := d.RawToken()
		p.err = err
		if err == io.EOF {
			p.err = nil
		}
		if err != nil {
			return n, nil
		}

		switch el := t.(type) {
		case xml.StartElement:
			start := p.prefixStartElement(el)
			p.start = append(p.start, start)
			t = start
		case xml.EndElement:
			end := p.start[len(p.start)-1]
			p.start = p.start[:len(p.start)-1]
			t = end.End()
		}

		err = p.e.EncodeToken(t)
		if err != nil {
			return n, err
		}

		err = p.e.Flush()
		if err != nil {
			return n, err
		}
	}
}
