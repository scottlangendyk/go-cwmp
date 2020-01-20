package cwmp

import (
	"../soap"
	"encoding/xml"
	"time"
)

const XMLSpace = "urn:dslforum-org:cwmp-1-0"

var (
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

type ID string
type HoldRequests bool
type SessionTimeout uint
type SupportedCWMPVersions []string
type UseCWMPVersion string

type Header []interface{}

func (h *Header) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var hdr interface{}

	switch start.Name.Local {
	case "ID":
		hdr = new(ID)
	case "HoldRequests":
		hdr = new(HoldRequests)
	case "SessionTimeout":
		hdr = new(SessionTimeout)
	case "SupportedCWMPVersions":
		hdr = new(SupportedCWMPVersions)
	case "UseCWMPVersion":
		hdr = new(UseCWMPVersion)
	default:
		return d.Skip()
	}

	err := d.DecodeElement(hdr, &start)
	if err != nil {
		return err
	}

	*h = append(*h, hdr)

	return nil
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
	CommandKey string   `xml:"CommandKey"`
}

type RebootResponse struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 RebootResponse"`
}

type GetRPCMethods struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 GetRPCMethods"`
}

type GetRPCMethodsResponse struct {
	MethodList []string
}

func (msg *GetRPCMethodsResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.Skip()

	l := []string{}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		el, ok := t.(xml.StartElement)
		if !ok {
			break
		}

		var s string

		err = d.DecodeElement(&s, &el)
		if err != nil {
			return err
		}

		l = append(l, s)
	}

	msg.MethodList = l

	return d.Skip()
}

type Fault struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 Fault"`
	Code    uint     `xml:"FaultCode"`
	String  string   `xml:"FaultString"`
}

type TransferComplete struct {
	XMLName      xml.Name  `xml:"urn:dslforum-org:cwmp-1-0 TransferComplete"`
	CommandKey   string    `xml:"CommandKey"`
	Fault        Fault     `xml:"FaultStruct"`
	StartTime    time.Time `xml"StartTime"`
	CompleteTime time.Time `xml:"CompleteTime"`
}

type TransferCompleteResponse struct {
	XMLName xml.Name `xml:"urn:dslforum-org:cwmp-1-0 TransferCompleteResponse"`
}

type AutonomousTransferComplete struct {
	XMLName        xml.Name  `xml:"urn:dslforum-org:cwmp-1-0 AutonomousTransferComplete"`
	AnnounceURL    string    `xml:"AnnounceURL"`
	TransferURL    string    `xml:"TranserURL"`
	IsDownload     bool      `xml:"IsDownload"`
	FileType       string    `xml:"FileType"`
	FileSize       uint      `xml:"FileSize"`
	TargetFileName string    `xml:"TargetFileName"`
	Fault          Fault     `xml:"FaultStruct"`
	StartTime      time.Time `xml:"StartTime"`
	CompleteTime   time.Time `xml:"CompleteTime"`
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
