package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"../cwmp"
	"../soap"
)

func handleMessage(r *http.Request) (*soap.Envelope, error) {
	defer r.Body.Close()

	if r.ContentLength == 0 {
		return nil, nil
	}

	d := xml.NewDecoder(r.Body)

	msg, err := cwmp.Decode(d)
	if err != nil {
		return nil, err
	}

	switch h := msg.Header.(type) {
	case *cwmp.Header:
		if h.ID != nil {
			fmt.Println(*h.ID)
		}
	}

	switch m := msg.Body.(type) {
	case *cwmp.Inform:
		fmt.Println(m)
		msg = &soap.Envelope{
			Body: &cwmp.InformResponse{},
		}
	case *cwmp.GetRPCMethods:
		msg = &soap.Envelope{
			Body: &cwmp.GetRPCMethodsResponse{
				MethodList: []string{
					"Inform",
					"GetRPCMethods",
					"TransferComplete",
				},
			},
		}
	case *cwmp.TransferComplete:
		msg = &soap.Envelope{
			Body: &cwmp.TransferCompleteResponse{},
		}
	default:
		msg = &soap.Envelope{
			Body: &soap.Fault{
				Code:   "Client",
				String: "CWMP fault",
				Detail: &cwmp.Fault{
					Code:   cwmp.ACSMethodNotSupported,
					String: "Method not supported",
				},
			},
		}
	}

	return msg, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	msg, err := handleMessage(r)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if msg == nil {
		w.WriteHeader(204)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.Header().Set("SOAPAction", "")

	p := soap.NewPrefixer(w, map[string]string{soap.XMLSpaceEnvelope: "soapenv", cwmp.XMLSpace: "cwmp"})

	e := xml.NewEncoder(p)

	err = e.Encode(msg)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	http.ListenAndServe("0.0.0.0:8081", nil)
}
