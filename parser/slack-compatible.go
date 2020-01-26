package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func SlackParserFunc(r *http.Request) (string, error) {
	// get alert data from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New(readErr)
	}

	alert := struct {
		Text        string `json:"text"`
		Attachments []struct {
			Title     string `json:"title"`
			TitleLink string `json:"title_link"`
			Text      string `json:"text"`
		} `json:"attachments"`
	}{}

	// parse body into the alert struct
	err = json.Unmarshal(body, &alert)
	if err != nil {
		return "", errors.New(parseErr)
	}

	// contruct alert message
	message := alert.Text
	for _, attachment := range alert.Attachments {
		if len(message) > 0 {
			message = message + "\n"
		}
		message += attachment.Title + "\n"
		message += attachment.TitleLink + "\n\n"
		message += attachment.Text
	}

	return message, nil
}
