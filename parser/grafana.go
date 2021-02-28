package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func GrafanaParserFunc(r *http.Request) (string, error) {
	// get alert data from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New(readErr)
	}

	alert := &struct {
		Title   string `json:"title"`
		RuleURL string `json:"ruleUrl"`
		State   string `json:"state"`
		Message string `json:"message"`
	}{}

	// parse body into the alert struct
	err = json.Unmarshal(body, &alert)
	if err != nil {
		return "", errors.New(parseErr)
	}

	// construct alert message
	var message string
	switch alert.State {
	case "ok":
		message = ":) " + alert.Title
	default:
		message = ":( " + alert.Title + "\n\n"
		message += alert.Message + "\n\n"
		message += alert.RuleURL
	}

	return message, nil
}
