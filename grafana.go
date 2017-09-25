package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
{
  "title": "My alert",
  "ruleId": 1,
  "ruleName": "Load peaking!",
  "ruleUrl": "http://url.to.grafana/db/dashboard/my_dashboard?panelId=2",
  "state": "alerting",
  "imageUrl": "http://s3.image.url",
  "message": "Load is peaking. Make sure the traffic is real and spin up more webfronts",
  "evalMatches": [
    {
      "metric": "requests",
      "tags": {},
      "value": 122
    }
  ]
}
*/

type GrafanaAlert struct {
	Title    string `json:"title"`
	RuleName string `json:"ruleName"`
	RuleUrl  string `json:"ruleUrl"`
	State    string `json:"state"`
	Message  string `json:"message"`
}

func makeGrafanaHandler(messages chan<- string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		fmt.Printf("%v", string(body))
		if err == nil {
			var alert GrafanaAlert
			err = json.Unmarshal(body, &alert)
			if err == nil {
				message := alert.State + ": " + alert.Title + "/" + alert.Message + "(" + alert.RuleUrl + ")"
				messages <- message
			}
		}
	}
}
