package main

import (
	"log"
	"os"

	"github.com/emgee/go-xmpp/src/xmpp"
)

const (
	xmppBotAnswer = "im a dumb bot"
)

// helper function for error checks
func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// starts xmpp session, sends initial presence and returns the xmpp client
func xmppLogin(id string, pass string) *xmpp.XMPP {
	// parse jid structure
	jid, err := xmpp.ParseJID(id)
	fatalOnErr(err)

	// extract/generate address:port from jid
	addr, err := xmpp.HomeServerAddrs(jid)
	fatalOnErr(err)

	// create xml stream to address
	stream, err := xmpp.NewStream(addr[0], nil)
	fatalOnErr(err)

	// create client (login)
	client, err := xmpp.NewClientXMPP(stream, jid, pass, nil)
	fatalOnErr(err)

	// announce presence
	client.Out <- xmpp.Presence{}

	return client
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
}

func main() {
	// get xmpp credentials from ENV
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")

	// check if xmpp credentials are supplied
	if len(xi) < 1 || len(xp) < 1 {
		log.Fatal("XMPP_ID or XMPP_PASS not set")
	}

	// start xmpp client
	xc := xmppLogin(xi, xp)

	// dispatch incoming stanzas to handler
	go handleXMPPStanza(xc.In, xc.Out)

	select {}
}
