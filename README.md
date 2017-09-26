# xmpp-webhook
- Multipurpose XMPP-Webhook (Built for Prometheus/Grafana Alerts)
- Based on https://github.com/atomatt/go-xmpp

## Status
`xmpp-webhook` currently only provides a hook for Grafana. I will implement a `parserFunc` for Prometheus ASAP. Check https://github.com/opthomas-prime/xmpp-webhook/blob/master/handler.go to learn how to support more source services.

## Usage
- `xmpp-webhook` is configured via environment variables:
    - `XMPP_ID` - The JID we want to use
    - `XMPP_PASS` - The password
    - `XMPP_RECEIVERS` - Comma-seperated list of JID's
- After startup, `xmpp-webhook` tries to connect to the XMPP server and provides the implemented HTTP enpoints (on `:4321`). e.g.:

```
curl -X POST -d @grafana-alert.json localhost:4321/grafana
```
- After parsing the body in the appropriate `parserFunc`, the notification is then distributed to the configured receivers.

## Installation
IMPORTANT NOTE: For the sake of simplicity, `xmpp-webhook` is not reconnecting to the XMPP server after a connection-loss. If you use the provided `xmpp-webhook.service` - Systemd will manage the reconnect by restarting the service.

- Download and extract the latest tarball (GitHub release page)
- Install the binary: `install -D -m 744 xmpp-webhook /usr/local/bin/xmpp-webhook`
- Install the service: `install -D -m 644 xmpp-webhook.service /etc/systemd/system/xmpp-webhook.service`
- Configure XMMP credentials in `/etc/xmpp-webhook.env`. e.g.:

```
XMPP_ID='bot@example.com'
XMPP_PASS='passw0rd'
XMPP_RECEIVERS='jdoe@example.com,ops@example.com'
```

- Enable and start the service:

```
systemctl daemon-reload
systemctl enable xmpp-webhook
systemctl start xmpp-webhook
```

## Building
- Get the sources: `go get -u github.com/opthomas-prime/xmpp-webhook`
- Install `dep`: https://github.com/golang/dep
- Change in the project folder: `cd $GOPATH/src/github.com/opthomas-prime/xmpp-webhook`
- Populate `vendor` folder: `dep ensure`
- Build `xmpp-webhook`: `go build`

## Need help?
Feel free to contact me!
