package xmlutil

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
)

type Prefixer interface {
	io.Writer
	Error() error
}

func NewPrefixer(w io.Writer, p map[string]string) Prefixer {
	return &prefixer{
		e: xml.NewEncoder(w),
		p: p,
	}
}

type stackItem struct {
	Start    *xml.StartElement
	Prefixes map[string]string
}

type prefixer struct {
	p     map[string]string
	stack []*stackItem
	b     bytes.Buffer
	e     *xml.Encoder
	err   error
}

func (p *prefixer) Error() error {
	return p.err
}

func (p *prefixer) prefixForNamespace(ns string) string {
	for i := len(p.stack) - 1; i >= 0; i-- {
		for space, prefix := range p.stack[i].Prefixes {
			if space == ns {
				return prefix
			}
		}
	}

	return ""
}

func (p *prefixer) push(start xml.StartElement) xml.StartElement {
	item := &stackItem{
		Start: &xml.StartElement{
			Name: start.Name,
		},
		Prefixes: make(map[string]string),
	}

	p.stack = append(p.stack, item)

	var attrs []xml.Attr

	if len(p.stack) == 1 {
		var spaces []string

		for space := range p.p {
			spaces = append(spaces, space)
		}

		sort.Strings(spaces)

		for _, space := range spaces {
			prefix := p.p[space]
			attr := xml.Attr{
				Name: xml.Name{
					Local: fmt.Sprintf("xmlns:%s", prefix),
				},
				Value: space,
			}

			attrs = append(attrs, attr)
			item.Prefixes[space] = prefix
		}
	}

	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" && attr.Name.Space == "" {
			pfx := p.prefixForNamespace(attr.Value)

			if pfx == "" {
				attrs = append(attrs, attr)
				continue
			}

			item.Start.Name.Space = ""
			item.Start.Name.Local = fmt.Sprintf("%s:%s", pfx, item.Start.Name.Local)

			continue
		}

		if attr.Name.Space == "xmlns" {
			pfx := p.prefixForNamespace(attr.Value)

			if pfx == "" {
				item.Prefixes[attr.Value] = attr.Name.Local
				attr.Name.Space = ""
				attr.Name.Local = fmt.Sprintf("xmlns:%s", attr.Name.Local)
				attrs = append(attrs, attr)
				continue
			}

			item.Prefixes[attr.Name.Local] = pfx

			continue
		}

		attrs = append(attrs, attr)
	}

	for _, attr := range attrs {
		pfx := p.prefixForNamespace(attr.Name.Space)

		if pfx != "" {
			attr.Name.Space = ""
			attr.Name.Local = fmt.Sprintf("%s:%s", pfx, attr.Name.Local)
		}

		item.Start.Attr = append(item.Start.Attr, attr)
	}

	if item.Start.Name.Space != "" {
		item.Start.Name.Local = fmt.Sprintf("%s:%s", item.Start.Name.Space, item.Start.Name.Local)
		item.Start.Name.Space = ""
	}

	return *item.Start
}

func (p *prefixer) pop() xml.EndElement {
	s := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return s.Start.End()
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
			t = p.push(el)
		case xml.EndElement:
			t = p.pop()
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
