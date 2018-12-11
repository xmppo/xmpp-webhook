package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/emgee/go-xmpp/src/xmpp"
)

// starts xmpp session and returns the xmpp client
func xmppLogin(id string, pass string) (*xmpp.XMPP, error) {
	// parse jid structure
	jid, err := xmpp.ParseJID(id)
	if err != nil {
		return nil, err
	}

	// extract/generate address:port from jid
	addr, err := xmpp.HomeServerAddrs(jid)
	if err != nil {
		return nil, err
	}

	// create xml stream to address
	stream, err := xmpp.NewStream(addr[0], nil)
	if err != nil {
		return nil, err
	}

	// create client (login)
	client, err := xmpp.NewClientXMPP(stream, jid, pass, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	// get xmpp credentials and message receivers from env
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")
	xr := os.Getenv("XMPP_RECEIVERS")

	// check if xmpp credentials and receiver list are supplied
	if len(xi) < 1 || len(xp) < 1 || len(xr) < 1 {
		log.Fatal("XMPP_ID, XMPP_PASS or XMPP_RECEIVERS not set")
	}

	// connect to xmpp server
	xc, err := xmppLogin(xi, xp)
	if err != nil {
		log.Fatal(err)
	}

	// announce initial presence
	xc.Out <- xmpp.Presence{}

	// listen for incoming xmpp stanzas
	go func() {
		for stanza := range xc.In {
			// check if stanza is a message
			m, ok := stanza.(*xmpp.Message)
			if ok && len(m.Body) > 0 {
				// echo the message
				xc.Out <- xmpp.Message{
					To:   m.From,
					Body: m.Body,
				}
			}
		}
		// xc.In is closed when the server closes the stream
		log.Fatal("connection lost")
	}()

	// create chan for messages (webhooks -> xmpp)
	messages := make(chan string)

	// wait for messages from the webhooks and send them to all receivers
	go func() {
		for m := range messages {
			for _, r := range strings.Split(xr, ",") {
				xc.Out <- xmpp.Message{
					To: r,
					Body: []xmpp.MessageBody{
						{
							Value: m,
						},
					},
				}
			}
		}
	}()

	// initialize handler for grafana alerts
	http.Handle("/grafana", newMessageHandler(messages, grafanaParserFunc))

	// initialize handler for grafana alerts
	http.Handle("/prometheus", newMessageHandler(messages, prometheusParserFunc))

	// listen for requests
	http.ListenAndServe(":4321", nil)
}
