FROM golang:1.13 as builder

ADD . /go/pocc
#ADD ./netrc /root/.netrc

WORKDIR /go/pocc

ENV GOPROXY=https://goproxy.io
ENV GOPRIVATE=*.bhex.io

RUN CGO_ENABLED=0 go build -tags netgo -v -o /go/bin/pocd ./cmd/pocd \
    && CGO_ENABLED=0 go build -tags netgo -v -o /go/bin/poccli ./cmd/poccli




FROM alpine:latest

WORKDIR /go/
COPY --from=builder /go/bin/pocd /go/
COPY --from=builder /go/bin/poccli /go/
COPY --from=builder /go/pocc/run.sh /go/

# p2p port
EXPOSE 26656
# RPC port
EXPOSE 26657

VOLUME [ "/root/.pocd", "/root/.poccli " ]

ENTRYPOINT [ "sh", "run.sh" ]
