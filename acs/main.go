package main

import (
	"../cwmp"
	"../soap"
	"encoding/xml"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Println(r.Cookies())

	if r.Header.Get("Content-Length") == "0" {
		w.WriteHeader(204)
		return
	}

	d := xml.NewDecoder(r.Body)

	msg, err := cwmp.Decode(d)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/xml")
	w.Header().Set("SOAPAction", "")

	e := xml.NewEncoder(w)

	switch m := msg.Body.(type) {
	case *cwmp.Inform:
		fmt.Println(m)
		http.SetCookie(w, &http.Cookie{Name: "test", Value: "test"})
		e.Encode(&soap.Envelope{
			Body: cwmp.InformResponse{},
		})
	case *soap.Fault:
		fmt.Println(m.Detail)
	default:
		w.Header().Del("Content-Type")
		w.Header().Del("SOAPAction")
		w.WriteHeader(204)
	}
}

func main() {
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8081", nil)
}
