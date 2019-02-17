FROM golang:1.11-alpine as build
RUN apk update && apk add --no-cache git gcc libc-dev
RUN adduser -D gitlab
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /go/src/github.com/sotomskir/gitlab-cli/
ADD . .
RUN go get -d -t -v ./...
RUN go test ./...
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM alpine
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/src/github.com/sotomskir/gitlab-cli/gitlab-cli /usr/local/bin/gitlab-cli
USER gitlab
