package soap

import (
	"encoding/xml"
)

type Header struct {
}

func (h *Header) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	t := xml.StartElement{
		Name: xml.Name{
			Local: "soapenv:Header",
		},
	}

	err := e.EncodeToken(t)
	if err != nil {
		return err
	}

	err = e.EncodeToken(t.End())
	if err != nil {
		return err
	}

	return e.Flush()
}
