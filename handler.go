package main

import (
	"net/http"
	"strings"
)

// interface for parser functions
type parserFunc func(*http.Request) (string, error)

type alertMessage struct {
	message    string
	recipients []string
}

type messageHandler struct {
	messages   chan<- alertMessage // chan to xmpp client
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
		var recipients []string = nil
		// check if the request contains per-request recipients
		value, overwrite := r.URL.Query()["recipients"]
		if overwrite {
			recipients = strings.Split(value[0], ",")
		}
		// send message to xmpp client
		h.messages <- alertMessage{
			message:    m,
			recipients: recipients,
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

// returns new handler with a given parser function
func newMessageHandler(m chan<- alertMessage, f parserFunc) *messageHandler {
	return &messageHandler{
		messages:   m,
		parserFunc: f,
	}
}
