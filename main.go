package main

import (
	"log"
	"os"
	"time"

	"github.com/emgee/go-xmpp/src/xmpp"
)

const (
	envXMPPID       = "XMPP_ID"
	envXMPPPASS     = "XMPP_PASS"
	errWrongArgs    = "XMPP_ID or XMPP_PASS not set"
	xmppBotAnswer   = "im a dumb bot"
	xmppConnErr     = "failed to connect "
	xmppFailedPause = 30
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

	// check if xmpp credentials are supplied
	if len(xi) < 1 || len(xp) < 1 {
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
				// send initial presence and dispatch channels to handler
				xc.Out <- xmpp.Presence{}
				handleXMPPStanza(xc.In, xc.Out)
			}
		}
	}()

	select {}
}
