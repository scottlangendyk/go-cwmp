package soap

import (
	"encoding/xml"
	"fmt"
)

type faultDetail struct {
	Contents interface{}
}

func (f *faultDetail) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if f.Contents != nil {
		err := d.DecodeElement(f.Contents, nil)
		if err != nil {
			return err
		}
	}

	return d.Skip()
}

type Fault struct {
	Code   string
	String string
	Factor string
	Detail interface{}
}

func (f *Fault) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if start.Name.Local != "Fault" {
		return fmt.Errorf("Expected Fault got (%s)", start.Name.Local)
	}

	detail := &faultDetail{
		Contents: f.Detail,
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch el := t.(type) {
		case xml.EndElement:
			if el.Name.Local == "Fault" {
				return nil
			}

			return fmt.Errorf("Unexpected EndElement")
		case xml.StartElement:
			var v interface{}

			switch el.Name.Local {
			case "faultcode":
				v = &f.Code
			case "faultstring":
				v = &f.String
			case "faultfactor":
				v = &f.Factor
			case "detail":
				v = &detail
			}

			err = d.DecodeElement(v, &el)
			if err != nil {
				return err
			}
		}
	}

	f.Detail = detail.Contents

	return nil
}

func (f *Fault) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Local: "soapenv:Fault",
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	err = e.EncodeElement(&f.Code, xml.StartElement{Name: xml.Name{Local: "faultcode"}})
	if err != nil {
		return err
	}

	err = e.EncodeElement(&f.String, xml.StartElement{Name: xml.Name{Local: "faultstring"}})
	if err != nil {
		return err
	}

	err = e.EncodeElement(&f.Factor, xml.StartElement{Name: xml.Name{Local: "faultfactor"}})
	if err != nil {
		return err
	}

	d := xml.StartElement{Name: xml.Name{Local: "detail"}}

	err = e.EncodeToken(d)
	if err != nil {
		return err
	}

	err = e.Encode(&f.Detail)
	if err != nil {
		return err
	}

	err = e.EncodeToken(d.End())
	if err != nil {
		return err
	}

	err = e.EncodeToken(t.End())
	if err != nil {
		return err
	}

	return e.Flush()
}
