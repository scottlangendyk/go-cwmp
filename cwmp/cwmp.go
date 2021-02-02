package cwmp

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/scottlangendyk/go-cwmp/soap"
)

const XMLSpace = "urn:dslforum-org:cwmp-1-0"

const (
	ACSMethodNotSupported = 8000
	ACSRequestDenied      = 8001
	ACSInternalError      = 8002
	ACSInvalidArguments   = 8003
	ACSResourcesExceeded  = 8004
	ACSRetryRequest       = 8005
	ACSIncompatible       = 8006

	CPEMethodNotSupported          = 9000
	CPERequestDenied               = 9001
	CPEInternalError               = 9002
	CPEInvalidArguments            = 9003
	CPEResourcedExceeded           = 9004
	CPEInvalidParameterName        = 9005
	CPEInvalidParameterType        = 9006
	CPEInvalidParameterValue       = 9007
	CPEParameterNotWritable        = 9008
	CPENotificationRequestRejected = 9009
	CPEFileTransferFailure         = 9010
	CPEUploadFailure               = 9011
	CPEInvalidUUID                 = 9022
)

type CWMPVersions []string

func (v CWMPVersions) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(strings.Join(v, ","), start)
}

func (v *CWMPVersions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	*v = strings.Split(s, ",")

	return nil
}

type Header struct {
	ID                    *string
	HoldRequests          *bool
	SessionTimeout        *uint
	SupportedCWMPVersions *CWMPVersions
	UseCWMPVersion        *string
}

func (h Header) startElement(local string, mustUnderstand bool) xml.StartElement {
	start := xml.StartElement{
		Name: xml.Name{
			Local: local,
			Space: XMLSpace,
		},
	}

	if mustUnderstand {
		start.Attr = []xml.Attr{
			xml.Attr{
				Name: xml.Name{
					Space: soap.XMLSpaceEnvelope,
					Local: "mustUnderstand",
				},
				Value: "1",
			},
		}
	}

	return start
}

func (h Header) encodeHoldRequests(e *xml.Encoder) error {
	if h.HoldRequests == nil {
		return nil
	}

	if *h.HoldRequests == true {
		return e.EncodeElement("1", h.startElement("HoldRequests", true))
	}

	return e.EncodeElement("0", h.startElement("HoldRequests", true))
}

func (h Header) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Header"
	start.Name.Space = soap.XMLSpaceEnvelope

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	err = e.EncodeElement(h.ID, h.startElement("ID", true))
	if err != nil {
		return err
	}

	err = h.encodeHoldRequests(e)
	if err != nil {
		return err
	}

	err = e.EncodeElement(h.SessionTimeout, h.startElement("SessionTimeout", false))
	if err != nil {
		return err
	}

	err = e.EncodeElement(h.SupportedCWMPVersions, h.startElement("SupportedCWMPVersions", false))
	if err != nil {
		return err
	}

	err = e.EncodeElement(h.UseCWMPVersion, h.startElement("UseCWMPVersion", true))
	if err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

func (h *Header) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var hdr interface{}

	switch start.Name.Local {
	case "ID":
		h.ID = new(string)
		hdr = h.ID
	case "HoldRequests":
		h.HoldRequests = new(bool)
		hdr = h.HoldRequests
	case "SessionTimeout":
		h.SessionTimeout = new(uint)
		hdr = h.SessionTimeout
	case "SupportedCWMPVersions":
		h.SupportedCWMPVersions = new(CWMPVersions)
		hdr = h.SupportedCWMPVersions
	case "UseCWMPVersion":
		h.UseCWMPVersion = new(string)
		hdr = h.UseCWMPVersion
	default:
		return d.Skip()
	}

	return d.DecodeElement(hdr, &start)
}

func Decode(d *xml.Decoder) (*soap.Envelope, error) {
	b := &body{}
	h := &Header{}

	e := &soap.Envelope{
		Header: h,
		Body:   b,
	}

	err := d.Decode(e)
	if err != nil {
		return nil, err
	}

	e.Body = b.Contents

	return e, nil
}

type body struct {
	Contents interface{}
}

func (b *body) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(&b.Contents)
}

func (b *body) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	switch start.Name.Local {
	case "Fault":
		b.Contents = &soap.Fault{
			Detail: &Fault{},
		}
	case "Inform":
		b.Contents = &Inform{}
	case "InformRespone":
		b.Contents = &InformResponse{}
	case "GetRPCMethods":
		b.Contents = &GetRPCMethods{}
	case "GetRPCMethodsResponse":
		b.Contents = &GetRPCMethodsResponse{}
	case "Reboot":
		b.Contents = &Reboot{}
	case "RebootResponse":
		b.Contents = &RebootResponse{}
	case "TransferComplete":
		b.Contents = &TransferComplete{}
	case "TransferCompleteResponse":
		b.Contents = &TransferCompleteResponse{}
	case "AutonomousTransferComplete":
		b.Contents = &AutonomousTransferComplete{}
	case "AutonomousTransferCompleteResponse":
		b.Contents = &AutonomousTransferCompleteResponse{}
	case "Download":
		b.Contents = &Download{}
	case "DownloadResponse":
		b.Contents = &DownloadResponse{}
	case "GetParameterValues":
		b.Contents = &GetParameterValues{}
	case "GetParameterValuesResponse":
		b.Contents = &GetParameterValuesResponse{}
	case "SetParameterValues":
		b.Contents = &SetParameterValues{}
	case "SetParameterValuesResponse":
		b.Contents = &SetParameterValuesResponse{}
	case "GetParameterNames":
		b.Contents = &GetParameterNames{}
	case "GetParameterNamesResponse":
		b.Contents = &GetParameterNamesResponse{}
	default:
		return d.Skip()
	}

	return d.DecodeElement(&b.Contents, &start)
}

type Reboot struct {
	XMLName    xml.Name `xml:"urn:dslforum-org:cwmp-1-0 Reboot"`
	CommandKey string
}

type RebootResponse struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 RebootResponse"`
}

type GetRPCMethods struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetRPCMethods"`
}

type GetRPCMethodsResponse struct {
	XMLName    xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetRPCMethodsResponse"`
	MethodList []string `xml:"MethodList>string"`
}

type Fault struct {
	XMLName                 xml.Name `xml:"urn:dslforum-org:cwmp-1-0 Fault"`
	Code                    uint     `xml:"FaultCode"`
	String                  string   `xml:"FaultString"`
	SetParameterValuesFault []SetParameterValuesFault
}

type TransferComplete struct {
	XMLName      xml.Name `xml:"urn:dslforum-org:cwmp-1-0 TransferComplete"`
	CommandKey   string
	Fault        Fault `xml:"FaultStruct"`
	StartTime    time.Time
	CompleteTime time.Time
}

type TransferCompleteResponse struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 TransferCompleteResponse"`
}

type AutonomousTransferComplete struct {
	XMLName        xml.Name `xml:"urn:dslforum-org:cwmp-1-0 AutonomousTransferComplete"`
	AnnounceURL    string
	TransferURL    string
	IsDownload     bool
	FileType       string
	FileSize       uint
	TargetFileName string
	Fault          Fault `xml:"FaultStruct"`
	StartTime      time.Time
	CompleteTime   time.Time
}

type AutonomousTransferCompleteResponse struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 AutonomousTransferCompleteResponse"`
}

type Download struct {
	XMLName        xml.Name `xml:"urn:dslforum-org:cwmp-1-0 Download"`
	CommandKey     string
	FileType       string
	URL            string
	Username       string
	Password       string
	FileSize       uint
	TargetFileName string
	DelaySeconds   uint
	SuccessURL     string
	FailureURL     string
}

type DownloadResponse struct {
	XMLName      xml.Name `xml:"urn:dslforum-org:cwmp-1-0 DownloadResponse"`
	Completed    bool
	StartTime    time.Time
	CompleteTime time.Time
}

type GetParameterNames struct {
	XMLName       xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetParameterNames"`
	ParameterPath string
	NextLevel     bool
}

type ParameterInfo struct {
	Name     string
	Writable bool
}

type GetParameterNamesResponse struct {
	XMLName       xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetParameterNamesResponse"`
	ParameterList []ParameterInfo
}

type GetParameterValues struct {
	XMLName        xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetParameterValues"`
	ParameterNames string
}

type GetParameterValuesResponse struct {
	XMLName       xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetParameterValuesResponse"`
	ParameterList []ParameterValue
}

type ParameterValue struct {
	Name  string
	Value string
}

type SetParameterValues struct {
	XMLName       xml.Name `xml:"urn:dslforum-org:cwmp-1-0 SetParameterValues"`
	ParameterList []ParameterValue
	ParameterKey  string
}

type SetParameterValuesResponse struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 SetParameterValuesResponse"`
	Status  bool
}

type DeviceID struct {
	Manufacturer string
	OUI          string
	ProductClass string
	SerialNumber string
}

type Event struct {
	EventCode  string
	CommandKey string
}

type Inform struct {
	XMLName       xml.Name `xml:"urn:dslforum-org:cwmp-1-0 Inform"`
	RetryCount    uint
	CurrentTime   time.Time
	MaxEnvelopes  uint
	DeviceID      DeviceID         `xml:"DeviceId"`
	Event         []Event          `xml:"Event>EventStruct"`
	ParameterList []ParameterValue `xml:"ParameterList>ParameterValueStruct"`
}

type InformResponse struct {
	MaxEnvelopes uint
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

type SetParameterValuesFault struct {
	Code   uint   `xml:"FaultCode"`
	String string `xml:"FaultString"`
	Name   string `xml:"ParameterName"`
}
