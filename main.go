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

func initXMPP(address jid.JID, pass string, skipTLSVerify bool, legacyTLS bool, forceStartTLS bool) (*xmpp.Session, error) {
	tlsConfig := tls.Config{InsecureSkipVerify: skipTLSVerify}
	var dialer dial.Dialer
	if skipTLSVerify {
		dialer = dial.Dialer{NoTLS: !legacyTLS, TLSConfig: &tlsConfig}
	} else {
		dialer = dial.Dialer{NoTLS: !legacyTLS}
	}
	conn, err := dialer.Dial(context.TODO(), "tcp", address)
	if err != nil {
		return nil, err
	}
	if !skipTLSVerify {
		tlsConfig.ServerName = address.Domainpart()
	}
	return xmpp.NegotiateSession(
		context.TODO(),
		address.Domain(),
		address,
		conn,
		false,
		xmpp.NewNegotiator(xmpp.StreamConfig{Features: []xmpp.StreamFeature{
			xmpp.BindResource(),
			xmpp.StartTLS(forceStartTLS, &tlsConfig),
			xmpp.SASL("", pass, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		}}),
	)
}

func closeXMPP(session *xmpp.Session) {
	session.Close()
	session.Conn().Close()
}

func main() {
	// get xmpp credentials, message receivers
	xi := os.Getenv("XMPP_ID")
	xp := os.Getenv("XMPP_PASS")
	xr := os.Getenv("XMPP_RECEIVERS")

	// get tls settings from env
	_, skipTLSVerify := os.LookupEnv("XMPP_SKIP_VERIFY")
	_, legacyTLS := os.LookupEnv("XMPP_OVER_TLS")
	_, forceStartTLS := os.LookupEnv("XMPP_FORCE_STARTTLS")

	// check if xmpp credentials and receiver list are supplied
	if xi == "" || xp == "1" || xr == "" {
		log.Fatal("XMPP_ID, XMPP_PASS or XMPP_RECEIVERS not set")
	}

	address, err := jid.Parse(xi)
	panicOnErr(err)

	// connect to xmpp server
	session, err := initXMPP(address, xp, skipTLSVerify, legacyTLS, forceStartTLS)
	panicOnErr(err)
	defer closeXMPP(session)

	// send initial presence
	panicOnErr(session.Send(context.TODO(), stanza.WrapPresence(jid.JID{}, stanza.AvailablePresence, nil)))

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
