FROM golang:1.15-alpine3.13 as builder
RUN apk add --no-cache git
COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.13
COPY --from=builder /build/recipient /recipient
CMD /recipient
