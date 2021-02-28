# xmpp-webhook
- Multipurpose XMPP-Webhook (Built for DevOps Alerts)
- Based on https://github.com/mellium/xmpp

## Status
`xmpp-webhook` currently support:

- Grafana Webhook alerts
- Alertmanager Webhooks
- Slack Incoming Webhooks (Feedback appreciated)

Check https://github.com/tmsmr/xmpp-webhook/blob/master/parser/ to learn how to support more source services.

## Usage
- `xmpp-webhook` is configured via environment variables:
    - `XMPP_ID` - The JID we want to use
    - `XMPP_PASS` - The password
    - `XMPP_RECEIVERS` - Comma-seperated list of JID's
    - `XMPP_SKIP_VERIFY` - Skip TLS verification (Optional)
    - `XMPP_OVER_TLS` - Use dedicated TLS port (Optional)
    - `XMPP_WEBHOOK_LISTEN_ADDRESS` - Bind address (Optional)
- After startup, `xmpp-webhook` tries to connect to the XMPP server and provides the implemented HTTP enpoints. e.g.:

```
curl -X POST -d @dev/grafana-webhook-alert-example.json localhost:4321/grafana
curl -X POST -d @dev/alertmanager-example.json localhost:4321/alertmanager
curl -X POST -d @dev/slack-compatible-notification-example.json localhost:4321/slack
```
- After parsing the body in the appropriate `parserFunc`, the notification is then distributed to the configured receivers.

## Run with Docker
### Build it
- Build image: `docker build -t xmpp-webhook .`
- Run: `docker run -e "XMPP_ID=alerts@example.org" -e "XMPP_PASS=xxx" -e "XMPP_RECEIVERS=receiver1@example.org,receiver2@example.org" -p 4321:4321 -d --name xmpp-webhook xmpp-webhook`
### Use prebuilt image from Docker Hub
- Run: `docker run -e "XMPP_ID=alerts@example.org" -e "XMPP_PASS=xxx" -e "XMPP_RECEIVERS=receiver1@example.org,receiver2@example.org" -p 4321:4321 -d --name xmpp-webhook tmsmr/xmpp-webhook:latest`

## Installation
- Download and extract the latest tarball (GitHub release page)
- Install the binary: `install -D -m 744 xmpp-webhook /usr/local/bin/xmpp-webhook`
- Install the service: `install -D -m 644 xmpp-webhook.service /etc/systemd/system/xmpp-webhook.service`
- Configure XMPP credentials in `/etc/xmpp-webhook.env`. e.g.:

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
- Dependencies are managed via Go Modules (https://github.com/golang/go/wiki/Modules).
- Clone the sources
- Change in the project folder:
- Build `xmpp-webhook`: `go build`
- `dev/xmpp-dev-stack` starts Prosody (With "auth_any" and "roster_allinall" enabled) and two XMPP-clients for easy testing

## Need help?
Feel free to contact me!
