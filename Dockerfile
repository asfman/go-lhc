FROM golang:1.11.0-alpine3.8 as builder

ARG target=go-lhc-api

ADD . /go/src/github.com/asfman/${target}/

RUN go install github.com/asfman/${target}/

FROM alpine:3.8

RUN apk update \
&& apk --update add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/${target} .

ENTRYPOINT ["/root/go-lhc-api"]
