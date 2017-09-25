package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emgee/go-xmpp/src/xmpp"
)

const (
	envXMPPID        = "XMPP_ID"
	envXMPPPASS      = "XMPP_PASS"
	envXMPPReceivers = "XMPP_RECEIVERS"
	errWrongArgs     = "XMPP_ID, XMPP_PASS or XMPP_RECEIVERS not set"
	xmppBotAnswer    = "im a dumb bot"
	xmppConnErr      = "failed to connect "
	xmppOfflineErr   = "not connected to XMPP server, dropped message"
	xmppFailedPause  = 30
	webHookAddr      = ":4321"
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

// creates MessageBody slice suitable for xmpp.Message
func xmppBodyCreate(message string) []xmpp.MessageBody {
	return []xmpp.MessageBody{
		xmpp.MessageBody{
			Value: message,
		},
	}
}

// handles incoming stanzas
func handleXMPPStanza(in <-chan interface{}, out chan<- interface{}) {
	for stanza := range in {
		// check if stanza is a message
		message, ok := stanza.(*xmpp.Message)
		if ok && len(message.Body) > 0 {
			// send constant as answer
			out <- xmpp.Message{
				To:   message.From,
				Body: xmppBodyCreate(xmppBotAnswer),
			}
		}
	}
	// func returns when in chan is closed (server terminated stream)
}

func main() {
	// get xmpp credentials from ENV
	xi := os.Getenv(envXMPPID)
	xp := os.Getenv(envXMPPPASS)
	xr := os.Getenv(envXMPPReceivers)

	// check if xmpp credentials and receiver list are supplied
	if len(xi) < 1 || len(xp) < 1 || len(xr) < 1 {
		log.Fatal(errWrongArgs)
	}

	// connect xmpp client and observe connection - reconnect if needed
	var xc *xmpp.XMPP
	go func() {
		for {
			// try to connect to xmpp server
			var err error
			xc, err = xmppLogin(xi, xp)
			if err != nil {
				// report failure and wait
				log.Print(xmppConnErr, err)
				time.Sleep(time.Second * time.Duration(xmppFailedPause))
			} else {
				// send initial presence and dispatch channels to handler for incoming messages
				xc.Out <- xmpp.Presence{}
				handleXMPPStanza(xc.In, xc.Out)
			}
		}
	}()

	// create channel for alerts (Webhook -> XMPP)
	alertChan := make(chan string)

	// create handler for outgoing XMPP messages
	go func() {
		for message := range alertChan {
			for _, r := range strings.Split(xr, ",") {
				if xc != nil {
					xc.Out <- xmpp.Message{
						To:   r,
						Body: xmppBodyCreate(message),
					}
				} else {
					log.Print(xmppOfflineErr)
				}
			}
		}
	}()

	// initialize HTTP handlers with chan for alerts
	grafanaHandler := makeGrafanaHandler(alertChan)
	http.HandleFunc("/grafana", grafanaHandler)

	// listen for requests
	http.ListenAndServe(webHookAddr, nil)
}
