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

func main() {
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")

	jid, err := xmpp.ParseJID(xi)
	fatalOnErr(err)

	addr, err := xmpp.HomeServerAddrs(jid)
	fatalOnErr(err)

	stream, err := xmpp.NewStream(addr[0], nil)
	fatalOnErr(err)

	client, err := xmpp.NewClientXMPP(stream, jid, xp, nil)
	fatalOnErr(err)

	client.Out <- xmpp.Presence{}

	go func() {
		for x := range client.In {
			log.Printf("* recv: %v\n", x)
		}
	}()

	select {}
}
