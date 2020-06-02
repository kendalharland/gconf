package main

import (
	"log"
	"net/http"
)

func newLoggingTransport(base http.RoundTripper) http.RoundTripper {
	return &modifyingTransport{
		base: base,
		modify: func(r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.String())
		},
	}
}

type modifyingTransport struct {
	base   http.RoundTripper
	modify func(*http.Request)
}

func (m *modifyingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	m.modify(r)
	return m.base.RoundTrip(r)
}
