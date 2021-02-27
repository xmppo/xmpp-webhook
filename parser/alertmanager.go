package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func AlertmanagerParserFunc(r *http.Request) (string, error) {
	// get alert data from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New(readErr)
	}

	payload := &struct {
		Alerts []struct {
			Status      string            `json:"status"`
			Labels      map[string]string `json:"labels"`
			Annotations map[string]string `json:"annotations"`
		} `json:"alerts"`
	}{}

	// parse body into the alert struct
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return "", errors.New(parseErr)
	}

	// contruct alert message
	var message string
	for _, alert := range payload.Alerts {
		if alert.Status == "resolved" {
			message = "Resolved" + "\n"
		} else {
			message = "Firing" + "\n"
		}

		message += "Labels" + "\n"
		for key, label := range alert.Labels {
			message += fmt.Sprintf("%s = %s\n", key, label)
		}

		message += "Annotations" + "\n"
		for key, annotation := range alert.Annotations {
			message += fmt.Sprintf("%s = %s\n", key, annotation)
		}

		message += "\n"
	}

	return message, nil
}
