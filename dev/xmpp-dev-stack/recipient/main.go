package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/dial"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
	"os"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type MessageBody struct {
	stanza.Message
	Body string `xml:"body"`
}

func main() {
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")

	if xi == "" || xp == "" {
		log.Fatal("XMPP_ID, XMPP_PASS not set")
	}

	address, err := jid.Parse(xi)
	panicOnErr(err)

	var dialer = dial.Dialer{NoTLS: true}
	conn, err := dialer.Dial(context.TODO(), "tcp", address)
	panicOnErr(err)

	tlsConfig := tls.Config{InsecureSkipVerify: true}

	session, err := xmpp.NegotiateSession(
		context.TODO(),
		address.Domain(),
		address,
		conn,
		false,
		xmpp.NewNegotiator(xmpp.StreamConfig{Features: []xmpp.StreamFeature{
			xmpp.BindResource(),
			xmpp.StartTLS(true, &tlsConfig),
			xmpp.SASL("", xp, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		}}),
	)
	panicOnErr(err)

	fmt.Println("connected")

	err = session.Send(context.TODO(), stanza.WrapPresence(address, stanza.AvailablePresence, nil))
	panicOnErr(err)

	err = session.Serve(xmpp.HandlerFunc(func(t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
		d := xml.NewTokenDecoder(t)
		if start.Name.Local != "message" {
			return nil
		}

		msg := MessageBody{}
		err = d.DecodeElement(&msg, start)
		if err != nil && err != io.EOF {
			return nil
		}

		if msg.Body == "" || msg.Type != stanza.ChatMessage {
			return nil
		}

		fmt.Printf("%s: %s\n", msg.From, msg.Body)

		return nil
	}))
	panicOnErr(err)
}
