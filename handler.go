package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// interface for parser functions (grafana, prometheus, ...)
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
	}
	// send message to xmpp client
	h.messages <- m
	w.WriteHeader(http.StatusNoContent)
}

// returns new handler with a given parser function
func newMessageHandler(m chan<- string, f parserFunc) *messageHandler {
	return &messageHandler{
		messages:   m,
		parserFunc: f,
	}
}

/*************
GRAFANA PARSER
*************/
func grafanaParserFunc(r *http.Request) (string, error) {
	// get alert data from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// grafana alert struct
	alert := &struct {
		Title   string `json:"title"`
		RuleURL string `json:"ruleUrl"`
		State   string `json:"state"`
		Message string `json:"message"`
	}{}

	// parse body into the alert struct
	err = json.Unmarshal(body, &alert)
	if err != nil {
		return "", err
	}

	// contruct alert message
	var message string
	switch alert.State {
	case "ok":
		message = ":) " + alert.Title
	default:
		message = ":( " + alert.Title + "\n" + alert.Message + "\n" + alert.RuleURL
	}

	return message, nil
}

/****************
PROMETHEUS PARSER
****************/
func prometheusParserFunc(r *http.Request) (string, error) {
	// get alert data from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// prometheus alert struct
	alert := &struct {
		Status      string `json:"status"`
		ExternalURL string `json:"externalURL"`
	}{}

	// parse body into the alert struct
	err = json.Unmarshal(body, &alert)
	if err != nil {
		return "", err
	}

	// contruct alert message
	var message string
	switch alert.Status {
	case "resolved":
		message = ":) " + alert.ExternalURL
	default:
		message = ":( " + alert.ExternalURL
	}

	return message, nil
}
