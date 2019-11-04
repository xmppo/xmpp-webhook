FROM golang:1.13-alpine3.10 as builder
RUN apk add --no-cache git
COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.10
COPY --from=builder /build/recipient /recipient
CMD /recipient
