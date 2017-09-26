# xmpp-webhook
- Multipurpose XMPP Webhook (Built for Prometheus/Grafana Alerts)
- Based on [](https://github.com/atomatt/go-xmpp)

## Status
`xmpp-webhook` currently only provides a hook for Grafana. I will implement a `parserFunc` for Prometheus ASAP. Check [](https://github.com/opthomas-prime/xmpp-webhook/blob/master/handler.go) to learn how to support more source services.

## Usage
- `xmpp-webhook` is configured via environment variables:
    - `XMPP_ID`
    - `XMPP_PAS`
    - `XMPP_RECEIVERS`
- 

```
XMPP_ID='bot@example.com'
XMPP_PASS='passw0rd'
XMPP_RECEIVERS='jdoe@example.com,ops@example.com'

/etc/systemd/system/xmpp-webhook.service

https://github.com/golang/dep
go get -u github.com/golang/dep/cmd/dep
```
