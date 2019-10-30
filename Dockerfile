FROM golang:1.13-alpine3.10 as builder
MAINTAINER Thomas Maier <contact@thomas-maier.net>
RUN apk add --no-cache git
COPY . /build
WORKDIR /build
RUN GOOS=linux GOARCH=amd64 go build

FROM alpine:3.10
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/xmpp-webhook /xmpp-webhook
RUN adduser -D -g '' xmpp-webhook
USER xmpp-webhook
ENV XMPP_ID="" \
	XMPP_PASS="" \
	XMPP_RECEIVERS=""
EXPOSE 4321
CMD ["/xmpp-webhook"]
