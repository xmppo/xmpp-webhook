# xmpp-webhook
- Multipurpose XMPP-Webhook (Built for Prometheus/Grafana Alerts)
- Based on https://github.com/atomatt/go-xmpp

## Status
`xmpp-webhook` currently only provides a hook for Grafana. I will implement a `parserFunc` for Prometheus ASAP. Check https://github.com/opthomas-prime/xmpp-webhook/blob/master/handler.go to learn how to support more source services.

## Usage
- `xmpp-webhook` is configured via environment variables:
    - `XMPP_ID` - The JID we want to use
    - `XMPP_PAS` - The password
    - `XMPP_RECEIVERS` - Comma-seperated list of JID's
- After startup `xmpp-webhooks` tries to connect to the XMPP server and provides the implemented HTTP enpoints (on `:4321`). e.g.:

```
curl -X POST -d @grafana-alert.json localhost:4321
```
- After parsing the body in the appropriate `parserFunc`, the notification is then distributed to the configured receivers.

```
XMPP_ID='bot@example.com'
XMPP_PASS='passw0rd'
XMPP_RECEIVERS='jdoe@example.com,ops@example.com'

/etc/systemd/system/xmpp-webhook.service

https://github.com/golang/dep
go get -u github.com/golang/dep/cmd/dep
```
