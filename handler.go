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

/*************
SLACK PARSER
*************/
type SlackMessage struct {
   Channel     string            `json:"channel"`
   IconEmoji   string            `json:"icon_emoji"`
   Username    string            `json:"username"`
   Text        string            `json:"text"`
   Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
   Color       string            `json:"color"`
   Title       string            `json:"title"`
   TitleLink   string            `json:"title_link"`
   Text        string            `json:"text"`
}

func nonemptyAppendNewline(message string) (string) {
   if len(message) == 0 {
      return message
   }

   return message+"\n"
}

func slackParserFunc(r *http.Request) (string, error) {
	// get alert data from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// grafana alert struct
	alert := SlackMessage{}

	// parse body into the alert struct
	err = json.Unmarshal(body, &alert)
	if err != nil {
		return "", err
	}

	// contruct alert message
   message := ""
   hasText := (alert.Text != "")
   if hasText {
      message = alert.Text
   }

   for _, attachment := range alert.Attachments {
      message = nonemptyAppendNewline(message)
      message = message + attachment.Title+": "+attachment.TitleLink+"\n"
      message = message + attachment.Text
   }


	return message, nil
}
