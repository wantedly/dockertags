FROM alpine:3.4

ENV GOPATH /go

RUN apk add --no-cache --update ca-certificates

COPY . /go/src/github.com/wantedly/dockertags

RUN apk add --no-cache --update --virtual=build-deps go git make mercurial \
    && cd /go/src/github.com/wantedly/dockertags \
    && make \
    && cp bin/dockertags /dockertags \
    && cd / \
    && rm -rf /go \
    && apk del build-deps

ENTRYPOINT ["/dockertags"]
