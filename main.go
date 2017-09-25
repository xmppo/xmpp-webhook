package main

import (
	"log"
	"os"

	"github.com/emgee/go-xmpp/src/xmpp"
)

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func xmppLogin(id string, pass string) *xmpp.XMPP {
	jid, err := xmpp.ParseJID(id)
	fatalOnErr(err)

	addr, err := xmpp.HomeServerAddrs(jid)
	fatalOnErr(err)

	stream, err := xmpp.NewStream(addr[0], nil)
	fatalOnErr(err)

	client, err := xmpp.NewClientXMPP(stream, jid, pass, nil)
	fatalOnErr(err)

	client.Out <- xmpp.Presence{}

	return client
}

func main() {
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")

	if len(xi) < 1 || len(xp) < 1 {
		log.Fatal("XMPP_ID or XMPP_PASS not set")
	}

	xc := xmppLogin(xi, xp)

	go func() {
		for msg := range xc.In {
			log.Printf("* recv: %v\n", msg)
		}
	}()

	select {}
}
