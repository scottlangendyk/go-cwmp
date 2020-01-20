package cwmp

import (
	"encoding/xml"
	"time"
)

type DeviceID struct {
	Manufacturer string `xml:"Manufacturer"`
	OUI          string `xml:"OUI"`
	ProductClass string `xml:"ProductClass"`
	SerialNumber string `xml:"SerialNumber"`
}

type Event struct {
	EventCode  string
	CommandKey string
}

type Inform struct {
	RetryCount   uint      `xml:"RetryCount"`
	CurrentTime  time.Time `xml:"CurrentTime"`
	MaxEnvelopes uint      `xml:"MaxEnvelopes"`
	DeviceID     DeviceID  `xml:"DeviceId"`
}

type InformResponse struct {
	MaxEnvelopes uint `xml:"MaxEnvelopes"`
}

func (r InformResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := xml.StartElement{Name: xml.Name{Space: XMLSpace, Local: "InformResponse"}}
	err := e.EncodeToken(s)
	if err != nil {
		return err
	}

	err = e.EncodeElement(1, xml.StartElement{Name: xml.Name{Local: "MaxEnvelopes"}})
	if err != nil {
		return err
	}

	return e.EncodeToken(s.End())
}
