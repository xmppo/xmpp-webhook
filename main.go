package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/dial"
	mjid "mellium.im/xmpp/jid"
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

func initXMPP(jid mjid.JID, pass string, notls bool) (*xmpp.Session, error) {
	dialer := dial.Dialer{NoTLS: notls}
	conn, err := dialer.Dial(context.TODO(), "tcp", jid)
	if err != nil {
		return nil, err
	}
	return xmpp.NegotiateSession(
		context.TODO(),
		jid.Domain(),
		jid,
		conn,
		false,
		xmpp.NewNegotiator(xmpp.StreamConfig{Features: []xmpp.StreamFeature{
			xmpp.BindResource(),
			xmpp.StartTLS(true, nil),
			xmpp.SASL("", pass, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		}}),
	)
}

func closeXMPP(session *xmpp.Session) {
	session.Close()
	session.Conn().Close()
}

func main() {
	// get xmpp credentials, message receivers and tls settings from env
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")
	xr := os.Getenv("XMPP_RECEIVERS")
	_, notls := os.LookupEnv("XMPP_NO_TLS")

	// check if xmpp credentials and receiver list are supplied
	if len(xi) < 1 || len(xp) < 1 || len(xr) < 1 {
		log.Fatal("XMPP_ID, XMPP_PASS or XMPP_RECEIVERS not set")
	}

	jid, err := mjid.Parse(xi)
	panicOnErr(err)

	session, err := initXMPP(jid, xp, notls)
	panicOnErr(err)
	defer closeXMPP(session)

	panicOnErr(session.Send(context.TODO(), stanza.WrapPresence(mjid.JID{}, stanza.AvailablePresence, nil)))

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

		reply := MessageBody{
			Message: stanza.Message{
				To: msg.From.Bare(),
			},
			Body: msg.Body,
		}

		err = t.Encode(reply)
		if err != nil {
			fmt.Printf("Error responding to message %q: %q", msg.ID, err)
		}
		return nil
	}))

	panicOnErr(err)

	/*// create chan for messages (webhooks -> xmpp)
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

	// listen for requests
	http.ListenAndServe(":4321", nil)*/
}
