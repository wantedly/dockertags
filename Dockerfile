FROM golang:1.14 AS build
WORKDIR /go/src/github.com/wantedly/dockertags
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install -tags netgo -ldflags "-extldflags -static" .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/bin/dockertags /bin/dockertags
ENTRYPOINT ["/bin/dockertags"]
