package cwmp

import (
	"../soap"
	"bytes"
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func assertEncode(t *testing.T, v interface{}, want string) {
	var b bytes.Buffer

	e := xml.NewEncoder(&b)

	err := e.Encode(v)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	assertEqual(t, want, b.String())
}

func assertEqual(t *testing.T, want, got interface{}) {
	if reflect.DeepEqual(want, got) {
		return
	}

	t.Fatalf("Not equal\nwant: %v\ngot:  %v", want, got)
}

func assertInform(t *testing.T, want, got Inform) {
	assertEqual(t, want.CurrentTime.String(), got.CurrentTime.String())
	assertEqual(t, want.RetryCount, got.RetryCount)
	assertEqual(t, want.MaxEnvelopes, got.MaxEnvelopes)

	assertEqual(t, want.DeviceID.Manufacturer, got.DeviceID.Manufacturer)
	assertEqual(t, want.DeviceID.OUI, got.DeviceID.OUI)
	assertEqual(t, want.DeviceID.ProductClass, got.DeviceID.ProductClass)
	assertEqual(t, want.DeviceID.SerialNumber, got.DeviceID.SerialNumber)

	if len(want.Event) != len(got.Event) {
		t.Fatalf("Event lengths aren't equal\nwant: %d\ngot:  %d", len(want.Event), len(got.Event))
	}

	for i, e := range want.Event {
		assertEqual(t, e.CommandKey, got.Event[i].CommandKey)
		assertEqual(t, e.EventCode, got.Event[i].EventCode)
	}

	if len(want.ParameterList) != len(got.ParameterList) {
		t.Fatalf("ParameterList lengths aren't equal\nwant: %d\ngot:  %d", len(want.ParameterList), len(got.ParameterList))
	}

	for i, p := range want.ParameterList {
		assertEqual(t, p.Name, got.ParameterList[i].Name)
		assertEqual(t, p.Value, got.ParameterList[i].Value)
	}
}

func testDecodeInform(t *testing.T, filename string, want Inform) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	d := xml.NewDecoder(f)

	e, err := Decode(d)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	got, ok := e.Body.(*Inform)
	if !ok {
		t.Fatal("Body is not type Inform")
	}

	assertInform(t, want, *got)
}

func assertFault(t *testing.T, want, got Fault) {
	assertEqual(t, want.Code, got.Code)
	assertEqual(t, want.String, got.String)

	if len(want.SetParameterValuesFault) != len(got.SetParameterValuesFault) {
		t.Fatalf("SetParameterValuesFault lengths aren't equal\nwant: %d\ngot:  %d", len(want.SetParameterValuesFault), len(got.SetParameterValuesFault))
	}

	for i, p := range want.SetParameterValuesFault {
		assertEqual(t, p.Name, got.SetParameterValuesFault[i].Name)
		assertEqual(t, p.String, got.SetParameterValuesFault[i].String)
		assertEqual(t, p.Code, got.SetParameterValuesFault[i].Code)
	}
}

func assertDetail(t *testing.T, want, got interface{}) {
	w, ok := want.(*Fault)
	if !ok {
		t.Fatal("Detail is not type cwmp.Fault")
	}

	g, ok := got.(*Fault)
	if !ok {
		t.Fatal("Detail is not type cwmp.Fault")
	}

	assertFault(t, *w, *g)
}

func testDecodeFault(t *testing.T, filename string, want soap.Fault) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	d := xml.NewDecoder(f)

	e, err := Decode(d)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, ok := e.Body.(*soap.Fault)
	if !ok {
		t.Fatal("Body is not type soap.Fault")
	}

	got := *body

	assertEqual(t, want.Code, got.Code)
	assertEqual(t, want.String, got.String)
	assertEqual(t, want.Factor, got.Factor)
	assertDetail(t, want.Detail, got.Detail)
}

func TestDecodeFault(t *testing.T) {
	testDecodeFault(t, "testdata/fault.xml", soap.Fault{
		Code:   "Client",
		String: "CWMP fault",
		Detail: &Fault{
			Code:   9000,
			String: "Upload method not supported",
		},
	})

	testDecodeFault(t, "testdata/fault.2.xml", soap.Fault{
		Code:   "Client",
		String: "CWMP fault",
		Detail: &Fault{
			Code:   9003,
			String: "Invalid arguments",
			SetParameterValuesFault: []SetParameterValuesFault{
				SetParameterValuesFault{
					Name: "Device.Time.NTPServer1",
					Code: 9007,
					String: "Invalid IP Address",
				},
				SetParameterValuesFault{
					Name: "Device.Time.LocalTimeZoneName",
					Code: 9007,
					String: "String too long",
				},
			},
		},
	})
}

func TestDecodeInform(t *testing.T) {
	testDecodeInform(t, "testdata/inform.xml", Inform{
		DeviceID: DeviceID{
			Manufacturer: "MikroTik",
			OUI:          "E48D8C",
			ProductClass: "hAP mini",
			SerialNumber: "B7B20A1DE3F0",
		},
		Event: []Event{
			Event{EventCode: "2 PERIODIC"},
		},
		MaxEnvelopes: 1,
		CurrentTime:  time.Date(2020, 01, 02, 20, 50, 49, 0, time.FixedZone("EST", -5*60*60)),
		RetryCount:   0,
		ParameterList: []ParameterValue{
			ParameterValue{Name: "Device.RootDataModelVersion", Value: "2.11"},
			ParameterValue{Name: "Device.DeviceInfo.SoftwareVersion", Value: "6.46.1"},
			ParameterValue{Name: "Device.DeviceInfo.ProvisioningCode", Value: ""},
			ParameterValue{Name: "Device.DeviceInfo.HardwareVersion", Value: "v1.0"},
			ParameterValue{Name: "Device.ManagementServer.ParameterKey", Value: ""},
			ParameterValue{Name: "Device.ManagementServer.ConnectionRequestURL", Value: "http://10.31.0.130:7547/e55b182787ec5d4c2d4bcaa39d234920bd35"},
			ParameterValue{Name: "Device.ManagementServer.AliasBasedAddressing", Value: "0"},
		},
	})
}

func TestDecodeHeader(t *testing.T) {
	input := `<soapenv:Envelope>
	<soapenv:Header>
		<cwmp:ID>MyID123Here</cwmp:ID>
		<skip>Skip this header</skip>
		<cwmp:SessionTimeout>2</cwmp:SessionTimeout>
		<cwmp:UseCWMPVersion>1.4</cwmp:UseCWMPVersion>
		<cwmp:HoldRequests>1</cwmp:HoldRequests>
		<cwmp:SupportedCWMPVersions>1.0,1.1,1.4</cwmp:SupportedCWMPVersions>
	</soapenv:Header>
	<soapenv:Body>
	</soapenv:Body>
</soapenv:Envelope>`

	d := xml.NewDecoder(strings.NewReader(input))

	e, err := Decode(d)
	if err != nil {
		t.Errorf(err.Error())
	}

	h, ok := e.Header.(*Header)
	if !ok {
		t.Fatal("Not a header")
	}

	want := Header{
		ID: new(string),
		SessionTimeout: new(uint),
		HoldRequests: new(HoldRequests),
		UseCWMPVersion: new(string),
		SupportedCWMPVersions: new(CWMPVersions),
	}

	*want.ID = "MyID123Here"
	*want.SessionTimeout = 2
	*want.HoldRequests = true
	*want.UseCWMPVersion = "1.4"
	*want.SupportedCWMPVersions = CWMPVersions{"1.0","1.1","1.4"}

	assertHeader(t, want, *h)
}

func TestDecodeHeaderEmpty(t *testing.T) {
	input := `<soapenv:Envelope>
	<soapenv:Header>
	</soapenv:Header>
	<soapenv:Body>
	</soapenv:Body>
</soapenv:Envelope>`

	d := xml.NewDecoder(strings.NewReader(input))

	e, err := Decode(d)
	if err != nil {
		t.Errorf(err.Error())
	}

	h, ok := e.Header.(*Header)
	if !ok {
		t.Fatal("Not a header")
	}

	assertHeader(t, Header{}, *h)
}

func TestDecodeHeaderPartial(t *testing.T) {
	input := `<soapenv:Envelope>
	<soapenv:Header>
		<cwmp:SessionTimeout>2</cwmp:SessionTimeout>
		<cwmp:SupportedCWMPVersions>1.0,1.1,1.4</cwmp:SupportedCWMPVersions>
	</soapenv:Header>
	<soapenv:Body>
	</soapenv:Body>
</soapenv:Envelope>`

	d := xml.NewDecoder(strings.NewReader(input))

	e, err := Decode(d)
	if err != nil {
		t.Errorf(err.Error())
	}

	h, ok := e.Header.(*Header)
	if !ok {
		t.Fatal("Not a header")
	}

	want := Header{
		SessionTimeout: new(uint),
		SupportedCWMPVersions: new(CWMPVersions),
	}

	*want.SessionTimeout = 2
	*want.SupportedCWMPVersions = CWMPVersions{"1.0","1.1","1.4"}

	assertHeader(t, want, *h)
}

func assertHeader(t *testing.T, want, got Header) {
	assertEqual(t, want.ID, got.ID)
	assertEqual(t, want.SessionTimeout, got.SessionTimeout)
	assertEqual(t, want.HoldRequests, got.HoldRequests)
	assertEqual(t, want.UseCWMPVersion, got.UseCWMPVersion)

	if want.SupportedCWMPVersions == nil && got.SupportedCWMPVersions == nil {
		return
	}

	if want.SupportedCWMPVersions == nil && got.SupportedCWMPVersions != nil {
		t.Fatal("Expected nil SupportedCWMPVersions")
	}

	if len(*want.SupportedCWMPVersions) != len(*got.SupportedCWMPVersions) {
		t.Fatalf("SupportedCWMPVersions lengths aren't equal\nwant: %d\ngot:  %d", len(*want.SupportedCWMPVersions), len(*got.SupportedCWMPVersions))
	}
}

func TestEncodeInformResponse(t *testing.T) {
	assertEncode(t, &InformResponse{}, `<InformResponse xmlns="urn:dslforum-org:cwmp-1-0"><MaxEnvelopes>1</MaxEnvelopes></InformResponse>`)
}

func TestEncodeRebootResponse(t *testing.T) {
	assertEncode(t, &RebootResponse{}, `<RebootResponse xmlns="urn:dslforum-org:cwmp-1-0"></RebootResponse>`)
}

func TestEncodeGetRPCMethodsResponse(t *testing.T) {
	v := &GetRPCMethodsResponse{
		MethodList: []string{"Method1", "Method2"},
	}

	want := `<GetRPCMethodsResponse><MethodList><string>Method1</string><string>Method2</string></MethodList></GetRPCMethodsResponse>`

	assertEncode(t, v, want)
}

func TestEncodeHeader(t *testing.T) {
	v := &Header{
		ID: new(string),
		HoldRequests: new(HoldRequests),
		SessionTimeout: new(uint),
		UseCWMPVersion: new(string),
		SupportedCWMPVersions: new(CWMPVersions),
	}

	*v.ID = "1234"
	*v.HoldRequests = true
	*v.SessionTimeout = 2
	*v.UseCWMPVersion = "1.4"
	*v.SupportedCWMPVersions = CWMPVersions{"1.0","1.1","1.4"}

	want := `<Header xmlns="urn:dslforum-org:cwmp-1-0"><ID>1234</ID><HoldRequests>1</HoldRequests><SessionTimeout>2</SessionTimeout><SupportedCWMPVersions>1.0,1.1,1.4</SupportedCWMPVersions><UseCWMPVersion>1.4</UseCWMPVersion></Header>`

	assertEncode(t, v, want)
}

func TestEncodeHeaderEmpty(t *testing.T) {
	want := `<Header xmlns="urn:dslforum-org:cwmp-1-0"></Header>`

	assertEncode(t, &Header{}, want)
}

func TestEncodeHeaderPartial(t *testing.T) {
	v := &Header{
		SessionTimeout: new(uint),
		SupportedCWMPVersions: new(CWMPVersions),
	}

	*v.SessionTimeout = 2
	*v.SupportedCWMPVersions = CWMPVersions{"1.0","1.1","1.4"}

	want := `<Header xmlns="urn:dslforum-org:cwmp-1-0"><SessionTimeout>2</SessionTimeout><SupportedCWMPVersions>1.0,1.1,1.4</SupportedCWMPVersions></Header>`

	assertEncode(t, v, want)
}