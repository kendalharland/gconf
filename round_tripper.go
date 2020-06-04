package main

import (
	"net/http"
)

func newModifyingTransport(base http.RoundTripper, modify func(*http.Request)) http.RoundTripper {
	return &modifyingTransport{
		base:   base,
		modify: modify,
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
