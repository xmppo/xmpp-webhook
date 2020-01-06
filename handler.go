package main

import (
	"net/http"
)

// interface for parser functions
type parserFunc func(*http.Request) (string, error)

type messageHandler struct {
	messages   chan<- string // chan to xmpp client
	parserFunc parserFunc
}

// http request handler
func (h *messageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// parse/generate message from http request
	m, err := h.parserFunc(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		// send message to xmpp client
		h.messages <- m
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

// returns new handler with a given parser function
func newMessageHandler(m chan<- string, f parserFunc) *messageHandler {
	return &messageHandler{
		messages:   m,
		parserFunc: f,
	}
}
